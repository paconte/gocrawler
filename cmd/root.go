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
	ms       int
	reqs     int
	// rootCmd represents the base command when called without any subcommands.
	rootCmd = NewRootCmd()
)

// NewRootCmd creates a new instance of the root command.
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

	// Define command flags
	cmd.PersistentFlags().StringVarP(&strategy, "strategy", "s", "", "The algorithm used for search, either OneLevel or Recursive")
	cmd.PersistentFlags().StringVarP(&url, "url", "u", "", "The url to search for subdomains")
	cmd.PersistentFlags().IntVarP(&ms, "milliseconds", "m", 0, "The ms to limit the search")
	cmd.PersistentFlags().IntVarP(&reqs, "requests", "r", 0, "The requests to limit the search")

	return cmd
}

// Execute executes the root command.
func Execute(cmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runCrawler runs the web crawler using the specified strategy and URL.
func runCrawler() {
	res, err := crawler.Run(url, strategy, crawler.Limits{Milliseconds: ms, Requests: reqs})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, link := range res {
		fmt.Println(link)
	}
}
