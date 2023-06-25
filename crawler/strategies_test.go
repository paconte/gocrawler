package crawler_test

import (
	"net/url"
	"paconte/parser/crawler"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

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
		alg := crawler.NewOneLevel()
		result := alg.Run(parsedUrl)
		assert.Equal(t, tt.expected, len(result))
	}
}

func TestRecursive(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responses
	htmlFiles := map[string]string{
		"http://www.parserdigital.com":   "testdata/treeLevel1.html",
		"http://www.parserdigital.com/A": "testdata/treeLevel2A.html",
		"http://www.parserdigital.com/B": "testdata/treeLevel2B.html",
		"http://www.parserdigital.com/C": "testdata/treeLevel3C.html",
		"http://www.parserdigital.com/D": "testdata/treeLevel3D.html",
		"http://www.parserdigital.com/E": "testdata/treeLevel3E.html",
		"http://www.parserdigital.com/F": "testdata/treeLevel3F.html",
	}
	for domain, file := range htmlFiles {
		// Load the stored response from the file
		fileContent := LoadFileAsString(t, file)

		// Mock http response
		httpmock.RegisterResponder("GET", domain,
			httpmock.NewStringResponder(200, fileContent))
	}

	// Test the function
	parsedUrl, _ := url.Parse("http://www.parserdigital.com")
	alg := crawler.NewRecursive()
	result := alg.Run(parsedUrl)
	info := httpmock.GetCallCountInfo()

	assert.Equal(t, 7, len(result))
	for link := range htmlFiles {
		assert.Equal(t, 1, info["GET "+link])
	}
}
