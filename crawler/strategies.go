package crawler

import (
	"net/url"
)

// This strategy implements asearch approach to crawl URLs recursively,
// discovering new URLs at each level and continuing the crawling process until
// there are no more unvisited URLs.
type Recursive struct {
	visited map[string]bool // Visited URLs
	found   map[string]bool // Found URLs
}

// NewRecursive creates a new instance of the Recursive strategy.
func NewRecursive() *Recursive {
	strategy := &Recursive{
		visited: map[string]bool{},
	}
	return strategy
}

// Run starts the web crawling process using the Recursive strategy.
// It takes the root URL as input and returns a list of visited URLs.
func (s *Recursive) Run(rootUrl *url.URL) []string {
	s.found = map[string]bool{rootUrl.String(): true}
	for len(s.visited) != len(s.found) {
		for link := range s.found {
			if s.visited[link] {
				continue
			}
			newFounds := CollectMap(Extract(Parse(Download(link)), rootUrl))
			for newFound := range newFounds {
				s.found[newFound] = true
			}
			s.visited[link] = true
		}
	}
	return MapToList(s.visited)
}

// This strategy crawls the root URL and collects URLs up to one level deep.
// It returns a list of collected URLs.
type OneLevel struct{}

// NewOneLevel creates a new instance of the OneLevel strategy.
func NewOneLevel() *OneLevel {
	strategy := &OneLevel{}
	return strategy
}

// Run starts the web crawling process using the OneLevel strategy.
// It takes the root URL as input and returns a list of collected URLs.
func (s *OneLevel) Run(rootUrl *url.URL) []string {
	return MapToList(CollectMap(Extract(Parse(Download(rootUrl.String())), rootUrl)))
}
