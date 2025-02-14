package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func NewTestClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
}

func HealthCheck(ctx context.Context, timeout time.Duration) error {
	url, err := url.Parse(os.Getenv("BASE_URL"))
	if err != nil {
		return err
	}
	url.Path = "healthz"
	return WaitForReady(ctx, timeout, url.String())
}

func WaitForReady(ctx context.Context, timeout time.Duration, url string) error {
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
				return fmt.Errorf("failed to create request: %w", err)
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed making request: %s\n", err)
				continue
			}
			resp.Body.Close()
			if resp.StatusCode != 200 {
				fmt.Fprintf(os.Stderr, "Request failed with status code: %d\n", resp.StatusCode)
			}
			if resp.StatusCode == 200 {
				return nil
			}
		}
	}
}
