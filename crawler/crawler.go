package crawler

import (
	"errors"
	"net/url"
	"sort"
)

// Strategy represents a web crawling strategy.
type Strategy interface {
	Run(url *url.URL) []string
}

// Run starts the web crawling process with the specified root URL.
// It returns the result of the crawl or an error if any occurred.
func Run(rootUrl string, strategy string) ([]string, error) {
	// Create the strategy
	st, err := createStrategy(strategy)
	if err != nil {
		return nil, err
	}
	// Parse the given URL
	parsedURL, err := url.Parse(rootUrl)
	if err != nil {
		return nil, err
	}
	// Run the algorithm
	result := st.Run(parsedURL)
	sort.Strings(result)
	return result, nil
}

// createStrategy creates a web crawling strategy based on the provided string.
func createStrategy(strategy string) (Strategy, error) {
	switch strategy {
	case "OneLevel":
		return NewOneLevel(), nil
	case "Recursive":
		return NewRecursive(), nil
	case "RecursiveParallel":
		return NewRecursiveParallel(), nil
	default:
		return nil, errors.New("error creating strategy")
	}
}
