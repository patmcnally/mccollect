package db

import (
	"database/sql"
	"fmt"

	"github.com/patmcnally/mccollect/model"
)

func (d *DB) EnsureCollection(name string) (int, error) {
	_, err := d.conn.Exec(`INSERT INTO collections (name) VALUES (?) ON CONFLICT(name) DO NOTHING`, name)
	if err != nil { return 0, fmt.Errorf("ensure collection: %w", err) }
	var id int
	err = d.conn.QueryRow("SELECT id FROM collections WHERE name = ?", name).Scan(&id)
	return id, err
}

func (d *DB) SetPackOwned(colID int, packCode string, owned bool) error {
	ownedInt := 0; if owned { ownedInt = 1 }
	_, err := d.conn.Exec(
		`INSERT INTO collection_packs (collection_id,pack_code,owned) VALUES (?,?,?)
		 ON CONFLICT(collection_id,pack_code) DO UPDATE SET owned=excluded.owned`,
		colID, packCode, ownedInt)
	return err
}

func (d *DB) ListPackOwnership(collectionID int) ([]model.PackOwnership, error) {
	rows, err := d.conn.Query(`
		SELECT p.code,p.name,p.cgdb_id,p.octgn_id,p.date_release,p.pack_type_code,p.position,p.size,
		       COALESCE(cp.owned,0)
		FROM packs p LEFT JOIN collection_packs cp ON p.code=cp.pack_code AND cp.collection_id=?
		ORDER BY p.date_release,p.position`, collectionID)
	if err != nil { return nil, err }
	defer rows.Close()
	var result []model.PackOwnership
	for rows.Next() {
		var po model.PackOwnership; var owned int
		rows.Scan(&po.Pack.Code,&po.Pack.Name,&po.Pack.CgdbID,&po.Pack.OctgnID,
			&po.Pack.DateRelease,&po.Pack.PackTypeCode,&po.Pack.Position,&po.Pack.Size,&owned)
		po.Owned = owned == 1; result = append(result, po)
	}
	return result, rows.Err()
}

func (d *DB) CollectionStats(collectionID int) (owned, total int, err error) {
	err = d.conn.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(CASE WHEN cp.owned=1 THEN 1 ELSE 0 END),0)
		FROM packs p LEFT JOIN collection_packs cp ON p.code=cp.pack_code AND cp.collection_id=?`,
		collectionID).Scan(&total, &owned)
	return
}

func (d *DB) DropAll() error {
	for _, t := range []string{"collection_packs","collections","cards","sets","packs","_meta"} {
		if _, err := d.conn.Exec("DROP TABLE IF EXISTS "+t); err != nil { return err }
	}
	return nil
}

var _ = sql.ErrNoRows
