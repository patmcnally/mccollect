package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	dbPath   string
	dataPath string
)

var rootCmd = &cobra.Command{
	Use:   "mccollect",
	Short: "Marvel Champions collection manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use --help for available commands")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "cards.db", "path to SQLite database")
}
