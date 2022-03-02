package pkg

import (
	"context"
	"sync"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestBot(t *testing.T) {
	botsNum := 2
	maxErrCount := 10
	onlyProxy := false
	targets := []Target{{URL: "https://gggdfgfgdfg.sw"}, {URL: "https://asdfasda.sad"}}
	proxies := []Proxy{}

	rootCtx, term := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	for _, target := range targets {
		if err := ValidateTarget(&target); err != nil {
			log.Errorf("Під час валідації даних про атаку (перевірте джерело): %v", err)
			continue
		}

		log.Infof("Дані із джерела підвантажено. Ціль: %s, К-ість проксі: %d\n",
			target.URL, len(proxies))

		botSched := NewBotScheduler(target, proxies, botsNum, maxErrCount, onlyProxy)
		if err := botSched.StartBots(rootCtx, &wg); err != nil {
			log.Errorf("Не вдалося запустити ботів: %v\n", err)
			continue
		}
	}

	time.AfterFunc(time.Second*2, func() {
		term()
	})

	wg.Wait()
}

// func TestValidateURL(t *testing.T) {
// 	assert.NoError(t, ValidateURL("http://google.com"))
// 	assert.NoError(t, ValidateURL("http://w.com/cn"))
// 	assert.NoError(t, ValidateURL("http://192.158.0.1:90"))
// 	assert.Error(t, ValidateURL("http://w"))
// 	assert.Error(t, ValidateURL("fsw"))
// 	assert.Error(t, ValidateURL("http://192.158.1/1"))
// }
