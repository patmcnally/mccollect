package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/patmcnally/mccollect/db"
	"github.com/patmcnally/mccollect/tui"
)

var tuiCollectionName string

func init() {
	rootCmd.Flags().StringVar(&tuiCollectionName, "name", "default", "collection name for TUI")
}

func runTUI() error {
	d, err := db.Open(dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer d.Close()

	if err := d.InitSchema(); err != nil {
		return fmt.Errorf("init schema: %w", err)
	}

	app := tui.NewApp(d, tuiCollectionName)
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
