package pkg

// func TestBot(t *testing.T) {
// 	botsNum := 2
// 	maxErrCount := 10
// 	onlyProxy := false
// 	targets := []Target{{URL: "https://gggdfgfgdfg.sw"}, {URL: "https://asdfasda.sad"}}
// 	proxies := []Proxy{}

// 	rootCtx, term := context.WithCancel(context.Background())

// 	var wg sync.WaitGroup
// 	for _, target := range targets {
// 		if err := ValidateTarget(&target); err != nil {
// 			log.Errorf("Під час валідації даних про атаку (перевірте джерело): %v", err)
// 			continue
// 		}

// 		log.Infof("Дані із джерела підвантажено. Ціль: %s, К-ість проксі: %d\n",
// 			target.URL, len(proxies))

// 		botSched := NewBotScheduler(target, proxies, botsNum, maxErrCount, onlyProxy)
// 		if err := botSched.Start(rootCtx, &wg); err != nil {
// 			log.Errorf("Не вдалося запустити ботів: %v\n", err)
// 			continue
// 		}
// 	}

// 	time.AfterFunc(time.Second*2, func() {
// 		term()
// 	})

// 	wg.Wait()
// }

// func TestClient(t *testing.T) {

// }

// func TestValidateAddress(t *testing.T) {
// 	addrs := []struct{
// 			in string
// 			out string
// 			resolve bool
// 		}{
// 		{"google.com:80", "http://142.250.203.206:80"},
// 		{"https://google.com/dfasfs/q=2323",
// 		{"http://google.com/q=hello", "http://google.com/q=hello"},

// 		{"142.250.217.78:443", "http://142.250.217.78:443",
// 		{"http://localhost:8080", "http://127.0.0.1:8080"},
// 		{"localhost", "http://localhost"},
// 		{"http://127.0.0.1:8080", "http://127.0.0.1:8080"},
// 		{"http://142.250.217.78", "http://142.250.217.78"},
// 		{"https://142.250.217.78:80", "https://142.250.217.78:80"},

// 	}

// 	// if strings.Contains()

// 	for _, a := range addrs {
// 		var err error
// 		a, err = ValidateAddress(context.TODO(), a, true)
// 		if err != nil {
// 			t.Error(err)
// 			continue
// 		}

// 		log.Println(a)

// 		if _, err := DefClient.Get(a); err != nil {
// 			t.Errorf("[ERR] %s: %v", a, err)
// 		} else {
// 			log.Errorf("[200] %s", a)
// 		}
// 	}
// }
