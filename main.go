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
	//rootUrl := "http://www.queoscontomateychampis.com/"
	args := os.Args[1:]
	crawler, err := crawler.NewCrawler(args[0])

	if err != nil {
		fmt.Println(err)
		return
	}

	crawler.Run(rootUrl)
	crawler.SortLinks()
	for _, link := range crawler.GetResult() {
		fmt.Println(link)
	}
}
