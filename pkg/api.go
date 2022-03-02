package pkg

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	GATEWAY_TIMEOUT = time.Second * 10
	SRC_TIMEOUT     = time.Second * 10
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

func GetSrcFromAPIGateway(rootCtx context.Context, gateway string) (string, error) {
	gCtx, cancel := context.WithTimeout(rootCtx, GATEWAY_TIMEOUT)
	defer cancel()

	req, err := newGetReq(gCtx, gateway)
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

func GetDataFromAPISrc(rootCtx context.Context, src string, apiVer int) ([]Target, []Proxy, error) {
	srcCtx, cancel := context.WithTimeout(rootCtx, SRC_TIMEOUT)
	defer cancel()

	req, err := newGetReq(srcCtx, src)
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

func ValidateTarget(target *Target) error {
	if target.URL == "" {
		if target.Page != "" {
			target.URL = target.Page
		} else {
			return errors.New("empty address; no target")
		}
	}

	target.URL = CleanupURL(target.URL)

	if _, err := ValidateURL(target.URL); err != nil {
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
