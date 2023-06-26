package crawler

import (
	"errors"
	"net/url"
	"sort"
)

// CrawlerInterface defines the interface for a web crawler.
type CrawlerInterface interface {
	Run()
	GetResult() []string
	SortLinks()
}

// Strategy represents a web crawling strategy.
type Strategy interface {
	Run(url *url.URL) []string
}

// Crawler represents a web crawler.
type Crawler struct {
	Strategy Strategy // Algorithm to use
	Result   []string // Result of the crawl
}

// NewCrawler creates a new web crawler with the specified strategy.
func NewCrawler(strategy string) (*Crawler, error) {
	// Create the strategy
	alg, err := createStrategy(strategy)
	if err != nil {
		return nil, err
	}
	// Create the crawler
	crawler := &Crawler{
		Strategy: alg,
		Result:   []string{},
	}
	return crawler, nil
}

// Run starts the web crawling process with the specified root URL.
// It returns the result of the crawl or an error if any occurred.
func (c *Crawler) Run(rootUrl string) ([]string, error) {
	// Parse the given URL
	parsedURL, err := url.Parse(rootUrl)
	if err != nil {
		return []string{}, errors.New("error parsing URL")
	}
	// Run the algorithm
	c.Result = c.Strategy.Run(parsedURL)
	return c.Result, nil
}

// GetResult returns the result of the web crawl.
func (c *Crawler) GetResult() []string {
	return c.Result
}

// SortLinks sorts the crawled links in ascending order.
func (c *Crawler) SortLinks() {
	sort.Strings(c.Result)
}

// createStrategy creates a web crawling strategy based on the provided string.
func createStrategy(strategy string) (Strategy, error) {
	switch strategy {
	case "OneLevel":
		return NewOneLevel(), nil
	case "Recursive":
		return NewRecursive(), nil
	default:
		return nil, errors.New("error creating strategy")
	}
}
