package main

import (
	"fmt"
	"os"

	"github.com/gabriel-dantas98/piltover-monorepo/tools/internal/cli"
)

const Version = "0.0.1"

func main() {
	root := cli.NewRootCmd()
	root.Version = Version
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "piltover:", err)
		os.Exit(1)
	}
}
