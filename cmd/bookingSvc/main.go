package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mnty4/booking/internal/app"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx, os.Getenv, os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
