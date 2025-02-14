//go:build integration

package api

import (
	"context"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mnty4/booking/internal/app"
	"github.com/mnty4/booking/utils"
)

func TestMain(m *testing.M) {
	logger := log.New(os.Stderr, "[Integration Setup] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	ctx := context.Background()
	args := []string{"-env=../.env.tcp.test"}
	if err := app.ParseFlags(args); err != nil {
		logger.Fatal(err)
	}
	server, err := app.NewServer(io.Discard)
	if err != nil {
		logger.Fatal(err)
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run(ctx)
	}()
	go func() {
		errCh <- utils.HealthCheck(ctx, 5*time.Second)
	}()
	if err := <-errCh; err != nil {
		logger.Fatal(err)
	}
	start := time.Now()
	code := m.Run()
	logger.Println("run all tests", time.Since(start))
	// empty database after tests finish
	start = time.Now()
	if err := utils.TruncateTables(server.Db, []string{"users"}); err != nil {
		logger.Fatal(err)
	}
	logger.Println("truncate", time.Since(start))
	start = time.Now()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
	logger.Println("shutdown", time.Since(start))
	os.Exit(code)
}
