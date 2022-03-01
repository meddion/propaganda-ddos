package cmd

import (
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var once sync.Once

func CreateDoneChan() (<-chan struct{}, func()) {
	done := make(chan struct{}, 1)
	return done, func() {
		once.Do(func() {
			close(done)
		})
	}
}
func startOsSignalHandler(terminate func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		terminate()
		log.Infoln("Готуюся до закриття.")
	}()
}

func startBotHandler(done <-chan struct{}, termBots func(), epoch *int64, curEpoch int64) {
	go func() {
		tk := time.NewTicker(time.Minute)
		select {
		case <-tk.C:
			if atomic.LoadInt64(epoch) != curEpoch {
				return
			}
		case <-done:
			termBots()
		}
	}()
}

var (
	botsNum     int
	reqPerEpoch int
	onlyProxy   bool
	maxErrCount int
)

func Execute() {
	var rootCmd = &cobra.Command{Use: "antirus"}
	rootCmd.PersistentFlags().IntVar(&botsNum, "bots", 200, "кількість ботів (активних з'єднань)")
	rootCmd.PersistentFlags().IntVar(&reqPerEpoch, "epoch", 10_000, "к-сть запитів перед новою ціллю")
	rootCmd.PersistentFlags().BoolVar(&onlyProxy, "onlyproxy", false, "з'єднання тільки через проксі")
	rootCmd.PersistentFlags().IntVar(&maxErrCount, "errcount", 100, "к-сть помилок на бота, щоб той закінчив роботу")

	rootCmd.AddCommand(ddatackCmd)
	rootCmd.AddCommand(targetCmd)
	rootCmd.Execute()
}
