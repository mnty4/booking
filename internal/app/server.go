package app

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/mnty4/booking/api"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	HTTPServer *http.Server
	Db         *sql.DB
	Logger     *log.Logger
}

// Start the server (blocking) and closes on parent ctx close, os.Interrupt, or syscall.SIGTERM
func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		s.Logger.Printf("Server listening on %s\n", os.Getenv("BASE_URL"))
		if err := s.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-ctx.Done()
		s.Logger.Println("Shutting down server...")
		return s.Shutdown(ctx)
	})
	return g.Wait()
}

// Gracefully shutdown the server and its components.
func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	start := time.Now()
	if err := s.Db.Close(); err != nil {
		return err
	}
	s.Logger.Println("database close", time.Since(start))
	start = time.Now()
	if err := s.HTTPServer.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	s.Logger.Println("http server shutdown", time.Since(start))
	return nil
}

func NewServer(out io.Writer) (*Server, error) {
	logger := NewLogger(out)
	db, err := NewDb()
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %v", err)
	}
	httpServer := NewHTTPServer(logger, db)

	return &Server{
		Db:         db,
		HTTPServer: httpServer,
		Logger:     logger,
	}, nil
}

// parse args as flags
func ParseFlags(args []string) error {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	envFile := fs.String("env", "", "Optionally specify a file to load env variables")
	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("failed to parse args %v: %v", args, err)
	}
	if *envFile != "" {
		err := godotenv.Load(*envFile)
		if err != nil {
			return fmt.Errorf("failed to load env file %q: %v", *envFile, err)
		}
	}
	return nil
}

// Create a http.Server to be served by the application
func NewHTTPServer(logger *log.Logger, db *sql.DB) *http.Server {
	mux := http.NewServeMux()
	validate := validator.New()
	api.AddRoutes(logger, db, validate, mux)
	return &http.Server{
		Handler: mux,
		Addr:    net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT")),
	}
}

func NewDb() (*sql.DB, error) {
	config := mysql.Config{
		DBName: os.Getenv("MYSQL_DATABASE"),
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Net:    os.Getenv("MYSQL_NET"),
		Addr:   os.Getenv("MYSQL_ADDR"),
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("error opening db: %v", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging db: %v", err)
	}
	return db, nil
}

func NewLogger(w io.Writer) *log.Logger {
	return log.New(w, "[Booking Service] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
}
