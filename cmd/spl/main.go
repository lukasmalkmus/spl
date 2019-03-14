package main

import (
	"os"

	"github.com/lukasmalkmus/spl/cmd/spl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
