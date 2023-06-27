package crawler

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

// Download asynchronously downloads the specified URLs and returns a channel of *http.Response.
// Each response will be sent on the channel as it becomes available.
// The returned channel will be closed once all downloads are complete.
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

// Parse asynchronously parses the HTML nodes in the *http.Response objects received on the input channel.
// It returns a channel of *html.Node containing the parsed nodes.
// The returned channel will be closed once all parsing is complete.
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

// Parse asynchronously parses the HTML nodes in the *http.Response objects received on the input channel.
// It returns a channel of *html.Node containing the parsed nodes.
// The returned channel will be closed once all parsing is complete.
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

// CollectMap reads strings from the input channel and collects them into a map.
// It returns a map containing the collected strings as keys, with a value of true.
func CollectMap(links <-chan string) map[string]bool {
	result := make(map[string]bool)
	for link := range links {
		result[link] = true
	}
	return result
}

// MapToList converts a map of strings to a slice of strings.
// It returns a slice containing all the keys from the input map.
func MapToList(links map[string]bool) []string {
	result := make([]string, 0, len(links))
	for link := range links {
		result = append(result, link)
	}
	return result
}
