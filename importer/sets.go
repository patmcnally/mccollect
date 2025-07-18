package importer

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/patmcnally/mccollect/model"
)

// rawSet matches the JSON structure in sets.json.
type rawSet struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	CardSetTypeCode string `json:"card_set_type_code"`
}

// LoadSets reads sets.json from the data root.
func LoadSets(dataRoot string) ([]model.Set, error) {
	data, err := os.ReadFile(filepath.Join(dataRoot, "sets.json"))
	if err != nil {
		return nil, err
	}
	var raw []rawSet
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	sets := make([]model.Set, len(raw))
	for i, r := range raw {
		sets[i] = model.Set{
			Code:            r.Code,
			Name:            r.Name,
			CardSetTypeCode: r.CardSetTypeCode,
		}
	}
	return sets, nil
}
