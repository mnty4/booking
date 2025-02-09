//go:build integration

package api

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/mnty4/booking/helpers"
	"github.com/mnty4/booking/internal/app"
)

func TestMain(m *testing.M) {
	logger := log.New(os.Stderr, "[Integration Setup] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	getEnv := helpers.NewTestEnv()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		app.RunTestServer(ctx, logger, getEnv, cancel)
		wg.Done()
	}()
	url := "http://localhost:" + getEnv("PORT") + "/healthz"
	if err := helpers.WaitForReady(ctx, 5*time.Second, url); err != nil {
		logger.Fatal(err)
	}
	code := m.Run()
	cancel()
	wg.Wait()
	os.Exit(code)
}
