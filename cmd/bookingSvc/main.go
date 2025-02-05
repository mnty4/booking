package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/mnty4/booking"
	"github.com/mnty4/booking/internal/db"
)

type ServerConfig struct{}

func run(ctx context.Context, getEnv func(string) string, w io.Writer, args []string) (err error) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()
	fs := flag.NewFlagSet("", flag.ExitOnError)
	envFile := fs.String("env", "", "Optionally specify a file to load env variables")
	err = fs.Parse(args)
	if err != nil {
		return fmt.Errorf("failed to parse args %v: %v", args, err)
	}
	if *envFile != "" {
		err := godotenv.Load(*envFile)
		if err != nil {
			return fmt.Errorf("failed to load environment file %s: %v", *envFile, err)
		}
	}
	logger := log.New(w, "[Booking Service] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	defer logger.Println("Server shutting down gracefully.")

	dbConf := mysql.Config{
		DBName: getEnv("MYSQL_DATABASE"),
		User:   getEnv("MYSQL_USER"),
		Passwd: getEnv("MYSQL_PASSWORD"),
		Net:    getEnv("MYSQL_NET"),
		Addr:   getEnv("MYSQL_ADDR"),
	}
	errCh := make(chan error, 3)
	db, err := db.NewDb(logger, "mysql", dbConf.FormatDSN())
	if err != nil {
		return fmt.Errorf("failed to create db: %v", err)
	}

	serverConf := &ServerConfig{}
	server := booking.NewServer(getEnv, serverConf, logger, db)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = server.Shutdown(ctx)
		if err != nil {
			errCh <- fmt.Errorf("failed shutting down server: %v", err)
		}
		err = db.Close()
		if err != nil {
			errCh <- fmt.Errorf("failed closing db: %v", err)
		}
		close(errCh)
		wg.Done()
	}()
	logger.Printf("Server listening on port: %s\n", getEnv("PORT"))
	if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("unexpected server error: %v", err)
	}
	wg.Wait()
	errs := []error{}
	for err := range errCh {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv, os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
