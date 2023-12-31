package crawler_test

import (
	"net/url"
	"testing"

	"github.com/paconte/gocrawler/crawler"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var HtmlFiles = map[string]string{
	"http://www.parserdigital.com":   "testdata/treeLevel1.html",
	"http://www.parserdigital.com/A": "testdata/treeLevel2A.html",
	"http://www.parserdigital.com/B": "testdata/treeLevel2B.html",
	"http://www.parserdigital.com/C": "testdata/treeLevel3C.html",
	"http://www.parserdigital.com/D": "testdata/treeLevel3D.html",
	"http://www.parserdigital.com/E": "testdata/treeLevel3E.html",
	"http://www.parserdigital.com/F": "testdata/treeLevel3F.html",
}

func TestOneLevel(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Load the stored response from the file
	filePath := "testdata/response.html"
	fileContent := LoadFileAsString(t, filePath)

	// Define test cases
	tests := []struct {
		url      string
		status   int
		content  string
		expected int
	}{
		{"https://parserdigital.com/", 200, fileContent, 23},
		{"https://google.com", 200, "", 0},
		{"A", 400, "", 0},
		{"B", 200, "", 0},
	}

	for _, tt := range tests {

		// Mock http response
		httpmock.RegisterResponder("GET", tt.url,
			httpmock.NewStringResponder(tt.status, tt.content))

		// Test the function
		parsedUrl, _ := url.Parse(tt.url)
		strategy := crawler.NewOneLevel(parsedUrl)
		result := strategy.Run()
		assert.Equal(t, tt.expected, len(result))
	}
}

func TestRecursive(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responses
	for domain, file := range HtmlFiles {
		// Load the stored response from the file
		fileContent := LoadFileAsString(t, file)

		// Mock http response
		httpmock.RegisterResponder("GET", domain,
			httpmock.NewStringResponder(200, fileContent))
	}

	// Test the function
	parsedUrl, _ := url.Parse("http://www.parserdigital.com")
	strategy := crawler.NewRecursive(parsedUrl)
	result := strategy.Run()
	info := httpmock.GetCallCountInfo()

	assert.Equal(t, 7, len(result))
	for link := range HtmlFiles {
		assert.Equal(t, 1, info["GET "+link])
	}
}

func TestRecursiveParallel(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responses
	for domain, file := range HtmlFiles {
		// Load the stored response from the file
		fileContent := LoadFileAsString(t, file)

		// Mock http response
		httpmock.RegisterResponder("GET", domain,
			httpmock.NewStringResponder(200, fileContent))
	}

	// Test the function
	parsedUrl, _ := url.Parse("http://www.parserdigital.com")
	strategy := crawler.NewRecursiveParallel(parsedUrl)
	result := strategy.Run()
	info := httpmock.GetCallCountInfo()

	assert.Equal(t, 7, len(result))
	for link := range HtmlFiles {
		assert.Equal(t, 1, info["GET "+link])
	}
}

func TestRecursiveWithLimits(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responses
	for domain, file := range HtmlFiles {
		// Load the stored response from the file
		fileContent := LoadFileAsString(t, file)

		// Mock http response
		httpmock.RegisterResponder("GET", domain,
			httpmock.NewStringResponder(200, fileContent))
	}

	// Define test cases
	tests := []struct {
		milliseconds int
		requests     int
		expected     int
	}{
		{100 * 1000, 100, 7},
		{100 * 1000, 0, 0},
		{0, 10, 0},
	}

	// Test the function
	parsedUrl, _ := url.Parse("http://www.parserdigital.com")
	for _, tt := range tests {
		strategy := crawler.NewRecursiveWithLimits(
			parsedUrl, crawler.Limits{Milliseconds: tt.milliseconds, Requests: tt.requests})
		result := strategy.Run()
		info := httpmock.GetCallCountInfo()

		assert.Equal(t, tt.expected, len(result))
		if tt.expected == 7 {
			for link := range HtmlFiles {
				assert.Equal(t, 1, info["GET "+link])
			}
		}
	}
}

func TestRecursiveParallelWithLimits(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responses
	for domain, file := range HtmlFiles {
		// Load the stored response from the file
		fileContent := LoadFileAsString(t, file)

		// Mock http response
		httpmock.RegisterResponder("GET", domain,
			httpmock.NewStringResponder(200, fileContent))
	}

	// Define test cases
	tests := []struct {
		milliseconds int
		requests     int
		expected     int
	}{
		{100 * 1000, 100, 7},
		{100 * 1000, 0, 0},
		{0, 10, 0},
	}

	// Test the function
	parsedUrl, _ := url.Parse("http://www.parserdigital.com")
	for _, tt := range tests {
		strategy := crawler.NewRecursiveParallelWithLimits(
			parsedUrl, crawler.Limits{Milliseconds: tt.milliseconds, Requests: tt.requests})
		result := strategy.Run()
		info := httpmock.GetCallCountInfo()

		assert.Equal(t, tt.expected, len(result))
		if tt.expected == 7 {
			for link := range HtmlFiles {
				assert.Equal(t, 1, info["GET "+link])
			}
		}
	}
}
