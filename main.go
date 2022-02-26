package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const HOST_AGGREGATOR = "http://rockstarbloggers.ru/hosts.json"

func GetHost() (string, error) {
	resp, err := http.Get(HOST_AGGREGATOR)
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

type HostResp struct {
	Site struct {
		ID           int    `json:"id"`
		URL          string `json:"url"`
		NeedParseURL int    `json:"need_parse_url"`
		PageTime     string `json:"page_time"`
		Attack       int    `json:"atack"`
	}

	Proxy []struct {
		ID   int    `json:"id"`
		IP   string `json:"ip"`
		Auth string `json:"auth"`
	}
}

func GetDataFromHost(host string) (*HostResp, error) {
	resp, err := http.Get(host)
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

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		host, err := GetHost()
		if err != nil {
			log.Printf("on getting hosts: %v\n", err)
			time.Sleep(time.Second * 5)
			continue
		}
		log.Printf("Using the host: %s\n", host)

		data, err := GetDataFromHost(host)

		if err != nil {
			log.Printf("on parsing data from host: %v\n", err)
			time.Sleep(time.Second * 5)
			continue
		}

		log.Println(data)
	}
}
