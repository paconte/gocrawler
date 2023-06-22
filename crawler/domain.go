package crawler

import (
	"net/url"

	"golang.org/x/net/html"
)

func IsSubdomain(link string, domain *url.URL) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}
	result := u.Hostname() == domain.Hostname() &&
		u.RequestURI() != domain.RequestURI() // Avoid same url
	return result
}

func GetSubdomains(node *html.Node, domain *url.URL) map[string]bool {
	var links = make(map[string]bool)
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" && IsSubdomain(attr.Val, domain) {
				links[attr.Val] = true
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		child_links := GetSubdomains(child, domain)
		for k := range child_links {
			links[k] = true
		}
	}

	return links
}
