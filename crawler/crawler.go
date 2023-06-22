package crawler

import (
	"errors"
	"fmt"
	"net/url"
	"sort"

	"golang.org/x/sync/syncmap"
)

// Crawler represents a web crawler.
type Crawler struct {
	BaseURL  *url.URL    // Base URL of the crawler
	Strategy Strategy    // Algorithm to use
	Result   []string    // Result of the crawl
	Visited  syncmap.Map // Visited URLs
}

// NewCrawler creates a new web crawler with the specified base URL and maximum visits.
func NewCrawler(baseURL string, strategy string) (*Crawler, error) {
	var alg Strategy
	parsedURL, err := url.Parse(baseURL)

	if err != nil {
		return nil, errors.New("error parsing URL")
	}

	if strategy == "bruteforce" {
		alg, err = NewBruteForce(parsedURL)
		if err != nil {
			return nil, errors.New("error creating bruteforce strategy")
		}
	} else if strategy == "singlevisit" {
		alg, err = NewSingleVisit(parsedURL)
		if err != nil {
			return nil, errors.New("error creating singlevisit strategy")
		}
	} else {
		return nil, errors.New("error creating strategy")
	}

	crawler := &Crawler{
		BaseURL:  parsedURL,
		Strategy: alg,
		Visited:  syncmap.Map{},
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
