package crawler

import (
	"errors"
	"net/url"
	"sort"
)

// Strategy represents a web crawling strategy.
type Strategy interface {
	Run() []string
}

// Run starts the web crawling process with the specified root URL.
// It returns the result of the crawl or an error if any occurred.
func Run(rootUrl string, strategy string, limits Limits) ([]string, error) {
	// Parse the given URL
	parsedURL, err := url.Parse(rootUrl)
	if err != nil {
		return nil, err
	}
	// Create the strategy
	st, err := createStrategy(parsedURL, strategy, limits)
	if err != nil {
		return nil, err
	}
	// Run the algorithm
	result := st.Run()
	sort.Strings(result)
	return result, nil
}

// createStrategy creates a web crawling strategy based on the provided string.
func createStrategy(url *url.URL, strategy string, limits Limits) (Strategy, error) {
	switch strategy {
	case "OneLevel":
		return NewOneLevel(url), nil
	case "Recursive":
		return NewRecursive(url), nil
	case "RecursiveParallel":
		return NewRecursiveParallel(url), nil
	case "RecursiveWithLimits":
		return NewRecursiveWithLimits(url, limits), nil
	case "RecursiveParallelWithLimits":
		return NewRecursiveParallelWithLimits(url, limits), nil
	default:
		return nil, errors.New("error creating strategy")
	}
}
