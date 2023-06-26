package main

import (
	"github.com/paconte/gocrawler/cmd"
)

func main() {
	cmd.Execute(cmd.NewRootCmd())
}
