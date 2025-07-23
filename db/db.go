package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	return &DB{conn: conn}, nil
}

func (d *DB) InitSchema() error {
	_, err := d.conn.Exec(schema)
	return err
}

func (d *DB) DropAll() error {
	for _, t := range []string{"cards", "sets", "packs"} {
		if _, err := d.conn.Exec("DROP TABLE IF EXISTS " + t); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) Conn() *sql.DB { return d.conn }
func (d *DB) Close() error  { return d.conn.Close() }
