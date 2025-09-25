package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/patmcnally/mccollect/db"
	"github.com/patmcnally/mccollect/model"
)

// packGroup is a section header + its packs for display.
type packGroup struct {
	label string
	packs []model.PackOwnership
}

// packsModel handles the pack selection view.
type packsModel struct {
	db           *db.DB
	collectionID int
	groups       []packGroup
	flatIndex    []flatEntry // flattened for cursor navigation
	cursor       int
	width        int
	height       int
}

type flatEntry struct {
	groupIdx int
	packIdx  int // -1 = header row
}

func newPacksModel(d *db.DB, collectionID int) packsModel {
	m := packsModel{
		db:           d,
		collectionID: collectionID,
	}
	m.loadPacks()
	return m
}

var typeOrder = []struct {
	code  string
	label string
}{
	{"core", "Core"},
	{"hero", "Hero Packs"},
	{"scenario", "Scenario Packs"},
	{"story", "Campaigns"},
	{"encounter", "Encounter Packs"},
}

func (m *packsModel) loadPacks() {
	packs, err := m.db.ListPackOwnership(m.collectionID)
	if err != nil {
		return
	}

	grouped := make(map[string][]model.PackOwnership)
	for _, po := range packs {
		grouped[po.Pack.PackTypeCode] = append(grouped[po.Pack.PackTypeCode], po)
	}

	m.groups = nil
	m.flatIndex = nil
	for _, t := range typeOrder {
		entries, ok := grouped[t.code]
		if !ok {
			continue
		}
		gIdx := len(m.groups)
		m.groups = append(m.groups, packGroup{label: t.label, packs: entries})
		m.flatIndex = append(m.flatIndex, flatEntry{groupIdx: gIdx, packIdx: -1})
		for i := range entries {
			m.flatIndex = append(m.flatIndex, flatEntry{groupIdx: gIdx, packIdx: i})
		}
	}
}

type packToggledMsg struct {
	packCode string
	owned    bool
}

func (m packsModel) Init() tea.Cmd {
	return nil
}

func (m packsModel) Update(msg tea.Msg) (packsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.flatIndex) - 1
			}
			// Skip headers
			if m.flatIndex[m.cursor].packIdx == -1 {
				m.cursor--
				if m.cursor < 0 {
					m.cursor = len(m.flatIndex) - 1
				}
			}
		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.flatIndex) {
				m.cursor = 0
			}
			// Skip headers
			if m.flatIndex[m.cursor].packIdx == -1 {
				m.cursor++
				if m.cursor >= len(m.flatIndex) {
					m.cursor = 0
				}
			}
		case " ", "enter":
			entry := m.flatIndex[m.cursor]
			if entry.packIdx >= 0 {
				po := &m.groups[entry.groupIdx].packs[entry.packIdx]
				newOwned, err := m.db.TogglePackOwned(m.collectionID, po.Pack.Code)
				if err == nil {
					po.Owned = newOwned
					return m, func() tea.Msg {
						return packToggledMsg{packCode: po.Pack.Code, owned: newOwned}
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m packsModel) View() string {
	var b strings.Builder

	// Calculate visible range for scrolling
	visibleStart, visibleEnd := m.visibleRange()

	for i, entry := range m.flatIndex {
		if i < visibleStart || i >= visibleEnd {
			continue
		}

		if entry.packIdx == -1 {
			// Section header
			b.WriteString(headerStyle.Render(m.groups[entry.groupIdx].label))
			b.WriteByte('\n')
			continue
		}

		po := m.groups[entry.groupIdx].packs[entry.packIdx]
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("> ")
		}

		check := "[ ]"
		style := unownedStyle
		if po.Owned {
			check = "[x]"
			style = ownedStyle
		}

		line := fmt.Sprintf("%s %s %s", cursor, check, po.Pack.Name)
		b.WriteString(style.Render(line))
		b.WriteByte('\n')
	}

	return b.String()
}

func (m packsModel) visibleRange() (int, int) {
	maxLines := m.height - 6 // reserve space for header/footer
	if maxLines <= 0 || maxLines >= len(m.flatIndex) {
		return 0, len(m.flatIndex)
	}

	start := m.cursor - maxLines/2
	if start < 0 {
		start = 0
	}
	end := start + maxLines
	if end > len(m.flatIndex) {
		end = len(m.flatIndex)
		start = end - maxLines
	}
	return start, end
}
