package pkg

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type (
	BotMsg struct {
		ID       int
		Err      error
		ErrCount int
		Done     func()
	}

	Bot struct {
		c         *http.Client
		req       *http.Request
		id        int
		withProxy bool
	}
)

func NewBot(id int, proxy *Proxy, withTLSproxy bool) (*Bot, error) {
	// TODO: adjust constants
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	tr.MaxIdleConns = 0    // 0 - no limit
	tr.MaxConnsPerHost = 0 // 0 - no limit
	tr.IdleConnTimeout = 0
	tr.ReadBufferSize = 10
	tr.DisableCompression = true

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

		proxyURL = CleanupURL(proxyURL)
		url, err := ValidateURL(proxyURL)
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

func (b *Bot) Start(ctx context.Context, target string, msgs chan<- BotMsg) {
	log := log.WithField("bot", b.id)

	select {
	case <-ctx.Done():
		log.Println("is canceld")
	default:
	}

	var once sync.Once
	botDone := make(chan struct{}, 1)
	termBot := func() {
		once.Do(func() {
			close(botDone)
		})
	}

	errCount := 0
	for {
		select {
		case <-botDone:
			return
		case <-ctx.Done():
			return
		default:
			var err error

			if req, err := newGetReq(ctx, target); err != nil {
				errCount++
				err = fmt.Errorf("Під час створення запиту: %v", err)
			} else if resp, err := b.c.Do(req); err != nil {
				errCount++
				if resp != nil {
					err = fmt.Errorf("%v (%d %s)", err, resp.StatusCode, http.StatusText(resp.StatusCode))
				} else if !errors.Is(err, context.Canceled) {
					err = fmt.Errorf("%v", err)
				}
			} else {
				if resp.StatusCode != http.StatusOK {
					errCount++
					err = fmt.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
				}
				toDevNull(resp.Body)
			}

			select {
			case <-botDone:
				return
			case <-ctx.Done():
				return
			case msgs <- BotMsg{ID: b.id, Err: err, Done: termBot, ErrCount: errCount}:
			}

		}
		runtime.Gosched()
	}
}
