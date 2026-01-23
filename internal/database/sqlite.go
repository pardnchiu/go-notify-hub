package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

type ContainerRecord struct {
	LocalDir  string `json:"local_dir"`
	RemoteDir string `json:"remote_dir"`
	File      string `json:"file"`
	Hostname  string `json:"hostname"`
	IP        string `json:"ip"`
	Content   string `json:"content"`
}

var (
	DB *SQLite
)

func NewSQLite(dbPath string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	s := &SQLite{db: db}
	if err := s.create(); err != nil {
		db.Close()
		return nil, err
	}

	DB = s

	return s, nil
}

func (s *SQLite) create() error {
	schema, err := os.ReadFile("sql/create.sql")
	if err != nil {
		return err
	}

	_, err = s.db.Exec(string(schema))
	return err
}

func (s *SQLite) Close() error {
	return s.db.Close()
}
