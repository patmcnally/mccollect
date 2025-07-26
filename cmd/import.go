package cmd

import (
	"fmt"

	"github.com/patmcnally/mccollect/db"
	"github.com/patmcnally/mccollect/importer"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Full import of card data from marvelsdb-json-data",
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
		fmt.Printf("imported %d cards from %d packs (%d sets)\n", result.Cards, result.Packs, result.Sets)
		return nil
	},
}

func init() {
	importCmd.Flags().StringVar(&dataPath, "data", "", "path to marvelsdb-json-data clone")
	importCmd.MarkFlagRequired("data")
	rootCmd.AddCommand(importCmd)
}
