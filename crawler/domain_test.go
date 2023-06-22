package crawler_test

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"os"
	crawler "paconte/parser/crawler"
	"path/filepath"
	"testing"

	"golang.org/x/net/html"
)

var TargetLinks []string = []string{
	"https://parserdigital.com/career-opportunities/",
	"https://parserdigital.com/category/aws/",
	"https://parserdigital.com/privacy-policy/",
	"https://parserdigital.com/life-at-parser/",
	"https://parserdigital.com/client-story/ey/",
	"https://parserdigital.com/recruitment-privacy-policy/",
	"https://parserdigital.com/how-we-work/",
	"https://parserdigital.com/aws-understanding-the-components-of-a-vpc-chapter-2-private-vpcs-and-complex-connections/",
	"https://parserdigital.com/client-story/modulr/",
	"https://parserdigital.com/what-we-do/",
	"https://parserdigital.com/career-accelerator/",
	"https://parserdigital.com/our-work/",
	"https://parserdigital.com/expertise/",
	"https://parserdigital.com/approach-to-esg/",
	"https://parserdigital.com/cookie-policy/",
	"https://parserdigital.com/contact-us/",
	"https://parserdigital.com/client-story/doctorlink/",
	"https://parserdigital.com/about-us/",
	"https://parserdigital.com/the-journey-to-our-new-values/",
	"https://parserdigital.com/aws-understanding-the-components-of-a-vpc-chapter-1-from-the-internet-to-the-public-vpc/",
	"https://parserdigital.com/category/qa/",
	"https://parserdigital.com/category/women-in-tech/",
	"https://parserdigital.com/my-experience-at-stareast/",
}

func OpenFile(relativePath string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(wd, relativePath)
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func TestIsSubdomain(t *testing.T) {
	tests := []struct {
		link     string
		domain   string
		expected bool
	}{
		{"https://parserdigital.com/what-we-do/", "https://parserdigital.com/", true},
		{"https://parserdigital.com/what-we-do/", "https://parserdigital.com", true},
		{"https://parserdigital.com/category/qa/", "https://parserdigital.com", true},
		{"https://parserdigital.com/", "https://parserdigital.com/", false}, // Same url is not a subdomain
		{"http://example.org", "https://parserdigital.com/", false},
		{"invalid-url", "https://parserdigital.com/", false},
		{"#", "https://parserdigital.com", false},
	}

	for _, tt := range tests {
		domain, err := url.Parse(tt.domain)
		if err != nil {
			t.Error(err)
		}
		result := crawler.IsSubdomain(tt.link, domain)
		if result != tt.expected {
			t.Errorf("isSubdomain(%q, %q) = %v, want %v", tt.link, tt.domain, result, tt.expected)
		}
	}
}

func TestGetSubdomains(t *testing.T) {
	// Load the stored response from the file
	filePath := "testdata/response.html"
	file, err := OpenFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	// Parse the HTML document
	body := ioutil.NopCloser(bytes.NewReader(file))
	doc, err := html.Parse(body)
	if err != nil {
		t.Fatal(err)
	}
	// Get the subdomains
	domain, err := url.Parse("https://parserdigital.com/")
	if err != nil {
		t.Fatal(err)
	}
	result := crawler.GetSubdomains(doc, domain)

	// Check the result
	expected := TargetLinks

	if len(result) != len(expected) {
		t.Errorf("GetSubdomains() returned %d links, want %d", len(result), len(expected))
	}

	for _, k := range expected {
		_, ok := result[k]
		if !ok {
			t.Errorf("GetSubdomains() missing link %q", k)
		}
	}
}
