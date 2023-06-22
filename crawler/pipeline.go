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
			resp, _ := http.Get(u)
			out <- resp
		}
		close(out)
	}()
	return out
}

func Parse(nodes <-chan *http.Response) <-chan *html.Node {
	out := make(chan *html.Node)
	go func() {
		for resp := range nodes {
			doc, _ := html.Parse(resp.Body)
			out <- doc
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

func Print(links <-chan string) {
	for link := range links {
		println(link)
	}
}
