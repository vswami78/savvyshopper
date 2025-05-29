package main

import (
	"context"
	"os"

	"savvyshopper/runner"
)

func main() {
	if err := runner.Run(context.Background(), os.Args[1:], os.Stdout); err != nil {
		os.Exit(1)
	}
}
