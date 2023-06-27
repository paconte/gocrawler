package crawler_test

import (
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	crawler "github.com/paconte/gocrawler/crawler"
)

func BLoadFileAsString(b *testing.B, filePath string) string {
	fileBytes, err := OpenFile(filePath)
	if err != nil {
		b.Fatalf("Error reading file: %v", err)
	}
	return string(fileBytes)
}

func BenchmarkOneLevel(b *testing.B) {
	// Register mock responses
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	domain := "http://www.parserdigital.com"
	fileContent := BLoadFileAsString(b, "testdata/treeLevel1.html")
	httpmock.RegisterResponder("GET", domain,
		httpmock.NewStringResponder(200, fileContent))

	// Test the function
	parsedUrl, _ := url.Parse(domain)
	c := crawler.NewOneLevel()
	for i := 0; i < b.N; i++ {
		c.Run(parsedUrl)
	}
}

func BenchmarkRecursive(b *testing.B) {
	// Register mock responses
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for domain, file := range HtmlFiles {
		fileContent := BLoadFileAsString(b, file)
		httpmock.RegisterResponder("GET", domain,
			httpmock.NewStringResponder(200, fileContent))
	}

	// Test the function
	parsedUrl, _ := url.Parse("http://www.parserdigital.com")
	c := crawler.NewRecursive()
	for i := 0; i < b.N; i++ {
		c.Run(parsedUrl)
	}
}

func BenchmarkRecursiveParallel(b *testing.B) {
	// Register mock responses
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for domain, file := range HtmlFiles {
		fileContent := BLoadFileAsString(b, file)
		httpmock.RegisterResponder("GET", domain,
			httpmock.NewStringResponder(200, fileContent))
	}

	// Test the function
	parsedUrl, _ := url.Parse("http://www.parserdigital.com")
	c := crawler.NewRecursiveParallel()
	for i := 0; i < b.N; i++ {
		c.Run(parsedUrl)
	}
}
