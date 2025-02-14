package app

// import (
// 	"context"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/mnty4/booking/helpers"
// )

// // test run start up
// func TestRun(t *testing.T) {
// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)
// 	t.Cleanup(cancel)

// 	getEnv := helpers.NewTestEnv()

// 	go func() {
// 		if err := Run(ctx, getEnv, os.Stdout, nil); err != nil {
// 			t.Errorf("starting server failed: %v", err)
// 		}
// 		cancel()
// 	}()
// 	url := "http://localhost:" + getEnv("PORT") + "/healthz"
// 	err := helpers.WaitForReady(ctx, 5*time.Second, url)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
