package crawler

import (
	"context"
	"net/url"
	"sync"
	"time"
)

// Limits represents the limits of a strategy.
// It contains the maximum number of seconds and requests.
type Limits struct {
	Milliseconds int
	Requests     int
}

/*
 * ######### RECURSIVE ###########
 */

// This strategy implements asearch approach to crawl URLs recursively,
// discovering new URLs at each level and continuing the crawling process until
// there are no more unvisited URLs.
type Recursive struct {
	visited map[string]bool // Visited URLs
	found   map[string]bool // Found URLs
	url     *url.URL        // Root URL
}

// NewRecursive creates a new instance of the Recursive strategy.
func NewRecursive(url *url.URL) *Recursive {
	strategy := &Recursive{
		url:     url,
		found:   map[string]bool{url.String(): true},
		visited: map[string]bool{},
	}
	return strategy
}

// Run starts the web crawling process using the Recursive strategy.
// It returns a list of visited URLs.
func (s *Recursive) Run() []string {
	for len(s.visited) != len(s.found) {
		for link := range s.found {
			if s.visited[link] {
				continue
			}
			newFounds := CollectMap(Extract(Parse(Download(link)), s.url))
			for newFound := range newFounds {
				s.found[newFound] = true
			}
			s.visited[link] = true
		}
	}
	return MapToList(s.visited)
}

/*
 * ######### RECURSIVE WITH LIMITS ###########
 */

// RecursiveWithLimits implements the same strategy as the Recursive strategy,
// but adding limits to the number of requests and time.
type RecursiveWithLimits struct {
	visited map[string]bool // Visited URLs
	found   map[string]bool // Found URLs
	url     *url.URL        // Root URL
	limits  Limits
}

// NewRecursiveWithLimits creates a new instance of the Recursive strategy.
func NewRecursiveWithLimits(url *url.URL, limits Limits) *RecursiveWithLimits {
	strategy := &RecursiveWithLimits{
		visited: map[string]bool{},
		found:   map[string]bool{url.String(): true},
		url:     url,
		limits:  limits,
	}
	return strategy
}

// Run starts the web crawling process using the Recursive strategy with limits.
// It returns a list of visited URLs.
func (s *RecursiveWithLimits) Run() []string {
	isTimeout := false
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Duration(s.limits.Milliseconds)*time.Millisecond)
	defer cancel()

	for len(s.visited) != len(s.found) {
		// Stop if the number of requests exceeds the limit
		if len(s.visited) >= s.limits.Requests {
			break
		}
		// Stop if the timeout is exceeded
		select {
		case <-ctx.Done():
			isTimeout = true
		default:
			for link := range s.found {
				if s.visited[link] {
					continue
				}
				newFounds := CollectMap(Extract(Parse(Download(link)), s.url))
				for newFound := range newFounds {
					s.found[newFound] = true
				}
				s.visited[link] = true
			}
		}
		if isTimeout {
			break
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
	url     *url.URL        // Root URL
	mutex   sync.Mutex
}

// NewRecursiveParallel creates a new instance of the RecursiveParallel strategy.
func NewRecursiveParallel(url *url.URL) *RecursiveParallel {
	strategy := &RecursiveParallel{
		visited: map[string]bool{},
		found:   map[string]bool{url.String(): true},
		url:     url,
		mutex:   sync.Mutex{},
	}
	return strategy
}

// Run starts the web crawling process using the RecursiveParallel strategy.
// It takes the root URL as input and returns a list of visited URLs.
func (s *RecursiveParallel) Run() []string {
	var wg sync.WaitGroup
	for len(s.visited) != len(s.found) {
		for link := range s.found {
			if s.isVisited(link) {
				continue
			}
			wg.Add(1)
			go s.job(link, s.url, &wg)
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
 * ######### RECURSIVE PARALLEL WITH LIMITS ###########
 */

// RecursiveParallelWithLimits implements a parallelized version of the Recursive strategy
// It also has a limit of http requests and time.
type RecursiveParallelWithLimits struct {
	visited map[string]bool // Visited URLs
	found   map[string]bool // Found URLs
	url     *url.URL        // Root URL
	mutex   sync.Mutex
	limits  Limits
}

// RecursiveParallelWithLimits creates a new instance of the RecursiveParallelWithLimits strategy.
func NewRecursiveParallelWithLimits(url *url.URL, limits Limits) *RecursiveParallelWithLimits {
	strategy := &RecursiveParallelWithLimits{
		visited: map[string]bool{},
		found:   map[string]bool{url.String(): true},
		url:     url,
		mutex:   sync.Mutex{},
		limits:  limits,
	}
	return strategy
}

// Run starts the web crawling process using the RecursiveParallelWithLimits strategy.
// It takes the root URL as input and returns a list of visited URLs.
func (s *RecursiveParallelWithLimits) Run() []string {
	var wg sync.WaitGroup
	isTimeout := false

	ctx, cancel := context.WithTimeout(
		context.Background(), time.Duration(s.limits.Milliseconds)*time.Millisecond)
	defer cancel()

	for len(s.visited) != len(s.found) {
		// Stop if the number of requests exceeds the limit
		if len(s.visited) >= s.limits.Requests {
			break
		}
		// Stop if the timeout is exceeded
		select {
		case <-ctx.Done():
			isTimeout = true
		default:
			for link := range s.found {
				if s.isVisited(link) {
					continue
				}
				wg.Add(1)
				go s.job(link, s.url, &wg)
			}
			wg.Wait()
		}
		if isTimeout {
			break
		}
	}
	return MapToList(s.visited)
}

// isVisited checks if a URL has already been visited.
func (s *RecursiveParallelWithLimits) isVisited(link string) bool {
	visited := false
	s.mutex.Lock()
	if s.visited[link] {
		visited = true
	}
	s.mutex.Unlock()
	return visited
}

// job performs the crawling job for a specific URL.
func (s *RecursiveParallelWithLimits) job(link string, rootUrl *url.URL, wg *sync.WaitGroup) {
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
type OneLevel struct {
	url *url.URL
}

// NewOneLevel creates a new instance of the OneLevel strategy.
func NewOneLevel(url *url.URL) *OneLevel {
	strategy := &OneLevel{url: url}
	return strategy
}

// Run starts the web crawling process using the OneLevel strategy.
// It takes the root URL as input and returns a list of collected URLs.
func (s *OneLevel) Run() []string {
	return MapToList(CollectMap(Extract(Parse(Download(s.url.String())), s.url)))
}
