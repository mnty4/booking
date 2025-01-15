package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

// test run start up

func TestRun(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	PORT := "8080"
	getEnv := func(key string) string {
		switch key {
		case "PORT":
			return PORT
		default:
			return ""
		}
	}
	go func() {
		err := run(ctx, getEnv, nil, nil)
		t.Errorf("Run failed: %v\n", err)
		cancel()
	}()
	url := "http://localhost:" + PORT + "/healthz"
	err := waitForReady(ctx, 5*time.Second, url)
	if err != nil {
		t.Error(err)
	}
}

func waitForReady(ctx context.Context, timeout time.Duration, url string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	client := http.Client{}
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return fmt.Errorf("Failed to create request: %w", err)
			}
			res, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed making request: %s\n", err)
				continue
			}
			res.Body.Close()
			if res.StatusCode != 200 {
				fmt.Fprintf(os.Stderr, "Request failed with status code: %d\n", res.StatusCode)
			}
			if res.StatusCode == 200 {
				return nil
			}
		}
	}
}

// test run and graceful shutdown
