package pkg

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	DIAL_TIMEOUT    = time.Second * 10
	TEST_PROXY_SITE = "https://google.com"
)
const (
	APIv1 = iota + 1
	APIv2
)

type (
	Target struct {
		ID           int    `json:"-"`
		URL          string `json:"url"`
		Page         string `json:"page"`
		NeedParseURL int    `json:"need_parse_url"`
		Attack       int    `json:"atack"`
	}

	Proxy struct {
		ID   int    `json:"id"`
		IP   string `json:"ip"`
		Auth string `json:"auth"`
	}

	srcDataAPIv1 struct {
		Target  Target  `json:"site"`
		Proxies []Proxy `json:"proxy"`
	}

	srcDataAPIv2 struct {
		Targets []Target `json:"site"`
		Proxies []Proxy  `json:"proxy"`
	}
)

var DefClient *http.Client

func init() {
	// Setup
	rand.Seed(time.Now().UnixNano())
	// Default http client
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// tr.IdleConnTimeout = time.Second * 5
	tr.ResponseHeaderTimeout = time.Second * 5
	tr.MaxConnsPerHost = 0 // no limit
	tr.ReadBufferSize = 100
	tr.DisableKeepAlives = true
	tr.Dial = (&net.Dialer{
		Timeout: DIAL_TIMEOUT,
	}).Dial
	DefClient = &http.Client{Transport: tr, Timeout: DIAL_TIMEOUT}
	// Logger
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}

func GetSrcFromAPIGateway(rootCtx context.Context, gateway string) (string, error) {
	gCtx, cancel := context.WithTimeout(rootCtx, DIAL_TIMEOUT)
	defer cancel()

	req, err := newReq(gCtx, "GET", gateway)
	if err != nil {
		return "", err
	}

	resp, err := DefClient.Do(req)
	if err != nil {
		return "", err
	}
	defer toDevNull(resp.Body)

	dec := json.NewDecoder(resp.Body)
	var sources []string
	if err := dec.Decode(&sources); err != nil {
		return "", err
	}

	return sources[rand.Intn(len(sources)-1)], nil
}

func GetProxyFromAPI(rootCtx context.Context, proxySrc string) ([]Proxy, error) {
	srcCtx, cancel := context.WithTimeout(rootCtx, DIAL_TIMEOUT)
	defer cancel()

	req, err := newReq(srcCtx, "GET", proxySrc)
	if err != nil {
		return nil, err
	}

	resp, err := DefClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer toDevNull(resp.Body)

	dec := json.NewDecoder(resp.Body)
	var proxies []Proxy
	if err := dec.Decode(&proxies); err != nil {
		return nil, err
	}

	return proxies, nil
}

func GetTargetsFromAPI(rootCtx context.Context, targetSrc string) ([]Target, error) {
	srcCtx, cancel := context.WithTimeout(rootCtx, DIAL_TIMEOUT)
	defer cancel()

	req, err := newReq(srcCtx, "GET", targetSrc)
	if err != nil {
		return nil, err
	}

	resp, err := DefClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer toDevNull(resp.Body)

	dec := json.NewDecoder(resp.Body)
	var targets []Target
	if err := dec.Decode(&targets); err != nil {
		return nil, err
	}

	return targets, nil
}

func GetDataFromAPISrc(rootCtx context.Context, src string, apiVer int) ([]Target, []Proxy, error) {
	srcCtx, cancel := context.WithTimeout(rootCtx, DIAL_TIMEOUT)
	defer cancel()

	req, err := newReq(srcCtx, "GET", src)
	if err != nil {
		return nil, nil, err
	}

	resp, err := DefClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer toDevNull(resp.Body)

	return DecodeJSON(resp.Body, apiVer)
}

func GetDataFromFile(rootCtx context.Context, path string, apiVer int) ([]Target, []Proxy, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	return DecodeJSON(bufio.NewReader(file), apiVer)
}

func DecodeJSON(content io.Reader, apiVer int) ([]Target, []Proxy, error) {
	dec := json.NewDecoder(content)
	switch apiVer {
	case APIv1:
		srcResp := new(srcDataAPIv1)
		if err := dec.Decode(srcResp); err != nil {
			return nil, nil, err
		}
		return []Target{srcResp.Target}, srcResp.Proxies, nil
	case APIv2:
		srcResp := new(srcDataAPIv2)
		if err := dec.Decode(srcResp); err != nil {
			return nil, nil, err
		}
		return srcResp.Targets, srcResp.Proxies, nil
	}

	log.Fatalf("Версія APIv%d не підтримується\n", apiVer)
	return nil, nil, nil
}

