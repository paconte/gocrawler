package crawler

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func Download(url ...string) <-chan *http.Response {
	out := make(chan *http.Response)
	go func() {
		for _, u := range url {
			resp, err := http.Get(u)
			if err == nil {
				out <- resp
			}
		}
		close(out)
	}()
	return out
}

func Parse(nodes <-chan *http.Response) <-chan *html.Node {
	out := make(chan *html.Node)
	go func() {
		for resp := range nodes {
			doc, err := html.Parse(resp.Body)
			if err == nil {
				out <- doc
			}
		}
		close(out)
	}()
	return out
}

func Extract(nodes <-chan *html.Node, url *url.URL) <-chan string {
	out := make(chan string)
	go func() {
		for node := range nodes {
			for link := range GetSubdomains(node, url) {
				out <- link
			}
		}
		close(out)
	}()
	return out
}

func CollectMap(links <-chan string) map[string]bool {
	result := make(map[string]bool)
	for link := range links {
		result[link] = true
	}
	return result
}

func MapToList(links map[string]bool) []string {
	result := make([]string, 0, len(links))
	for link := range links {
		result = append(result, link)
	}
	return result
}
