package db

import (
	"database/sql"
	"fmt"

	"github.com/patmcnally/mccollect/model"
)

// EnsureCollection creates a collection by name if it doesn't exist, and returns its ID.
func (d *DB) EnsureCollection(name string) (int, error) {
	_, err := d.conn.Exec(
		`INSERT INTO collections (name) VALUES (?)
		 ON CONFLICT(name) DO UPDATE SET updated_at = datetime('now')`,
		name,
	)
	if err != nil {
		return 0, fmt.Errorf("ensure collection %q: %w", name, err)
	}
	var id int
	err = d.conn.QueryRow("SELECT id FROM collections WHERE name = ?", name).Scan(&id)
	return id, err
}

// SetPackOwned sets the ownership status of a single pack in a collection.
func (d *DB) SetPackOwned(collectionID int, packCode string, owned bool) error {
	ownedInt := 0
	if owned {
		ownedInt = 1
	}
	_, err := d.conn.Exec(
		`INSERT INTO collection_packs (collection_id, pack_code, owned) VALUES (?, ?, ?)
		 ON CONFLICT(collection_id, pack_code) DO UPDATE SET owned = excluded.owned`,
		collectionID, packCode, ownedInt,
	)
	return err
}

// ImportCollectionBulk replaces all pack ownership for a collection.
// entries maps pack_code to owned status.
func (d *DB) ImportCollectionBulk(collectionID int, entries map[string]bool) error {
	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM collection_packs WHERE collection_id = ?", collectionID); err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO collection_packs (collection_id, pack_code, owned) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for packCode, owned := range entries {
		ownedInt := 0
		if owned {
			ownedInt = 1
		}
		if _, err := stmt.Exec(collectionID, packCode, ownedInt); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ListPackOwnership returns all packs with their ownership status for a collection.
func (d *DB) ListPackOwnership(collectionID int) ([]model.PackOwnership, error) {
	rows, err := d.conn.Query(`
		SELECT p.code, p.name, p.cgdb_id, p.octgn_id, p.date_release, p.pack_type_code, p.position, p.size,
		       COALESCE(cp.owned, 0)
		FROM packs p
		LEFT JOIN collection_packs cp ON p.code = cp.pack_code AND cp.collection_id = ?
		ORDER BY p.date_release, p.position`, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.PackOwnership
	for rows.Next() {
		var po model.PackOwnership
		var owned int
		if err := rows.Scan(
			&po.Pack.Code, &po.Pack.Name, &po.Pack.CgdbID, &po.Pack.OctgnID,
			&po.Pack.DateRelease, &po.Pack.PackTypeCode, &po.Pack.Position, &po.Pack.Size,
			&owned,
		); err != nil {
			return nil, err
		}
		po.Owned = owned == 1
		result = append(result, po)
	}
	return result, rows.Err()
}

// CollectionStats returns owned and total pack counts for a collection.
func (d *DB) CollectionStats(collectionID int) (owned, total int, err error) {
	err = d.conn.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(CASE WHEN cp.owned = 1 THEN 1 ELSE 0 END), 0)
		FROM packs p
		LEFT JOIN collection_packs cp ON p.code = cp.pack_code AND cp.collection_id = ?`,
		collectionID,
	).Scan(&total, &owned)
	return
}

// CollectionStatsByType returns owned and total counts grouped by pack_type_code.
func (d *DB) CollectionStatsByType(collectionID int) ([]PackTypeStat, error) {
	rows, err := d.conn.Query(`
		SELECT p.pack_type_code,
		       COUNT(*),
		       COALESCE(SUM(CASE WHEN cp.owned = 1 THEN 1 ELSE 0 END), 0)
		FROM packs p
		LEFT JOIN collection_packs cp ON p.code = cp.pack_code AND cp.collection_id = ?
		GROUP BY p.pack_type_code
		ORDER BY p.pack_type_code`, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []PackTypeStat
	for rows.Next() {
		var s PackTypeStat
		if err := rows.Scan(&s.PackTypeCode, &s.Total, &s.Owned); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

// GetCollection returns a collection by name, or sql.ErrNoRows if not found.
func (d *DB) GetCollection(name string) (model.Collection, error) {
	var c model.Collection
	err := d.conn.QueryRow(
		"SELECT id, name, description, created_at, updated_at FROM collections WHERE name = ?",
		name,
	).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

// PackTypeStat holds counts for a single pack type.
type PackTypeStat struct {
	PackTypeCode string `json:"pack_type_code"`
	Total        int    `json:"total"`
	Owned        int    `json:"owned"`
}

// TogglePackOwned flips the ownership status of a pack and returns the new state.
func (d *DB) TogglePackOwned(collectionID int, packCode string) (bool, error) {
	var current int
	err := d.conn.QueryRow(
		"SELECT COALESCE((SELECT owned FROM collection_packs WHERE collection_id = ? AND pack_code = ?), 0)",
		collectionID, packCode,
	).Scan(&current)
	if err != nil {
		return false, err
	}
	newOwned := current == 0
	if err := d.SetPackOwned(collectionID, packCode, newOwned); err != nil {
		return false, err
	}
	return newOwned, nil
}

// unused import guard
var _ = sql.ErrNoRows
