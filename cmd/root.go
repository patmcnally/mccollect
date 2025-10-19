package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	dbPath   string
	jsonOut  bool
	dataPath string
)

var rootCmd = &cobra.Command{
	Use:   "mccollect",
	Short: "Marvel Champions collection manager",
	Long:  "Manage your Marvel Champions: The Card Game collection — import card data, track pack ownership, and browse your collection.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTUI()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// defaultDBPath returns ~/.config/mccollect/cards.db, falling back to
// "references/cards.db" if the home directory cannot be determined.
func defaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "references/cards.db"
	}
	return filepath.Join(home, ".config", "mccollect", "cards.db")
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", defaultDBPath(), "path to SQLite database")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output JSON instead of human-readable text")
}
