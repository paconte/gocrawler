package crawler_test

import (
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
	url := "http://www.parserdigital.com"
	fileContent := BLoadFileAsString(b, "testdata/treeLevel1.html")
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, fileContent))

	// Test the function
	for i := 0; i < b.N; i++ {
		crawler.Run(url, "OneLevel")
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
	for i := 0; i < b.N; i++ {
		crawler.Run("http://www.parserdigital.com", "Recursive")
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
	for i := 0; i < b.N; i++ {
		crawler.Run("http://www.parserdigital.com", "RecursiveParallel")
	}
}
