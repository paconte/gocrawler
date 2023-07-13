# gocrawler
A simple web crawler written in go.

## Project Description

The project provides different strategies for crawling websites and extracting information from them.


## Installation

1. Make sure you have Go installed on your system.
2. Clone the project repository:

   ```shell
   git clone https://github.com/paconte/gocrawler.git
   ```
3. Change to the project directory:

   ```shell
   cd gocrawler
   ```
4. Build the project:
   ```shell
   go build
   ```


## Usage

Once you have built the executable file you can run the crawler and utilie different strategies:

Example usage:
   ```shell
   # Crawl only the pricipal site for subdomains
   ./gocrawler -a OneLevel  -u https://as.com
   ```

   ```shell
   # Crawl recursively all the found subdomains
   ./gocrawler -a Recursive  -u https://as.com
   ```

## Development
To use the project, you can import the relevant packages into your own Go code and utilize the provided strategies.

```go
package main

import (
	"fmt"
	"net/url"

	"github.com/paconte/gocrawler/crawler"
)

func main() {
   crawler, err := crawler.NewCrawler("Recursive")
	result, err := crawler.Run("https://example.com")
	fmt.Println(result)
}
```

## License
This project is licensed under the GNU General Public License.