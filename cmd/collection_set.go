package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/patmcnally/mccollect/db"
	"github.com/spf13/cobra"
)

var (
	setPack    string
	setOwned   bool
	setUnowned bool
)

var collectionSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set ownership status of a pack",
	RunE: func(cmd *cobra.Command, args []string) error {
		if setPack == "" {
			return fmt.Errorf("--pack is required")
		}
		if !cmd.Flags().Changed("owned") && !cmd.Flags().Changed("not-owned") {
			return fmt.Errorf("either --owned or --not-owned is required")
		}

		owned := setOwned && !setUnowned

		d, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer d.Close()

		colID, err := d.EnsureCollection(collectionName)
		if err != nil {
			return err
		}

		if err := d.SetPackOwned(colID, setPack, owned); err != nil {
			return err
		}

		if jsonOut {
			return json.NewEncoder(os.Stdout).Encode(map[string]any{
				"collection": collectionName,
				"pack":       setPack,
				"owned":      owned,
			})
		}

		status := "owned"
		if !owned {
			status = "not owned"
		}
		fmt.Printf("Set %s → %s in collection %q\n", setPack, status, collectionName)
		return nil
	},
}

func init() {
	collectionSetCmd.Flags().StringVar(&setPack, "pack", "", "pack code to update")
	collectionSetCmd.Flags().BoolVar(&setOwned, "owned", false, "mark pack as owned")
	collectionSetCmd.Flags().BoolVar(&setUnowned, "not-owned", false, "mark pack as not owned")
	collectionCmd.AddCommand(collectionSetCmd)
}
