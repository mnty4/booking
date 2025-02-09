//go:build integration

package main

// test run start up
// func TestRun(t *testing.T) {
// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)
// 	t.Cleanup(cancel)
// 	PORT := "8080"
// 	getEnv := func(key string) string {
// 		switch key {
// 		case "PORT":
// 			return PORT
// 		case "MYSQL_DATABASE":
// 			return "bookingtest"
// 		case "MYSQL_USER":
// 			return "test_user"
// 		case "MYSQL_PASSWORD":
// 			return "devpass"
// 		case "MYSQL_NET":
// 			return "tcp"
// 		case "MYSQL_ADDR":
// 			return "127.0.0.1:3306"
// 		default:
// 			return ""
// 		}
// 	}
// 	go func() {
// 		if err := Run(ctx, getEnv, os.Stdout, nil); err != nil {
// 			t.Errorf("starting server failed: %v", err)
// 		}
// 		cancel()
// 	}()
// 	url := "http://localhost:" + PORT + "/healthz"
// 	err := waitForReady(ctx, 20*time.Second, url)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func waitForReady(ctx context.Context, timeout time.Duration, url string) error {
// 	ctx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()
// 	client := http.Client{}
// 	ticker := time.NewTicker(500 * time.Millisecond)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case <-ticker.C:
// 			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
// 			if err != nil {
// 				return fmt.Errorf("Failed to create request: %w", err)
// 			}
// 			res, err := client.Do(req)
// 			if err != nil {
// 				fmt.Fprintf(os.Stderr, "Failed making request: %s\n", err)
// 				continue
// 			}
// 			res.Body.Close()
// 			if res.StatusCode != 200 {
// 				fmt.Fprintf(os.Stderr, "Request failed with status code: %d\n", res.StatusCode)
// 			}
// 			if res.StatusCode == 200 {
// 				return nil
// 			}
// 		}
// 	}
// }
