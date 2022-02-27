package main

import (
	"encoding/json"
	"math/rand"
)

type (
	Proxy struct {
		ID   int    `json:"id"`
		IP   string `json:"ip"`
		Auth string `json:"auth"`
	}

	HostResp struct {
		Site struct {
			ID           int    `json:"-"`
			URL          string `json:"url"`
			Page         string `json:"page"`
			NeedParseURL int    `json:"need_parse_url"`
			Attack       int    `json:"atack"`
		}

		Proxies []Proxy `json:"proxy"`
	}
)

func GetSrcURL() (string, error) {
	resp, err := DefClient.Get(HOST_AGGREGATOR)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var hosts []string
	if err := dec.Decode(&hosts); err != nil {
		return "", err
	}

	return hosts[rand.Intn(len(hosts)-1)], nil
}

func GetDataFromSrc(host string) (*HostResp, error) {
	resp, err := DefClient.Get(host)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	hostResp := new(HostResp)
	if err := dec.Decode(hostResp); err != nil {
		return nil, err
	}

	return hostResp, err
}
