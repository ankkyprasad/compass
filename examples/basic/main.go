package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/ankkyprasad/compass/v2"
)

// Screen identifiers.
const (
	HomeScreen compass.Screen = iota
	DetailScreen
	AnotherDetailScreen
)

// --- Navigation messages ---

type NavigateMsg struct{ To compass.Screen }
type NavigateBackMsg struct{}

type RootModel struct {
	compass compass.Model
}

func NewRootModel() RootModel {
	compass := compass.New(map[compass.Screen]tea.Model{
		HomeScreen:   NewHomeModel(),
		DetailScreen: NewDetailModel(),
	}, HomeScreen)

	return RootModel{compass: compass}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NavigateMsg:
		var cmd tea.Cmd
		m.compass, cmd = m.compass.Push(msg.To, nil)
		return m, cmd

	case NavigateBackMsg:
		m.compass = m.compass.Pop()
		return m, nil
	}

	// Forward everything else to the active screen.
	var cmd tea.Cmd
	m.compass, cmd = m.compass.Update(msg)
	return m, cmd
}

func (m RootModel) View() tea.View {
	return m.compass.View()
}

type HomeModel struct {
	choices  []string
	cursor   int
	selected int
}

func NewHomeModel() HomeModel {
	return HomeModel{
		choices: []string{"View Details", "Quit"},
	}
}

func (m HomeModel) Init() tea.Cmd { return nil }

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.cursor
			switch m.cursor {
			case 0:
				return m, func() tea.Msg { return NavigateMsg{To: DetailScreen} }
			case 1:
				return m, tea.Quit
			}

		case "esc":
			return m, func() tea.Msg { return NavigateBackMsg{} }
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m HomeModel) View() tea.View {
	var s strings.Builder
	s.WriteString("Welcome! Choose an option:\n\n")
	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}
		s.WriteString(cursor + choice + "\n")
	}
	s.WriteString("\nPress q to quit.\n")
	return tea.NewView(s.String())
}

type DetailModel struct{}

func NewDetailModel() DetailModel { return DetailModel{} }

func (m DetailModel) Init() tea.Cmd { return nil }

func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			return m, func() tea.Msg { return NavigateBackMsg{} }
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m DetailModel) View() tea.View {
	return tea.NewView("📄 Detail Screen\n\nThis is the detail view.\n\nPress esc to go back, q to quit.\n")
}

func main() {
	p := tea.NewProgram(NewRootModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
