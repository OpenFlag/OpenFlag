package main

import (
	"os"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/cmd"
)

const (
	exitFailure = 1
)

func main() {
	root := cmd.NewRootCommand()

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(exitFailure)
		}
	}
}
