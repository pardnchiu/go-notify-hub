package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/lib/pq"
)

var (
	PG *sql.DB
)

type db struct {
	host       string
	port       string
	user       string
	password   string
	dbName     string
	sslMode    string
	searchPath string
}

func NewDB() (*sql.DB, error) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	config := db{
		host:       "192.168.0.243",
		port:       "5432",
		user:       "postgres",
		password:   "0123456789",
		dbName:     "daily_push",
		sslMode:    "disable",
		searchPath: "public",
	}

	link := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		config.host, config.port, config.user,
		strings.ReplaceAll(config.password, "'", "\\'"),
		config.dbName, config.sslMode, config.searchPath,
	)
	db, err := sql.Open("postgres", link)
	if err != nil {
		slog.Error("open db failed", "error", err)
		os.Exit(1)
	}

	if err := db.PingContext(ctx); err != nil {
		slog.Error("ping db failed", "error", err)
		os.Exit(1)
	}
	PG = db

	slog.Info("database connected",
		"host", config.host,
		"port", config.port,
		"database", config.dbName,
	)

	return PG, nil
}