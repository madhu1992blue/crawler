package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}
	if cfg.baseURL.Host != currentURL.Host {
		return
	}
	rawCurrentURLNormalized, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	_, ok := cfg.pages[rawCurrentURLNormalized]
	if ok {
		cfg.pages[rawCurrentURLNormalized] += 1
		cfg.mu.Unlock()
		return
	} else {
		cfg.pages[rawCurrentURLNormalized] = 1
	}
	cfg.mu.Unlock()

	htmlString, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}
	fmt.Println(htmlString)
	raw_urls, err := getURLsFromHTML(htmlString, cfg.baseURL.String())
	if err != nil {
		return
	}
	for _, ru := range raw_urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(ru)
	}
}
