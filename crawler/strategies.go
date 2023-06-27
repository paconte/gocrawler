package crawler

import (
	"net/url"
	"sync"
)

/*
 * ######### RECURSIVE ###########
 */

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

/*
 * ######### RECURSIVE PARALLEL ###########
 */

// RecursiveParallel implements a parallelized version of the Recursive strategy.
type RecursiveParallel struct {
	visited map[string]bool // Visited URLs
	found   map[string]bool // Found URLs
	mutex   sync.Mutex
}

// NewRecursiveParallel creates a new instance of the RecursiveParallel strategy.
func NewRecursiveParallel() *RecursiveParallel {
	strategy := &RecursiveParallel{
		visited: map[string]bool{},
		found:   map[string]bool{},
		mutex:   sync.Mutex{},
	}
	return strategy
}

// Run starts the web crawling process using the RecursiveParallel strategy.
// It takes the root URL as input and returns a list of visited URLs.
func (s *RecursiveParallel) Run(rootUrl *url.URL) []string {
	s.found = map[string]bool{rootUrl.String(): true}
	var wg sync.WaitGroup
	for len(s.visited) != len(s.found) {
		for link := range s.found {
			if s.isVisited(link) {
				continue
			}
			wg.Add(1)
			go s.job(link, rootUrl, &wg)
		}
		wg.Wait()
	}
	return MapToList(s.visited)
}

// isVisited checks if a URL has already been visited.
func (s *RecursiveParallel) isVisited(link string) bool {
	visited := false
	s.mutex.Lock()
	if s.visited[link] {
		visited = true
	}
	s.mutex.Unlock()
	return visited
}

// job performs the crawling job for a specific URL.
func (s *RecursiveParallel) job(link string, rootUrl *url.URL, wg *sync.WaitGroup) {
	newFounds := CollectMap(Extract(Parse(Download(link)), rootUrl))

	s.mutex.Lock()
	for newFound := range newFounds {
		s.found[newFound] = true
	}
	s.visited[link] = true
	s.mutex.Unlock()

	wg.Done()
}

/*
 * ######### ONELEVEL ###########
 */

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
