package cmd

import (
	"fmt"
	"os"

	"github.com/paconte/gocrawler/crawler"

	"github.com/spf13/cobra"
)

var (
	// flags
	strategy string
	url      string

	rootCmd = NewRootCmd()
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crawler",
		Short: "Crawler generates a tree of subdomains for a given domain.",
		Long:  `A crawler that support different algorithms when searching for subdomains of a web site.`,
		Args:  cobra.MatchAll(cobra.MaximumNArgs(0)),
		Run: func(cmd *cobra.Command, args []string) {
			runCrawler()
		},
	}
	cmd.PersistentFlags().StringVarP(&strategy, "alg", "a", "", "The algorithm used for search, either OneLevel or Recursive")
	cmd.PersistentFlags().StringVarP(&url, "url", "u", "", "The url to search for subdomains")
	return cmd
}

func Execute(cmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runCrawler() {
	crawler, err := crawler.NewCrawler(strategy)

	if err != nil {
		fmt.Println(err)
		return
	}

	crawler.Run(url)
	crawler.SortLinks()

	for _, link := range crawler.GetResult() {
		fmt.Println(link)
	}
}
