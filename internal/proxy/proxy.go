package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Proxy struct {
	IP   string
	Port string
}

func scrapeProxiesFromFreeProxyList() ([]Proxy, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Fetch the webpage
	resp, err := client.Get("https://free-proxy-list.net/")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch webpage: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch webpage: status %d", resp.StatusCode)
	}

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var proxies []Proxy

	// Select the proxy table rows
	doc.Find(".fpl-list .table tbody tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		if tds.Length() >= 8 {
			proxy := Proxy{
				IP:   strings.TrimSpace(tds.Eq(0).Text()),
				Port: strings.TrimSpace(tds.Eq(1).Text()),
			}
			proxies = append(proxies, proxy)
		}
	})

	return proxies, nil
}

// Test if a proxy works by trying to access reddit.com
func testProxy(proxy Proxy) bool {
	proxyURL := fmt.Sprintf("http://%s:%s", proxy.IP, proxy.Port)

	parsedProxy, err := url.Parse(proxyURL)
	if err != nil {
		return false
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(parsedProxy),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // Shorter timeout for testing
	}

	resp, err := client.Get("https://reddit.com")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Consider proxy working if we get any HTTP response (even redirects)
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

// Result struct to hold proxy test results
type ProxyResult struct {
	Proxy   Proxy
	Working bool
	Index   int
}

// Alternative: use sync.Once for safe closing
func findWorkingProxyFastest(proxies []Proxy) (*Proxy, error) {
	fmt.Printf("Testing %d proxies concurrently (returning fastest)...\n", len(proxies))

	resultChan := make(chan *Proxy, 1)
	doneChan := make(chan struct{})

	var wg sync.WaitGroup
	var closeOnce sync.Once // Ensures doneChan is only closed once

	maxConcurrent := 50
	semaphore := make(chan struct{}, maxConcurrent)

	for _, proxy := range proxies {
		wg.Add(1)
		go func(p Proxy) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			select {
			case <-doneChan:
				return
			default:
				if testProxy(p) {
					if resultChan == nil {
						fmt.Printf("Found working proxy: %s:%s\n", p.IP, p.Port)

					}
					select {
					case resultChan <- &p:
						closeOnce.Do(func() { close(doneChan) }) // Safe close
					case <-doneChan:
					}
				}
			}
		}(proxy)
	}

	go func() {
		wg.Wait()
		closeOnce.Do(func() { close(resultChan) }) // Safe close
	}()

	select {
	case proxy := <-resultChan:
		if proxy != nil {
			return proxy, nil
		}
		return nil, fmt.Errorf("no working proxies found")
	case <-time.After(60 * time.Second):
		closeOnce.Do(func() { close(doneChan) }) // Safe close
		return nil, fmt.Errorf("timeout: no working proxies found within 60 seconds")
	}
}

// Create HTTP client with first working proxy from scraped list
func CreateSpysProxyClient() (*http.Client, *Proxy, error) {
	proxies, err := scrapeProxiesFromFreeProxyList()
	if err != nil {
		return nil, nil, err
	}
	if len(proxies) == 0 {
		return nil, nil, fmt.Errorf("no proxies available")
	}

	// Find first working proxy (you can use either version)
	// selectedProxy, err := findWorkingProxy(proxies) // Returns first in list that works
	selectedProxy, err := findWorkingProxyFastest(proxies) // Returns fastest responding proxy

	if err != nil {
		return nil, nil, err
	}

	proxyURL := fmt.Sprintf("http://%s:%s", selectedProxy.IP, selectedProxy.Port)
	fmt.Printf("Using proxy: %s\n", proxyURL)

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return client, selectedProxy, nil
}
