package crawler

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
)

type CrawlerInterface interface {
	Run()
	SortLinks()
	PrintLinks()
}

// Strategy represents a web crawling strategy.
type Strategy interface {
	Run(url *url.URL) []string
}

// Crawler represents a web crawler.
type Crawler struct {
	BaseURL  *url.URL // Base URL of the crawler
	Strategy Strategy // Algorithm to use
	Result   []string // Result of the crawl
}

// NewCrawler creates a new web crawler with the specified base URL and maximum visits.
func NewCrawler(baseURL string, strategy string) (*Crawler, error) {
	// Parse the given URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, errors.New("error parsing URL")
	}
	// Create the strategy
	alg, err := createStrategy(strategy, parsedURL)
	if err != nil {
		return nil, err
	}
	// Create the crawler
	crawler := &Crawler{
		BaseURL:  parsedURL,
		Strategy: alg,
	}

	return crawler, nil
}

func (c *Crawler) Run() {
	c.Result = c.Strategy.Run(c.BaseURL)
}

func (c *Crawler) SortLinks() {
	sort.Strings(c.Result)
}

func (c *Crawler) PrintLinks() {
	for _, link := range c.Result {
		fmt.Println(link)
	}
}

func createStrategy(strategy string, parsedURL *url.URL) (Strategy, error) {
	switch strategy {
	case "OneLevel":
		return NewOneLevel(parsedURL), nil
	case "Recursive":
		return NewRecursive(parsedURL), nil
	default:
		return nil, errors.New("error creating strategy")
	}
}
