package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/mnty4/booking"
	"github.com/mnty4/booking/internal/db"
)

type ServerConfig struct{}

func run(ctx context.Context, getEnv func(string) string, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	fs := flag.NewFlagSet("name", flag.ExitOnError)
	envFile := fs.String("env", "", "Optionally specify a file to load env variables")
	fs.Parse(args)
	if *envFile != "" {
		godotenv.Load(*envFile)
	}
	logger := log.New(os.Stdout, "[Booking Service]", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

	dbConfig := mysql.Config{
		DBName: getEnv("MYSQL_DATABASE"),
		User:   getEnv("MYSQL_USER"),
		Passwd: getEnv("MYSQL_PASSWORD"),
		Net:    getEnv("MYSQL_NET"),
		Addr:   getEnv("MYSQL_ADDR"),
	}
	db, err := db.NewDb(logger, "mysql", dbConfig.FormatDSN())
	if err != nil {
		return fmt.Errorf("failed to create db: %v", err)
	}

	config := &ServerConfig{}
	mux := booking.NewServer(config, logger, db)

	return http.ListenAndServe(":"+getEnv("PORT"), mux)
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv, os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
