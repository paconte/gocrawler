package crawler

import (
	"net/url"
	"sync"

	"golang.org/x/sync/syncmap"
)

type Strategy interface {
	Run(rootUrl *url.URL) []string
}

type BruteForce struct {
}

type SingleVisit struct {
	Visited syncmap.Map // Visited URLs
}

func NewBruteForce(baseURL *url.URL) (*BruteForce, error) {
	strategy := &BruteForce{}
	return strategy, nil
}

func NewSingleVisit(baseURL *url.URL) (*SingleVisit, error) {
	strategy := &SingleVisit{
		Visited: sync.Map{},
	}
	return strategy, nil
}

func (b BruteForce) Run(rootUrl *url.URL) []string {
	var result []string = []string{rootUrl.String()}
	for url := range Extract(Parse(Download(rootUrl.String())), rootUrl) {
		result = append(result, url)
		for url2 := range Extract(Parse(Download(url)), rootUrl) {
			result = append(result, url2)
		}
	}
	return result
}

func (s *SingleVisit) Run(rootUrl *url.URL) []string {
	s.Visited.Store(rootUrl.String(), true)
	for level1Urls := range Extract(Parse(Download(rootUrl.String())), rootUrl) {
		s.Visited.Store(level1Urls, true)
		for level2Urls := range Extract(Parse(Download(level1Urls)), rootUrl) {
			if _, stored := s.Visited.LoadOrStore(level2Urls, true); !stored { // If not visited
				s.Visited.Store(level2Urls, true)
			}
		}
	}
	return s.mapToSlice()
}

func (s *SingleVisit) mapToSlice() []string {
	var result []string
	s.Visited.Range(func(key, value interface{}) bool {
		result = append(result, key.(string))
		return true
	})
	return result
}
