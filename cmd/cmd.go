package cmd

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func startOsSignalHandler(terminate func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		terminate()
		log.Infoln("Готуюся до закриття.")
	}()
}

var (
	botsNum     int
	reqPerEpoch int
	onlyProxy   bool
	maxErrCount int
	srcFile     string
)

func Execute() {
	var rootCmd = &cobra.Command{Use: "antirus"}
	rootCmd.PersistentFlags().IntVar(&botsNum, "bots", 200, "кількість ботів (активних з'єднань)")
	rootCmd.PersistentFlags().IntVar(&reqPerEpoch, "epoch", 10_000, "к-сть запитів перед новою ціллю")
	rootCmd.PersistentFlags().BoolVar(&onlyProxy, "onlyproxy", false, "з'єднання тільки через проксі")
	rootCmd.PersistentFlags().IntVar(&maxErrCount, "errcount", 100, "к-сть помилок на бота, щоб той закінчив роботу")
	rootCmd.PersistentFlags().StringVar(&srcFile, "file", "", "файл із цілями та проксі")

	rootCmd.AddCommand(ddatackCmd)
	rootCmd.Execute()
}
