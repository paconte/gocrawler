package crawler_test

import (
	"sort"
	"testing"

	crawler "github.com/paconte/gocrawler/crawler"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNewCrawler(t *testing.T) {

	tests := []struct {
		strategy string
		fails    bool
	}{
		{"OneLevel", false},
		{"Recursive", false},
		{"Something", true},
		{"", true},
	}

	for _, tt := range tests {
		_, err := crawler.NewCrawler(tt.strategy)
		if tt.fails {
			assert.NotNil(t, err)
		}
	}
}

func TestRun(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Load the stored response from the file
	filePath := "testdata/response.html"
	fileContent := LoadFileAsString(t, filePath)

	// Mock http response
	httpmock.RegisterResponder("GET", "https://parserdigital.com/",
		httpmock.NewStringResponder(200, fileContent))

	// Test Run
	c, _ := crawler.NewCrawler("OneLevel")
	tests := []struct {
		url   string
		fails bool
	}{
		{"https://parserdigital.com/", false},
		{" http://foo.com", true},
		{"1http://foo.com", true},
		{"cache_object:foo/bar", false},
	}
	for _, tt := range tests {
		_, err := c.Run(tt.url)
		if tt.fails {
			assert.NotNil(t, err)
		}
	}
}

func TestSort(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Load the stored response from the file
	filePath := "testdata/response.html"
	fileContent := LoadFileAsString(t, filePath)

	// Mock http response
	httpmock.RegisterResponder("GET", "https://parserdigital.com/",
		httpmock.NewStringResponder(200, fileContent))

	// Test sort links
	c, _ := crawler.NewCrawler("OneLevel")
	resultsA, _ := c.Run("https://parserdigital.com/")
	sort.Strings(resultsA)
	c.SortLinks()
	resultsB := c.GetResult()
	assert.Equal(t, len(resultsA), len(resultsB))
	for i, link := range resultsA {
		assert.Equal(t, link, resultsB[i])
	}
}
