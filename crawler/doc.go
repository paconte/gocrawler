/*
Package crawler implements a web crawler.

A web crawler is a program that automatically traverses the web by downloading
the pages and following the links from page to page.

This package contains a web crawler that can be configured to use different
strategies to crawl the web. The strategies are implemented in the strategies.go
file.

The crawler can be configured to use a strategy and a set of limits. The limits
are the maximum number of seconds and requests that the crawler can make. The
crawler will stop when it reaches the limits. The limits are not hard limits and
the crawler may exceed them by a small amount. The crawler will stop as soon as
it can.

# Types of strategies

The crawler can be configured to use one of the following strategies:

## Recursive

This strategy implements a search approach to crawl URLs recursively,
discovering new URLs at each level and continuing the crawling process until
there are no more unvisited URLs.

## Recursive with limits

This strategy implements a search approach to crawl URLs recursively,
discovering new URLs at each level and continuing the crawling process until
there are no more unvisited URLs or the limits are reached.

## Parallel

This strategy implements a search approach to crawl URLs in parallel,
discovering new URLs at each level and continuing the crawling process until
there are no more unvisited URLs.

## Parallel with limits

This strategy implements a search approach to crawl URLs in parallel,
discovering new URLs at each level and continuing the crawling process until
there are no more unvisited URLs or the limits are reached.

# Usage

The following example shows how to use the crawler package to crawl a website
using the Recursive strategy:

	package main

	import (
		"fmt"
		"log"
		"net/url"

		"github.com/paconte/gocrawler"
		"github.com/paconte/crawler/strategies"
	)

	func main() {

		// Create a new URL to crawl.
		url, err := url.Parse("https://www.example.com")
		if err != nil {
			log.Fatal(err)
		}

		// Create a new crawler with the Recursive strategy.
		c := crawler.New(url, strategies.NewRecursive(url))

		// Start crawling.
		visited := c.Crawl()

		// Print the visited URLs.
		fmt.Println(visited)
	}

The following example shows how to use the crawler package to crawl a website
using the Recursive with limits strategy:

	package main

	import (
		"fmt"
		"log"
		"net/url"
		"time"

		"github.com/paconte/gocrawler"
		"github.com/paconte/crawler/strategies"
	)

	func main() {

		// Create a new URL to crawl.
		url, err := url.Parse("https://www.example.com")
		if err != nil {
			log.Fatal(err)
		}

		// Create a new crawler with the Recursive with limits strategy.
		c := crawler.New(url, strategies.NewRecursiveWithLimits(url, crawler.Limits{
			Milliseconds: 1000,
			Requests:     100,
		}))

		// Start crawling.
		visited := c.Crawl()

		// Print the visited URLs.
		fmt.Println(visited)
	}

# Pipeline

The crawler package uses a pipeline to crawl the web. The pipeline is composed
of the following stages:

## Download

# The Download stage downloads the content of a URL and returns a string

## Parse

# The Parse stage parses the content of a URL and returns a slice of URLs

## Extract

# The Extract stage extracts the URLs that match the root URL and returns a slice

## Collect

# The Collect stage collects the URLs and returns a map

## MapToList

The MapToList stage converts a map to a slice
*/
package crawler
