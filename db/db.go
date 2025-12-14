package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps a SQLite connection with domain-specific methods.
type DB struct {
	conn *sql.DB
}

// Open opens (or creates) a SQLite database at path and sets pragmas.
// It creates the parent directory if it does not exist.
func Open(path string) (*DB, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db dir: %w", err)
		}
	}
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if _, err := conn.Exec("PRAGMA journal_mode = WAL"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("set WAL: %w", err)
	}
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("set FK: %w", err)
	}
	return &DB{conn: conn}, nil
}

// InitSchema creates all tables and indexes if they don't exist.
func (d *DB) InitSchema() error {
	_, err := d.conn.Exec(schema)
	return err
}

// DropAll drops all tables for a full rebuild.
func (d *DB) DropAll() error {
	tables := []string{"collection_packs", "collections", "cards", "sets", "packs", "_meta"}
	for _, t := range tables {
		if _, err := d.conn.Exec("DROP TABLE IF EXISTS " + t); err != nil {
			return fmt.Errorf("drop %s: %w", t, err)
		}
	}
	return nil
}

// Conn returns the underlying *sql.DB for use in transactions.
func (d *DB) Conn() *sql.DB {
	return d.conn
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.conn.Close()
}
