package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/patmcnally/mccollect/db"
)

type statsModel struct {
	database     *db.DB
	collectionID int
	owned        int
	total        int
}

func newStatsModel(d *db.DB, collectionID int) statsModel {
	m := statsModel{database: d, collectionID: collectionID}
	m.refresh()
	return m
}

func (m *statsModel) refresh() {
	m.owned, m.total, _ = m.database.CollectionStats(m.collectionID)
}

func (m statsModel) Init() tea.Cmd { return nil }
func (m statsModel) Update(msg tea.Msg) (statsModel, tea.Cmd) { return m, nil }
func (m statsModel) View() string {
	return headerStyle.Render("Stats") + "\n\n  (coming soon)\n"
}
