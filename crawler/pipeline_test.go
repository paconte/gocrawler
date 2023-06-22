package crawler_test

import (
	"io/ioutil"
	"net/url"
	crawler "paconte/parser/crawler"
	"testing"

	"github.com/jarcoal/httpmock"
	"golang.org/x/exp/slices"
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
	fileBytes, err := OpenFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	fileContent := string(fileBytes)

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
	url := "https://parserdigital.com/"
	status := 200

	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Load the stored response from the file
	filePath := "testdata/response.html"
	fileBytes, err := OpenFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	fileContent := string(fileBytes)

	// Mock http response
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(status, fileContent))

	// Test the function
	out := crawler.Parse(crawler.Download(url))
	node := <-out
	_ = html.Node(*node) // Check if node is a html.Node
}

func TestExtract(t *testing.T) {
	baseUrl := "https://parserdigital.com/"
	status := 200

	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Load the stored response from the file
	filePath := "testdata/response.html"
	fileBytes, err := OpenFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	fileContent := string(fileBytes)

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
		if !slices.Contains(result, link) {
			t.Errorf("Extract() = %v, want %v", result, expected)
		}
	}
}
