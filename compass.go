package compass

import (
	tea "charm.land/bubbletea/v2"
)

type Screen int

type Model struct {
	screenEntries map[Screen]tea.Model
	stack         []Screen
	topIndex      int

	options Options
}

// New creates a new Model with the given map of screens.
func New(screenEntries map[Screen]tea.Model, initialScreen Screen, opts ...OptionFunc) Model {
	if _, ok := screenEntries[initialScreen]; !ok {
		panic("compass: initialScreen entry cannot be nil")
	}

	options := Options{}
	for _, opt := range opts {
		opt(&options)
	}

	stack := []Screen{initialScreen}
	return Model{
		screenEntries: screenEntries,
		stack:         stack,
		options:       options,
	}
}

// Update forwards the message to the model at the top of the navigation stack.
func (m Model) Update(message tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch message := message.(type) {
	case tea.KeyPressMsg:
		// key binds to exit if the stack is empty
		if m.isStackEmpty() {
			switch message.String() {
			case "esc", "ctrl+c":
				return m, tea.Quit
			}
		}

	case NavigateBackMsg:
		m = m.Pop()
		if m.isStackEmpty() && m.options.AutoQuitOnEmpty {
			return m, tea.Quit
		}

		return m, nil

	case NavigateMsg:
		m, cmd = m.Push(message.To, message.Argument)
		return m, cmd
	}

	screen := m.stack[m.topIndex]
	model, cmd := m.screenEntries[screen].Update(message)
	m.screenEntries[screen] = model

	return m, cmd
}

// View renders the screen at the top of the navigation stack.
// Returns an empty view if the stack is empty.
func (m Model) View() tea.View {
	if m.isStackEmpty() {
		return m.options.FallbackView
	}

	screen := m.stack[m.topIndex]
	return m.screenEntries[screen].View()
}

// Push adds the given screen to the top of the navigation stack,
// calls Init() on it, and forwards args to it via Update.
func (m Model) Push(screen Screen, args tea.Msg) (Model, tea.Cmd) {
	registeredModel, ok := m.screenEntries[screen]
	if !ok {
		panic("compass: screen not found in the screen entries; consider registering it.")
	}

	m.topIndex++

	if m.topIndex < len(m.stack) {
		m.stack[m.topIndex] = screen
	} else {
		m.stack = append(m.stack, screen)
	}

	initCmd := registeredModel.Init()
	if args == nil {
		return m, initCmd
	}

	registeredModel, cmd := registeredModel.Update(args)
	m.screenEntries[screen] = registeredModel
	return m, tea.Sequence(initCmd, cmd)
}

// Pop navigates back to the last pushed screen.
func (m Model) Pop() Model {
	if m.isStackEmpty() {
		panic("compass: pop called on an empty navigation stack")
	}

	m.topIndex--
	return m
}

// RegisterScreen associates a model with a screen identifier.
func (m Model) RegisterScreen(screen Screen, model tea.Model) Model {
	if model == nil {
		panic("compass: model cannot be nil")
	}

	m.screenEntries[screen] = model
	return m
}

// Model returns the tea.Model associated with the given screen identifier.
func (m Model) Model(screen Screen) tea.Model {
	return m.screenEntries[screen]
}

// Broadcast sends a message to every registered screen.
func (m Model) Broadcast(message tea.Msg) (Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, len(m.screenEntries))

	for screen, model := range m.screenEntries {
		updatedModel, cmd := model.Update(message)
		m.screenEntries[screen] = updatedModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) isStackEmpty() bool {
	return m.topIndex == -1
}
