package main

import (
	"fmt"
	"os"

	"github.com/nokamoto/kaas-operator-prototype/internal/cli"
)

func main() {
	cmd := cli.New()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
