package main

import (
	"fmt"
	crawler "paconte/parser/crawler"
)

func main() {
	runCrawler()
}

func runCrawler() {
	rootUrl := "https://parserdigital.com/"
	//crawler, err := net.NewCrawler(rootUrl, "bruteforce")
	crawler, err := crawler.NewCrawler(rootUrl, "singlevisit")
	if err != nil {
		fmt.Println(err)
		return
	}
	crawler.Run()
	crawler.SortLinks()
	crawler.PrintLinks()
}
