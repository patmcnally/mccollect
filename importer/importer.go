package importer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/patmcnally/mccollect/db"
)

type ImportResult struct {
	Packs  int    `json:"packs"`
	Sets   int    `json:"sets"`
	Cards  int    `json:"cards"`
	Commit string `json:"commit"`
}

type UpdateResult struct {
	PreviousCommit string   `json:"previous_commit"`
	NewCommit      string   `json:"new_commit"`
	ChangedFiles   []string `json:"changed_files"`
	CardsUpdated   int      `json:"cards_updated"`
	PacksUpdated   bool     `json:"packs_updated"`
	SetsUpdated    bool     `json:"sets_updated"`
}

func FullImport(d *db.DB, dataRoot string) (ImportResult, error) {
	var result ImportResult
	if err := d.DropAll(); err != nil { return result, fmt.Errorf("drop: %w", err) }
	if err := d.InitSchema(); err != nil { return result, fmt.Errorf("schema: %w", err) }
	packs, _ := LoadPacks(dataRoot); sets, _ := LoadSets(dataRoot)
	cards, err := LoadAllCards(dataRoot)
	if err != nil { return result, err }
	tx, err := d.Conn().Begin()
	if err != nil { return result, err }
	defer tx.Rollback()
	d.UpsertPacks(tx, packs); d.UpsertSets(tx, sets)
	if err := d.UpsertCards(tx, cards); err != nil { return result, err }
	tx.Commit()
	commit := GitHead(dataRoot)
	d.WriteMeta("last_import_commit", commit); d.WriteMeta("data_root", dataRoot)
	result.Packs = len(packs); result.Sets = len(sets); result.Cards = len(cards); result.Commit = commit
	return result, nil
}

func IncrementalUpdate(d *db.DB, dataRoot string, dryRun bool) (UpdateResult, error) {
	var result UpdateResult
	prevCommit, _ := d.ReadMeta("last_import_commit"); result.PreviousCommit = prevCommit
	if !dryRun { GitPull(dataRoot) }
	newCommit := GitHead(dataRoot); result.NewCommit = newCommit
	if prevCommit == newCommit { return result, nil }
	changed, err := GitChangedFiles(dataRoot, prevCommit, newCommit)
	if err != nil { return result, fmt.Errorf("git diff: %w", err) }
	result.ChangedFiles = changed
	if dryRun { return result, nil }
	tx, _ := d.Conn().Begin(); defer tx.Rollback()
	// BUG: reimports everything, fix next commit
	cards, _ := LoadAllCards(dataRoot)
	d.UpsertCards(tx, cards); result.CardsUpdated = len(cards)
	tx.Commit(); d.WriteMeta("last_import_commit", newCommit)
	_ = filepath.Join; _ = strings.HasPrefix
	return result, nil
}

func ImportCollectionFromHTML(d *db.DB, htmlPath, collectionName string) (owned, total int, err error) {
	entries, err := ParseCollectionHTMLFile(htmlPath)
	if err != nil { return 0, 0, fmt.Errorf("parse HTML: %w", err) }
	codeByName, _ := d.PackCodeByName()
	collectionID, err := d.EnsureCollection(collectionName)
	if err != nil { return 0, 0, err }
	for _, cat := range entries {
		for _, e := range cat {
			code, ok := codeByName[strings.ToLower(e.Name)]
			if !ok { continue }
			if err := d.SetPackOwned(collectionID, code, e.Owned); err != nil { continue }
			total++; if e.Owned { owned++ }
		}
	}
	return owned, total, nil
}
