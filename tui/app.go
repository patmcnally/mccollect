package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/patmcnally/mccollect/db"
)

type App struct {
	db             *db.DB
	collectionName string
	collectionID   int
}

func NewApp(d *db.DB, collectionName string) App {
	colID, _ := d.EnsureCollection(collectionName)
	return App{db: d, collectionName: collectionName, collectionID: colID}
}

func (a App) Init() tea.Cmd { return nil }

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a App) View() string {
	return titleStyle.Render("Marvel Champions Collection") + "\n\npress q to quit\n"
}
