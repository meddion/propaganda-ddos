package pkg

import (
	"context"
	"encoding/json"
	"math/rand"
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
		Site
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
