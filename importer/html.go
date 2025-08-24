package importer

import (
	"os"
	"regexp"
	"strings"
)

var (
	blockSplit = regexp.MustCompile(`<div class="col-md-3 col-sm-4 cycle">`)
	catLabel   = regexp.MustCompile(`<label>([^<]+)</label>`)
	packEntry  = regexp.MustCompile(`<label([^>]+)>([^<]+)</label>`)
	dataID     = regexp.MustCompile(`\bdata-id="(\d+)"`)
	classAttr  = regexp.MustCompile(`\bclass="([^"]*)"`)
)

type CollectionEntry struct {
	Name  string
	Owned bool
}

func ParseCollectionHTML(html string) map[string][]CollectionEntry {
	categoryMap := map[string]string{
		"Core":          "core",
		"Scenario Pack": "scenario_packs",
		"Hero Pack":     "hero_packs",
		"Campaign":      "campaigns",
		"Encounter":     "encounter_packs",
	}
	result := make(map[string][]CollectionEntry)
	for _, v := range categoryMap {
		result[v] = nil
	}
	blocks := blockSplit.Split(html, -1)
	if len(blocks) < 2 {
		return result
	}
	for _, block := range blocks[1:] {
		catMatch := catLabel.FindStringSubmatch(block)
		if catMatch == nil {
			continue
		}
		catKey, ok := categoryMap[strings.TrimSpace(catMatch[1])]
		if !ok {
			continue
		}
		for _, m := range packEntry.FindAllStringSubmatch(block, -1) {
			attrs := m[1]
			if !dataID.MatchString(attrs) {
				continue
			}
			name := strings.TrimSpace(m[2])
			classMatch := classAttr.FindStringSubmatch(attrs)
			owned := classMatch != nil && strings.Contains(classMatch[1], "active")
			result[catKey] = append(result[catKey], CollectionEntry{Name: name, Owned: owned})
		}
	}
	return result
}

func ParseCollectionHTMLFile(path string) (map[string][]CollectionEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseCollectionHTML(string(data)), nil
}
