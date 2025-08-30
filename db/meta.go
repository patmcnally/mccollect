package db

// ReadMeta reads a value from the _meta table.
func (d *DB) ReadMeta(key string) (string, error) {
	var val string
	err := d.conn.QueryRow("SELECT value FROM _meta WHERE key = ?", key).Scan(&val)
	return val, err
}

// WriteMeta upserts a key-value pair in the _meta table.
func (d *DB) WriteMeta(key, value string) error {
	_, err := d.conn.Exec(
		"INSERT INTO _meta (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value",
		key, value,
	)
	return err
}
