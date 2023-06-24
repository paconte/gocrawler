package main

import (
	"fmt"
	"os"
	"paconte/parser/crawler"
)

func main() {
	runCrawler()
}

func runCrawler() {
	rootUrl := "https://parserdigital.com/"
	args := os.Args[1:]
	crawler, err := crawler.NewCrawler(rootUrl, args[0])

	if err != nil {
		fmt.Println(err)
		return
	}

	crawler.Run()
	crawler.SortLinks()
	crawler.PrintLinks()
}
