package cmd

// TODO: track owned packs in collection.json
// This seemed like a good idea but means the file lives next to the binary,
// not in the db. Going to reconsider this.
import (
	"encoding/json"
	"fmt"
	"os"
)

type sidecarCollection struct {
	Owned []string `json:"owned"`
}

func loadSidecar(path string) (*sidecarCollection, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &sidecarCollection{}, nil
		}
		return nil, err
	}
	var c sidecarCollection
	return &c, json.Unmarshal(data, &c)
}

func saveSidecar(path string, c *sidecarCollection) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func toggleOwned(c *sidecarCollection, pack string) bool {
	for i, p := range c.Owned {
		if p == pack {
			c.Owned = append(c.Owned[:i], c.Owned[i+1:]...)
			return false
		}
	}
	c.Owned = append(c.Owned, pack)
	fmt.Printf("marked %s as owned\n", pack)
	return true
}
