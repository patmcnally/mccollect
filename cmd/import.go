package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/patmcnally/mccollect/importer"
	"github.com/spf13/cobra"
)

var dataPath string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import card data from marvelsdb-json-data",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dataPath == "" {
			return fmt.Errorf("--data is required")
		}
		packs, err := importer.LoadPacks(dataPath)
		if err != nil {
			return err
		}
		sets, err := importer.LoadSets(dataPath)
		if err != nil {
			return err
		}
		fmt.Printf("loaded %d packs, %d sets from %s\n", len(packs), len(sets), filepath.Base(dataPath))
		_ = json.Marshal // suppress import
		_ = os.Stdout
		return nil
	},
}

func init() {
	importCmd.Flags().StringVar(&dataPath, "data", "", "path to marvelsdb-json-data clone")
	importCmd.MarkFlagRequired("data")
	rootCmd.AddCommand(importCmd)
}
