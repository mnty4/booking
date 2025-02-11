package app

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
)

type ServerConfig struct{}

func Run(ctx context.Context, getEnv func(string) string, w io.Writer, args []string) (err error) {
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
	db, err := NewDb(logger, "mysql", dbConf.FormatDSN())
	if err != nil {
		return fmt.Errorf("failed to create db: %v", err)
	}

	serverConf := &ServerConfig{}
	server := NewServer(getEnv, serverConf, logger, db)
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

func RunTestServer(ctx context.Context, logger *log.Logger, getEnv func(string) string, cancel context.CancelFunc) {
	defer cancel()
	if err := Run(ctx, getEnv, io.Discard, nil); err != nil {
		logger.Fatalf("starting server failed: %v", err)
	}
}

// func RunTestServer(ctx context.Context) error {
// 	ctx, cancel := context.WithCancel(ctx)
// 	errCh := make(chan error)
// 	getEnv := helpers.NewTestEnv()
// 	go func() {
// 		if err := Run(ctx, getEnv, os.Stdout, nil); err != nil {
// 			errCh <- fmt.Errorf("starting server failed: %v", err)
// 		}
// 		cancel()
// 	}()
// 	go func() {
// 		url := "http://localhost:" + getEnv("PORT") + "/healthz"
// 		err := helpers.WaitForReady(ctx, 5*time.Second, url)
// 		if err != nil {
// 			errCh <- err
// 		}
// 		cancel()
// 	}()
// 	<-ctx.Done()
// 	close(errCh)
// 	select {
// 	case err := <-errCh:
// 		return err
// 	default:
// 		return nil
// 	}

// }
