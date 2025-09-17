package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/patmcnally/mccollect/db"
	"github.com/patmcnally/mccollect/importer"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Full import of card data from marvelsdb-json-data",
	Long:  "Wipe and rebuild the database from a local marvelsdb-json-data clone.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dataPath == "" {
			return fmt.Errorf("--data is required")
		}

		d, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer d.Close()

		result, err := importer.FullImport(d, dataPath)
		if err != nil {
			return err
		}

		if jsonOut {
			return json.NewEncoder(os.Stdout).Encode(result)
		}

		fmt.Printf("Imported %d cards from %d packs (%d sets) → %s\n", result.Cards, result.Packs, result.Sets, dbPath)
		fmt.Printf("Data commit: %s\n", result.Commit)
		return nil
	},
}

func init() {
	importCmd.Flags().StringVar(&dataPath, "data", "", "path to marvelsdb-json-data clone")
	importCmd.MarkFlagRequired("data")
	rootCmd.AddCommand(importCmd)
}
