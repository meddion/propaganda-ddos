package pkg

import (
	"context"
	"sync"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestBot(t *testing.T) {
	ctx, termBots := context.WithCancel(context.Background())
	// done := make(chan struct{}, 1)
	counter := make(chan bool, 1)

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		c, err := NewBot(ctx, 1, nil, false)
		if err != nil {
			t.Fatal(err)
		}

		go func() {
			defer wg.Done()
			c.Start(ctx, "google.com", 10, counter)
		}()
	}

	time.AfterFunc(time.Second, func() {
		log.Println("canceling the bot")
		termBots()
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
