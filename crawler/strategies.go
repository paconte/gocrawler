package crawler

import "net/url"

// This strategy implements asearch approach to crawl URLs recursively,
// discovering new URLs at each level and continuing the crawling process until
// there are no more unvisited URLs.
type Recursive struct {
	visited map[string]bool // Visited URLs
	found   map[string]bool // Found URLs
}

func NewRecursive(baseURL *url.URL) *Recursive {
	strategy := &Recursive{
		visited: map[string]bool{},
		found:   map[string]bool{baseURL.String(): true},
	}
	return strategy
}

func (s *Recursive) Run(rootUrl *url.URL) []string {
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

func NewOneLevel(baseURL *url.URL) *OneLevel {
	strategy := &OneLevel{}
	return strategy
}

func (s *OneLevel) Run(rootUrl *url.URL) []string {
	return MapToList(CollectMap(Extract(Parse(Download(rootUrl.String())), rootUrl)))
}
