package main

import (
	"context"
	"log"
	"os"

	"github.com/mnty4/booking/internal/app"
)

func main() {
	ctx := context.Background()
	if err := app.ParseFlags(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	server, err := app.NewServer(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server shutdown gracefully.")
}
