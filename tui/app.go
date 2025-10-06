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

// App is the top-level Bubbletea model.
type App struct {
	db             *db.DB
	collectionName string
	collectionID   int
	currentView    view
	packs          packsModel
	stats          statsModel
	width          int
	height         int
	err            error
}

// NewApp creates a new TUI application.
func NewApp(d *db.DB, collectionName string) App {
	colID, err := d.EnsureCollection(collectionName)
	if err != nil {
		return App{err: err}
	}

	return App{
		db:             d,
		collectionName: collectionName,
		collectionID:   colID,
		packs:          newPacksModel(d, colID),
		stats:          newStatsModel(d, colID),
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "tab":
			if a.currentView == viewPacks {
				a.currentView = viewStats
				a.stats.refresh()
			} else {
				a.currentView = viewPacks
			}
			return a, nil
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	}

	// Forward to active sub-model
	var cmd tea.Cmd
	switch a.currentView {
	case viewPacks:
		a.packs, cmd = a.packs.Update(msg)
		// Also forward toggle events to stats
		if _, ok := msg.(packToggledMsg); ok {
			a.stats, _ = a.stats.Update(msg)
		}
	case viewStats:
		a.stats, cmd = a.stats.Update(msg)
	}

	return a, cmd
}

func (a App) View() string {
	if a.err != nil {
		return fmt.Sprintf("Error: %v\n", a.err)
	}

	title := titleStyle.Render("Marvel Champions Collection")

	// Tab bar
	packsTab := " Packs "
	statsTab := " Stats "
	if a.currentView == viewPacks {
		packsTab = cursorStyle.Render("[Packs]")
		statsTab = helpStyle.Render(" Stats ")
	} else {
		packsTab = helpStyle.Render(" Packs ")
		statsTab = cursorStyle.Render("[Stats]")
	}
	tabs := fmt.Sprintf("%s  %s", packsTab, statsTab)

	var content string
	switch a.currentView {
	case viewPacks:
		content = a.packs.View()
	case viewStats:
		content = a.stats.View()
	}

	help := helpStyle.Render("tab: switch view • space/enter: toggle • j/k: navigate • q: quit")

	return fmt.Sprintf("%s\n%s\n\n%s\n%s\n", title, tabs, content, help)
}
