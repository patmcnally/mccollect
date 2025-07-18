package importer

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/patmcnally/mccollect/model"
)

// rawPack matches the JSON structure in packs.json.
type rawPack struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	CgdbID       *int    `json:"cgdb_id"`
	OctgnID      *string `json:"octgn_id"`
	DateRelease  *string `json:"date_release"`
	PackTypeCode string  `json:"pack_type_code"`
	Position     *int    `json:"position"`
	Size         *int    `json:"size"`
}

// LoadPacks reads packs.json from the data root.
func LoadPacks(dataRoot string) ([]model.Pack, error) {
	data, err := os.ReadFile(filepath.Join(dataRoot, "packs.json"))
	if err != nil {
		return nil, err
	}
	var raw []rawPack
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	packs := make([]model.Pack, len(raw))
	for i, r := range raw {
		packs[i] = model.Pack{
			Code:         r.Code,
			Name:         r.Name,
			CgdbID:       r.CgdbID,
			OctgnID:      r.OctgnID,
			DateRelease:  r.DateRelease,
			PackTypeCode: r.PackTypeCode,
			Position:     r.Position,
			Size:         r.Size,
		}
	}
	return packs, nil
}
