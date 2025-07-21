package importer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/patmcnally/mccollect/model"
)

// LoadPackFile reads a single pack JSON file.
func LoadPackFile(path string) ([]model.Card, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw []map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	var cards []model.Card
	for _, r := range raw {
		name, _ := r["name"].(string)
		code, _ := r["code"].(string)
		typeCode, _ := r["type_code"].(string)
		packCode, _ := r["pack_code"].(string)
		if typeCode == "" || name == "" {
			continue
		}
		cards = append(cards, model.Card{
			Code:        code,
			Name:        name,
			TypeCode:    typeCode,
			FactionCode: func() string { s, _ := r["faction_code"].(string); return s }(),
			PackCode:    packCode,
		})
	}
	return cards, nil
}

// LoadAllCards reads all pack/*.json files.
func LoadAllCards(dataRoot string) ([]model.Card, error) {
	packDir := filepath.Join(dataRoot, "pack")
	entries, err := os.ReadDir(packDir)
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	var all []model.Card
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		cards, err := LoadPackFile(filepath.Join(packDir, e.Name()))
		if err != nil {
			return nil, err
		}
		all = append(all, cards...)
	}
	return all, nil
}
