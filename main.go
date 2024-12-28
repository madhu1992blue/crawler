package main

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	cmdArgs := os.Args[1:]
	if len(cmdArgs) < 3 {
		fmt.Println("usage: ./crawler concurrency maxpages")
		os.Exit(1)
	}
	if len(cmdArgs) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	BASE_URL := cmdArgs[0]
	concurrencyLimit, err := strconv.Atoi(cmdArgs[1])
	if err != nil {
		fmt.Printf("concurrency should be a number")
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(cmdArgs[2])
	if err != nil {
		fmt.Printf("maxPages should be a number")
		os.Exit(1)
	}
	fmt.Printf("starting crawl of: %s\n", BASE_URL)

	var pages map[string]int = make(map[string]int)

	baseURL, err := url.Parse(BASE_URL)
	if err != nil {
		fmt.Println("Something wrong with BASE_URL")
		os.Exit(1)
	}
	cfg := config{
		pages:              pages,
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, concurrencyLimit),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}
	cfg.wg.Add(1)
	go cfg.crawlPage(BASE_URL)
	cfg.wg.Wait()
	printReport(pages, BASE_URL)
}

func printReport(pages map[string]int, baseURL string) {
	type urlCount struct {
		url   string
		count int
	}
	var urlCounts []urlCount
	for k, v := range pages {
		urlCounts = append(urlCounts, urlCount{url: k, count: v})
	}
	sort.Slice(urlCounts, func(i, j int) bool {
		return urlCounts[i].count >= urlCounts[j].count
	})
	fmt.Printf(`
=============================
  REPORT for %s
=============================
`, baseURL)
	for _, uc := range urlCounts {
		fmt.Printf("Found %d internal links to %s\n", uc.count, uc.url)
	}

}
