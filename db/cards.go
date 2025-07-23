package db

import (
	"database/sql"

	"github.com/patmcnally/mccollect/model"
)

func (d *DB) UpsertCards(tx *sql.Tx, cards []model.Card) error {
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO cards (code,pack_code,name,type_code,faction_code) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, c := range cards {
		if _, err := stmt.Exec(c.Code, c.PackCode, c.Name, c.TypeCode, c.FactionCode); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) DeleteCardsByPack(tx *sql.Tx, packCode string) error {
	_, err := tx.Exec("DELETE FROM cards WHERE pack_code = ?", packCode)
	return err
}
