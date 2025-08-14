package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/patmcnally/mccollect/db"
	"github.com/patmcnally/mccollect/importer"
	"github.com/spf13/cobra"
)

var collectionImportHTMLCmd = &cobra.Command{
	Use:   "import-html FILE",
	Short: "Import collection from saved marvelcdb.com HTML page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		htmlPath := args[0]

		d, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer d.Close()

		owned, total, err := importer.ImportCollectionFromHTML(d, htmlPath, collectionName)
		if err != nil {
			return err
		}

		if jsonOut {
			return json.NewEncoder(os.Stdout).Encode(map[string]any{
				"collection": collectionName,
				"owned":      owned,
				"total":      total,
			})
		}

		fmt.Printf("Collection %q: %d/%d packs owned\n", collectionName, owned, total)
		return nil
	},
}

func init() {
	collectionCmd.AddCommand(collectionImportHTMLCmd)
}
