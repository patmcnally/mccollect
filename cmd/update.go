package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/patmcnally/mccollect/db"
	"github.com/patmcnally/mccollect/importer"
	"github.com/spf13/cobra"
)

var dryRun bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Incremental update from marvelsdb-json-data",
	Long:  "Git pull the data repo and re-import only changed pack files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if dataPath == "" {
			return fmt.Errorf("--data is required")
		}

		d, err := db.Open(dbPath)
		if err != nil {
			return err
		}
		defer d.Close()

		result, err := importer.IncrementalUpdate(d, dataPath, dryRun)
		if err != nil {
			return err
		}

		if jsonOut {
			return json.NewEncoder(os.Stdout).Encode(result)
		}

		if result.PreviousCommit == result.NewCommit {
			fmt.Println("Already up to date.")
			return nil
		}

		if dryRun {
			fmt.Printf("Would update from %s to %s\n", result.PreviousCommit[:8], result.NewCommit[:8])
			fmt.Printf("Changed files: %d\n", len(result.ChangedFiles))
			for _, f := range result.ChangedFiles {
				fmt.Printf("  %s\n", f)
			}
			return nil
		}

		fmt.Printf("Updated %s → %s\n", result.PreviousCommit[:8], result.NewCommit[:8])
		fmt.Printf("Changed files: %d, cards updated: %d\n", len(result.ChangedFiles), result.CardsUpdated)
		return nil
	},
}

func init() {
	updateCmd.Flags().StringVar(&dataPath, "data", "", "path to marvelsdb-json-data clone")
	updateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without applying")
	updateCmd.MarkFlagRequired("data")
	rootCmd.AddCommand(updateCmd)
}
