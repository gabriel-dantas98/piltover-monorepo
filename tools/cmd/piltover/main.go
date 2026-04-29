package main

import (
	"fmt"
	"os"
)

const Version = "0.0.1"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("piltover", Version)
		return
	}
	fmt.Fprintln(os.Stderr, "piltover: not implemented yet")
	os.Exit(1)
}
