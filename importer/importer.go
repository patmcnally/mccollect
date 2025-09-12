package importer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/patmcnally/mccollect/db"
)

// ImportResult summarizes what was imported.
type ImportResult struct {
	Packs int `json:"packs"`
	Sets  int `json:"sets"`
	Cards int `json:"cards"`
	Commit string `json:"commit"`
}

// FullImport wipes and rebuilds the database from a marvelsdb-json-data clone.
func FullImport(d *db.DB, dataRoot string) (ImportResult, error) {
	var result ImportResult

	// Drop and recreate
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
		return result, fmt.Errorf("upsert packs: %w", err)
	}
	if err := d.UpsertSets(tx, sets); err != nil {
		return result, fmt.Errorf("upsert sets: %w", err)
	}
	if err := d.UpsertCards(tx, cards); err != nil {
		return result, fmt.Errorf("upsert cards: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return result, err
	}

	commit := GitHead(dataRoot)
	d.WriteMeta("last_import_commit", commit)
	d.WriteMeta("data_root", dataRoot)

	result.Packs = len(packs)
	result.Sets = len(sets)
	result.Cards = len(cards)
	result.Commit = commit
	return result, nil
}

// UpdateResult summarizes an incremental update.
type UpdateResult struct {
	PreviousCommit string   `json:"previous_commit"`
	NewCommit      string   `json:"new_commit"`
	ChangedFiles   []string `json:"changed_files"`
	CardsUpdated   int      `json:"cards_updated"`
	PacksUpdated   bool     `json:"packs_updated"`
	SetsUpdated    bool     `json:"sets_updated"`
}

// IncrementalUpdate pulls the data repo and re-imports only changed files.
func IncrementalUpdate(d *db.DB, dataRoot string, dryRun bool) (UpdateResult, error) {
	var result UpdateResult

	prevCommit, _ := d.ReadMeta("last_import_commit")
	result.PreviousCommit = prevCommit

	if !dryRun {
		if err := GitPull(dataRoot); err != nil {
			return result, fmt.Errorf("git pull: %w", err)
		}
	}

	newCommit := GitHead(dataRoot)
	result.NewCommit = newCommit

	if prevCommit == newCommit {
		return result, nil // nothing to do
	}

	changed, err := GitChangedFiles(dataRoot, prevCommit, newCommit)
	if err != nil {
		return result, fmt.Errorf("git diff: %w", err)
	}
	result.ChangedFiles = changed

	if dryRun {
		return result, nil
	}

	tx, err := d.Conn().Begin()
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	for _, f := range changed {
		switch {
		case f == "packs.json":
			packs, err := LoadPacks(dataRoot)
			if err != nil {
				return result, err
			}
			if err := d.UpsertPacks(tx, packs); err != nil {
				return result, err
			}
			result.PacksUpdated = true

		case f == "sets.json":
			sets, err := LoadSets(dataRoot)
			if err != nil {
				return result, err
			}
			if err := d.UpsertSets(tx, sets); err != nil {
				return result, err
			}
			result.SetsUpdated = true

		case strings.HasPrefix(f, "pack/") && strings.HasSuffix(f, ".json"):
			path := filepath.Join(dataRoot, f)
			cards, err := LoadPackFile(path)
			if err != nil {
				return result, fmt.Errorf("load %s: %w", f, err)
			}
			if err := d.UpsertCards(tx, cards); err != nil {
				return result, fmt.Errorf("upsert %s: %w", f, err)
			}
			result.CardsUpdated += len(cards)
		}
	}

	if err := tx.Commit(); err != nil {
		return result, err
	}

	d.WriteMeta("last_import_commit", newCommit)
	return result, nil
}

// ImportCollectionFromHTML parses HTML and imports ownership into the database.
func ImportCollectionFromHTML(d *db.DB, htmlPath, collectionName string) (owned, total int, err error) {
	entries, err := ParseCollectionHTMLFile(htmlPath)
	if err != nil {
		return 0, 0, fmt.Errorf("parse HTML: %w", err)
	}

	codeByName, err := d.PackCodeByName()
	if err != nil {
		return 0, 0, fmt.Errorf("load pack names: %w", err)
	}

	collectionID, err := d.EnsureCollection(collectionName)
	if err != nil {
		return 0, 0, err
	}

	packOwnership := make(map[string]bool)
	for _, catEntries := range entries {
		for _, e := range catEntries {
			code, ok := codeByName[strings.ToLower(e.Name)]
			if !ok {
				continue
			}
			packOwnership[code] = e.Owned
			total++
			if e.Owned {
				owned++
			}
		}
	}

	if err := d.ImportCollectionBulk(collectionID, packOwnership); err != nil {
		return 0, 0, err
	}

	return owned, total, nil
}
