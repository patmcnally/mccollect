package db

import (
	"database/sql"

	"github.com/patmcnally/mccollect/model"
)

func (d *DB) UpsertPacks(tx *sql.Tx, packs []model.Pack) error {
	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO packs (code,name,cgdb_id,octgn_id,date_release,pack_type_code,position,size) VALUES (?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, p := range packs {
		if _, err := stmt.Exec(p.Code, p.Name, p.CgdbID, p.OctgnID, p.DateRelease, p.PackTypeCode, p.Position, p.Size); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) ListPacks() ([]model.Pack, error) {
	rows, err := d.conn.Query("SELECT code,name,cgdb_id,octgn_id,date_release,pack_type_code,position,size FROM packs ORDER BY date_release,position")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var packs []model.Pack
	for rows.Next() {
		var p model.Pack
		if err := rows.Scan(&p.Code, &p.Name, &p.CgdbID, &p.OctgnID, &p.DateRelease, &p.PackTypeCode, &p.Position, &p.Size); err != nil {
			return nil, err
		}
		packs = append(packs, p)
	}
	return packs, rows.Err()
}
