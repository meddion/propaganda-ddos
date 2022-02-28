package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	done := make(chan struct{}, 1)
	var once sync.Once
	terminate := func() {
		once.Do(func() {
			close(done)
		})
	}
	// Os signal handler
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		terminate()
		log.Infoln("Готуюся до закриття.")
	}()

	for epoch := 0; true; epoch++ {
		select {
		case <-done:
			return
		default:
		}

		ctx, termBots := context.WithCancel(context.Background())
		epochDone := make(chan struct{}, 1)
		go func() {
			select {
			case <-epochDone:
			case <-done:
				termBots()
			}
		}()

		log.Infof("Готуюся до атаки русні :) Сесія: %d \n", epoch)
		log.Infof("Кожні %d запитів ціль і проксі можуть змінюватися.\n", reqPerEpoch)
		log.Infof("Ціль: %s\n", args[0])

		targetResp := pkg.NewTargetData(args[0])

		var wg sync.WaitGroup
		counter := make(chan bool, botsNum)

		if err := pkg.StartBots(ctx, &wg, targetResp, maxErrCount, botsNum, counter); err != nil {
			log.Errorf("Не вдалося запустити ботів: %v\n", err)
			continue
		}

		go pkg.RequestLimiter(ctx, termBots, &wg, 0, counter)

		wg.Wait()
		close(epochDone)
	}
}
