package importer

import (
	"fmt"

	"github.com/patmcnally/mccollect/db"
)

type ImportResult struct {
	Packs  int    `json:"packs"`
	Sets   int    `json:"sets"`
	Cards  int    `json:"cards"`
	Commit string `json:"commit"`
}

func FullImport(d *db.DB, dataRoot string) (ImportResult, error) {
	var result ImportResult

	if err := d.DropAll(); err != nil {
		return result, fmt.Errorf("drop tables: %w", err)
	}
	if err := d.InitSchema(); err != nil {
		return result, fmt.Errorf("init schema: %w", err)
	}

	packs, err := LoadPacks(dataRoot)
	if err != nil {
		return result, fmt.Errorf("load packs: %w", err)
	}
	sets, err := LoadSets(dataRoot)
	if err != nil {
		return result, fmt.Errorf("load sets: %w", err)
	}
	cards, err := LoadAllCards(dataRoot)
	if err != nil {
		return result, fmt.Errorf("load cards: %w", err)
	}

	tx, err := d.Conn().Begin()
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	if err := d.UpsertPacks(tx, packs); err != nil {
		return result, err
	}
	if err := d.UpsertSets(tx, sets); err != nil {
		return result, err
	}
	if err := d.UpsertCards(tx, cards); err != nil {
		return result, err
	}
	if err := tx.Commit(); err != nil {
		return result, err
	}

	result.Packs = len(packs)
	result.Sets = len(sets)
	result.Cards = len(cards)
	return result, nil
}
