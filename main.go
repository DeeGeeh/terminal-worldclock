package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	labelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	timeStyle  = lipgloss.NewStyle().Bold(true)
	dateStyle  = lipgloss.NewStyle().Faint(true)
	hintStyle  = lipgloss.NewStyle().Faint(true)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Align(lipgloss.Center)
	focusedStyle = panelStyle.BorderForeground(lipgloss.Color("39"))
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

type zoneItem string

func (z zoneItem) FilterValue() string { return string(z) }
func (z zoneItem) Title() string       { return string(z) }
func (z zoneItem) Description() string { return "" }

type model struct {
	cfg     config
	locs    [2]*time.Location
	focus   int
	picking bool
	picker  list.Model
	width   int
	height  int
}

func initialModel() model {
	cfg := loadConfig()
	m := model{cfg: cfg}
	for i, z := range cfg.Zones {
		m.locs[i] = mustLoad(z)
	}

	names := zoneNames()
	items := make([]list.Item, len(names))
	for i, n := range names {
		items[i] = zoneItem(n)
	}
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	l := list.New(items, d, 30, 20)
	l.Title = "select timezone"
	m.picker = l
	return m
}

func mustLoad(name string) *time.Location {
	if loc, err := loadZone(name); err == nil {
		return loc
	}
	return time.UTC
}

func (m model) Init() tea.Cmd { return tick() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.picker.SetSize(msg.Width, msg.Height-1)
		return m, nil
	case tickMsg:
		return m, tick()
	case tea.KeyMsg:
		if m.picking {
			return m.updatePicker(msg)
		}
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab", "left", "right", "h", "l":
			m.focus = 1 - m.focus
		case " ", "enter":
			m.picking = true
			m.picker.ResetFilter()
		}
	}
	return m, nil
}

func (m model) updatePicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if it, ok := m.picker.SelectedItem().(zoneItem); ok {
			name := string(it)
			m.locs[m.focus] = mustLoad(name)
			m.cfg.Zones[m.focus] = name
			m.cfg.save()
		}
		m.picking = false
		return m, nil
	case "esc":
		if m.picker.FilterState() != list.Filtering {
			m.picking = false
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.picker, cmd = m.picker.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	if m.picking {
		return m.picker.View()
	}
	now := time.Now()
	row := lipgloss.JoinHorizontal(lipgloss.Top, m.panel(0, now), m.panel(1, now))
	hint := hintStyle.Render("tab switch · space change zone · q quit")
	return lipgloss.JoinVertical(lipgloss.Center, row, hint)
}

func (m model) panel(i int, now time.Time) string {
	contentW := m.width/2 - 2
	contentH := m.height - 4
	t := now.In(m.locs[i])

	body := lipgloss.JoinVertical(lipgloss.Center,
		labelStyle.Render(m.cfg.Zones[i]),
		renderClock(t, contentW, contentH-3),
		timeStyle.Render(t.Format("15:04")),
		dateStyle.Render(t.Format("Mon, 2 Jan")),
	)

	style := panelStyle
	if i == m.focus {
		style = focusedStyle
	}
	return style.Width(contentW).Height(contentH).Render(body)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
