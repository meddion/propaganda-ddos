package pkg

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
)

type (
	BotMsg struct {
		ID       int
		Err      error
		ErrCount int
		Continue chan<- bool
	}

	Bot struct {
		c  *http.Client
		id int
	}
)

func NewBot(id int, getProxy func() *Proxy) (*Bot, error) {
	// TODO: adjust constants
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	tr.MaxIdleConns = 10   // 0 - no limit
	tr.MaxConnsPerHost = 0 // 0 - no limit
	tr.IdleConnTimeout = DIAL_TIMEOUT
	tr.ReadBufferSize = 10
	tr.DisableCompression = true
	tr.DisableKeepAlives = true

	tr.Proxy = nil

	// tr.Proxy = func(req *http.Request) (*url.URL, error) {
	// 	proxy := getProxy()
	// 	// proxy := (*Proxy)(nil)
	// 	if proxy == nil {
	// 		log.WithField("id", id).Println("Надсилаю без проксі")
	// 		return nil, nil
	// 	}

	// 	url, err := url.Parse(proxy.IP)

	// 	if err != nil {
	// 		url, err = url.Parse(req.URL.Scheme + proxy.IP)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// 	a := strings.Split(proxy.Auth, ":")
	// 	if len(a) > 2 {
	// 		req.SetBasicAuth(a[0], a[1])
	// 	}

	// 	return url, nil
	// }

	b := &Bot{
		id: id,
		c:  &http.Client{Transport: tr},
	}

	return b, nil
}

func newReq(ctx context.Context, method string, target string) (*http.Request, error) {
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
	select {
	case <-ctx.Done():
		return
	default:
	}

	cont := make(chan bool, 1)
	errCount := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			req, err := newReq(ctx, "GET", target)
			switch {
			case err != nil:
				errCount++
				err = fmt.Errorf("Під час створення запиту: %v", err)
			default:
				var resp *http.Response
				switch resp, err = b.c.Do(req); {
				case err != nil:
					errCount++
					if resp != nil {
						err = fmt.Errorf("%v (%d %s)", err, resp.StatusCode, http.StatusText(resp.StatusCode))
					}
					// else if !errors.Is(err, context.Canceled) {
					// }
				default:
					if resp.StatusCode != http.StatusOK {
						errCount++
						err = fmt.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
					}
					toDevNull(resp.Body)
				}
			}

			select {
			case <-ctx.Done():
				return
			case msgs <- BotMsg{ID: b.id, Err: err, Continue: cont, ErrCount: errCount}:
				select {
				case <-ctx.Done():
					return
				case v := <-cont:
					if !v {
						return
					}
				}
			}

		}
		runtime.Gosched()
	}
}