func ValidateTargets(ctxRoot context.Context, targets []Target, dnsResolv bool) []Target {
	var wg sync.WaitGroup
	wg.Add(len(targets))
	res := make(chan Target, len(targets))
	validaTargets := make([]Target, 0, len(targets))

	for _, target := range targets {
		go func(target Target) {
			defer wg.Done()

			if target.URL == "" {
				if target.Page != "" {
					target.URL = target.Page
				} else {
					log.Errorln("empty address; no target")
					return
				}
			}

			if url, err := ValidateAddress(ctxRoot, target.URL, dnsResolv); err != nil {
				log.Errorf("on parsing target url %s: %v", url, err)
				return
			} else {
				target.URL = url
				res <- target
			}
		}(target)
	}
	go func() {
		wg.Wait()
		close(res)
	}()

	for t := range res {
		validaTargets = append(validaTargets, t)
	}

	return validaTargets
}

func ValidateAddress(ctxRoot context.Context, addr string, dnsResolve bool) (string, error) {
	CleanupURL := func(targetURL string) string {
		targetURL = strings.Trim(targetURL, "\r")
		targetURL = strings.Trim(targetURL, "\n")
		return strings.TrimFunc(targetURL, func(r rune) bool {
			return r == ' ' || r == '\n' || r == '\r'
		})
	}

	addr = CleanupURL(addr)
	if !strings.Contains(addr, "http") {
		addr = "http://" + addr
	}

	url, err := url.Parse(addr)
	if err != nil {
		return "", err
	}

	var port string
	pair := strings.Split(url.Host, ":")
	if len(pair) > 1 {
		port = pair[1]
	}
	host := pair[0]

	// if dnsResolve && net.ParseIP(host) == nil {}
	if dnsResolve && !isIP(host) && isDNS(host) {
		if host, err = resolveHost(ctxRoot, host); err != nil {
			return "", err
		}
	}

	if port != "" {
		url.Host = host + ":" + port
	} else {
		url.Host = host
	}

	return url.String(), nil
}

// Taken from: https://github.com/Arriven/db1000n
const (
	IP_REGEX  = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"
	DNS_REGEX = "^(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z]|[A-Za-z][A-Za-z0-9\\-]*[A-Za-z0-9])$"
)

func isIP(host string) bool {
	res, err := regexp.MatchString(IP_REGEX, host)
	if err != nil {
		// log.Errorf("Під час зіставлення IP regex із %s: %v", host, err)
		return false
	}

	return res
}

func isDNS(host string) bool {
	res, err := regexp.MatchString(DNS_REGEX, host)
	if err != nil {
		// log.Fatalf("ід час зіставлення DNS regex із %s: %v", host, err)
		return false
	}

	return res
}

// resolveHost function gets a string and returns the ip address while deciding it is an ip address already or DNS record
func resolveHost(ctxRoot context.Context, host string) (string, error) {
	ipRecords, err := net.DefaultResolver.LookupIP(ctxRoot, "ip4", host)
	if err != nil {
		return "", err
	}
	ip := ipRecords[0].String()
	if ip == "127.0.0.1" || ip == "0.0.0.0" {
		return "", errors.New("couldn't resolve")
	}

	return ip, nil
}

func ProxyValidation(rootCtx context.Context, proxies []Proxy) []Proxy {
	checkProxy := func(rootCtx context.Context, proxy Proxy) error {
		bot, err := NewBot(0, func() *Proxy { return &proxy })
		if err != nil {
			return err
		}

		proxyCtx, cancel := context.WithTimeout(rootCtx, time.Second*2)
		defer cancel()

		msgs := make(chan BotMsg)
		go bot.Start(proxyCtx, TEST_PROXY_SITE, msgs)

		select {
		case <-proxyCtx.Done():
			return fmt.Errorf("ніякої відповіді від проксі: %s", proxy.IP)
		case msg := <-msgs:
			return msg.Err
		}
	}

	if len(proxies) == 0 {
		return proxies
	}

	validProxies := make([]Proxy, 0, len(proxies))

	type result struct {
		p Proxy
		e error
	}
	resChan := make(chan result)
	var wg sync.WaitGroup
	wg.Add(len(proxies))

	go func() {
		wg.Wait()
		close(resChan)
	}()
	for _, proxy := range proxies {
		go func(proxy Proxy) {
			defer wg.Done()
			resChan <- result{proxy, checkProxy(rootCtx, proxy)}
		}(proxy)
	}

	for res := range resChan {
		if res.e != nil {
			log.Errorln(res.e)
		} else {
			validProxies = append(validProxies, res.p)
		}
	}

	return validProxies
}
