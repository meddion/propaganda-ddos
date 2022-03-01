package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strings"
)

type (
	Site struct {
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

	TargetData struct {
		Site    Site    `json:"site"`
		Proxies []Proxy `json:"proxy"`
	}
)

func NewTargetData(targetURL string, proxy ...Proxy) *TargetData {
	return &TargetData{
		Site: Site{
			URL: targetURL,
		},
		Proxies: proxy,
	}
}

func GetSrcURL(ctx context.Context, gateway string) (string, error) {
	req, err := newGetReq(ctx, gateway)
	if err != nil {
		return "", err
	}

	resp, err := DefClient.Do(req)
	if err != nil {
		return "", err
	}
	defer toDevNull(resp.Body)

	dec := json.NewDecoder(resp.Body)
	var hosts []string
	if err := dec.Decode(&hosts); err != nil {
		return "", err
	}

	return hosts[rand.Intn(len(hosts)-1)], nil
}

func GetDataFromSrc(ctx context.Context, src string) (*TargetData, error) {
	req, err := newGetReq(ctx, src)
	if err != nil {
		return nil, err
	}

	resp, err := DefClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer toDevNull(resp.Body)

	dec := json.NewDecoder(resp.Body)
	hostResp := new(TargetData)
	if err := dec.Decode(hostResp); err != nil {
		return nil, err
	}

	return hostResp, err
}

func ValidateTargetData(data *TargetData) error {
	if data.Site.URL == "" {
		if data.Site.Page != "" {
			data.Site.URL = data.Site.Page
		} else {
			return errors.New("empty address; no target")
		}
	}

	data.Site.URL = CleanupURL(data.Site.URL)

	if _, err := ValidateURL(data.Site.URL); err != nil {
		return fmt.Errorf("on parsing target url: %w", err)
	}

	return nil
}

func CleanupURL(targetURL string) string {
	targetURL = strings.Trim(targetURL, "\r")
	targetURL = strings.Trim(targetURL, "\n")
	return strings.TrimFunc(targetURL, func(r rune) bool {
		return r == ' ' || r == '\n' || r == '\r'
	})
}

func ValidateURL(targetURL string) (*url.URL, error) {
	url, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return nil, err
	}

	if net.ParseIP(url.Host) == nil && !strings.Contains(url.Host, ".") {
		return nil, errors.New("invalid host")
	}

	return url, nil
}
