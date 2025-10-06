package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/patmcnally/mccollect/db"
)

type statsModel struct {
	database     *db.DB
	collectionID int
	stats        []db.PackTypeStat
	owned        int
	total        int
}

func newStatsModel(d *db.DB, collectionID int) statsModel {
	m := statsModel{
		database:     d,
		collectionID: collectionID,
	}
	m.refresh()
	return m
}

func (m *statsModel) refresh() {
	m.owned, m.total, _ = m.database.CollectionStats(m.collectionID)
	m.stats, _ = m.database.CollectionStatsByType(m.collectionID)
}

func (m statsModel) Init() tea.Cmd {
	return nil
}

func (m statsModel) Update(msg tea.Msg) (statsModel, tea.Cmd) {
	switch msg.(type) {
	case packToggledMsg:
		m.refresh()
	}
	return m, nil
}

func (m statsModel) View() string {
	var b strings.Builder

	b.WriteString(headerStyle.Render("Collection Statistics"))
	b.WriteByte('\n')
	b.WriteByte('\n')

	b.WriteString(fmt.Sprintf("  Total: %d / %d packs owned\n\n", m.owned, m.total))

	for _, s := range m.stats {
		label := packTypeDisplayName(s.PackTypeCode)
		bar := progressBar(s.Owned, s.Total, 20)
		b.WriteString(fmt.Sprintf("  %-20s %s %d/%d\n", label, bar, s.Owned, s.Total))
	}

	return b.String()
}

func packTypeDisplayName(code string) string {
	for _, t := range typeOrder {
		if t.code == code {
			return t.label
		}
	}
	return code
}

func progressBar(owned, total, width int) string {
	if total == 0 {
		return strings.Repeat("░", width)
	}
	filled := owned * width / total
	return ownedStyle.Render(strings.Repeat("█", filled)) +
		unownedStyle.Render(strings.Repeat("░", width-filled))
}
