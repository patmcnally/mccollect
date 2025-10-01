package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/patmcnally/mccollect/db"
)

type view int

const (
	viewPacks view = iota
	viewStats
)

type App struct {
	db             *db.DB
	collectionName string
	collectionID   int
	currentView    view
	packs          packsModel
	stats          statsModel
	err            error
}

func NewApp(d *db.DB, collectionName string) App {
	colID, err := d.EnsureCollection(collectionName)
	if err != nil { return App{err: err} }
	return App{
		db: d, collectionName: collectionName, collectionID: colID,
		packs: newPacksModel(d, colID), stats: newStatsModel(d, colID),
	}
}

func (a App) Init() tea.Cmd { return nil }

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "tab":
			if a.currentView == viewPacks { a.currentView = viewStats; a.stats.refresh() } else { a.currentView = viewPacks }
			return a, nil
		}
	}
	var cmd tea.Cmd
	if a.currentView == viewPacks { a.packs, cmd = a.packs.Update(msg) } else { a.stats, cmd = a.stats.Update(msg) }
	return a, cmd
}

func (a App) View() string {
	if a.err != nil { return fmt.Sprintf("Error: %v\n", a.err) }
	var content string
	if a.currentView == viewPacks { content = a.packs.View() } else { content = a.stats.View() }
	return titleStyle.Render("Marvel Champions Collection") + "\n\n" + content + helpStyle.Render("tab: switch • space/enter: toggle • j/k: move • q: quit") + "\n"
}
