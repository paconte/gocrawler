package crawler_test

import (
	"io/ioutil"
	"net/url"
	crawler "paconte/parser/crawler"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestDownload(t *testing.T) {
	url := "https://parserdigital.com/"
	status := 200

	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Load the stored response from the file
	filePath := "testdata/response.html"
	fileContent := LoadFileAsString(t, filePath)

	// Mock http response
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(status, fileContent))

	// Test the function
	out := crawler.Download(url)
	response := <-out

	if response.StatusCode != status {
		t.Errorf("Download(%q) = %v, want %v", url, response.StatusCode, 200)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Download(%q) = %v, want %v", url, err, nil)
	}
	if string(data) != fileContent {
		t.Errorf("Download(%q) = %v, want %v", url, response.Body, fileContent)
	}
}

func TestParse(t *testing.T) {
	status := 200
	url := "https://parserdigital.com/"
	filePath := "testdata/response.html"

	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Load the stored response from the file
	fileContent := LoadFileAsString(t, filePath)

	// Mock http response
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(status, fileContent))

	// Test the function
	out := crawler.Parse(crawler.Download(url))
	node := <-out
	_ = html.Node(*node) // Check if node is a html.Node
}

func TestExtract(t *testing.T) {
	filePath := "testdata/response.html"
	baseUrl := "https://parserdigital.com/"
	status := 200

	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Load the stored response from the file
	fileContent := LoadFileAsString(t, filePath)

	// Mock http response
	httpmock.RegisterResponder("GET", baseUrl,
		httpmock.NewStringResponder(status, fileContent))

	// Call the function
	parsedURL, err := url.Parse(baseUrl)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	out := crawler.Extract(crawler.Parse(crawler.Download(baseUrl)), parsedURL)

	// Check the result
	result := []string{}
	for link := range out {
		result = append(result, link)
	}
	expected := TargetLinks
	for _, link := range expected {
		//assert.Containsf(t, result, link, "Extract() = %v, want %v", result, expected)
		assert.Contains(t, result, link)
	}
}

func TestCollectMap(t *testing.T) {
	ch := make(chan string, 2*len(TargetLinks))
	for _, link := range TargetLinks {
		ch <- link
		ch <- link
	}
	close(ch)

	result := crawler.CollectMap(ch)
	assert.Equal(t, len(TargetLinks), len(result))
}

func TestMapToList(t *testing.T) {
	targetsMap := make(map[string]bool)
	for _, target := range TargetLinks {
		targetsMap[target] = true
	}
	result := crawler.MapToList(targetsMap)
	assert.Equal(t, len(result), len(TargetLinks))
}
