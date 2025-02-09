/* defines useful helpers for testing the service */
package helpers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

// func NewTestDb(logger *log.Logger) (*sql.DB, error) {

// 	return db.NewDb(logger, "mysql", dsn)
// }

func NewTestEnv() func(key string) string {
	return func(key string) string {
		switch key {
		case "PORT":
			return "8080"
		case "MYSQL_DATABASE":
			return "bookingtest"
		case "MYSQL_USER":
			return "test_user"
		case "MYSQL_PASSWORD":
			return "devpass"
		case "MYSQL_NET":
			return "tcp"
		case "MYSQL_ADDR":
			return "127.0.0.1:3306"
		default:
			return ""
		}
	}
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
