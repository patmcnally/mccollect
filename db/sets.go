package db

import (
	"database/sql"

	"github.com/patmcnally/mccollect/model"
)

// UpsertSets inserts or replaces sets within a transaction.
func (d *DB) UpsertSets(tx *sql.Tx, sets []model.Set) error {
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO sets (code, name, card_set_type_code) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, s := range sets {
		if _, err := stmt.Exec(s.Code, s.Name, s.CardSetTypeCode); err != nil {
			return err
		}
	}
	return nil
}
