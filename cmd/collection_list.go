package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/patmcnally/mccollect/db"
	"github.com/patmcnally/mccollect/model"
	"github.com/spf13/cobra"
)

var collectionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List packs with ownership status",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer d.Close()

		col, err := d.GetCollection(collectionName)
		if err != nil {
			return fmt.Errorf("collection %q not found (run import-html first)", collectionName)
		}

		packs, err := d.ListPackOwnership(col.ID)
		if err != nil {
			return err
		}

		if jsonOut {
			return json.NewEncoder(os.Stdout).Encode(packs)
		}

		// Group by pack type
		grouped := make(map[string][]model.PackOwnership)
		typeOrder := []string{"core", "hero", "scenario", "story", "encounter"}
		for _, po := range packs {
			grouped[po.Pack.PackTypeCode] = append(grouped[po.Pack.PackTypeCode], po)
		}

		owned, total := 0, 0
		for _, t := range typeOrder {
			entries, ok := grouped[t]
			if !ok {
				continue
			}
			fmt.Printf("\n%s:\n", packTypeLabel(t))
			for _, po := range entries {
				mark := "  "
				if po.Owned {
					mark = "* "
					owned++
				}
				total++
				fmt.Printf("  %s%s\n", mark, po.Pack.Name)
			}
		}
		fmt.Printf("\n%d/%d packs owned\n", owned, total)
		return nil
	},
}

func packTypeLabel(code string) string {
	switch code {
	case "core":
		return "Core"
	case "hero":
		return "Hero Packs"
	case "scenario":
		return "Scenario Packs"
	case "story":
		return "Campaigns"
	case "encounter":
		return "Encounter Packs"
	default:
		return code
	}
}

func init() {
	collectionCmd.AddCommand(collectionListCmd)
}
