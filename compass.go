package compass

import (
	"errors"

	tea "charm.land/bubbletea/v2"
)

type Screen int

type Model struct {
	screenEntries map[Screen]tea.Model
	stack         []Screen
	topIndex      int
}

// New creates a new Model with the given map of screens.
func New(screenEntries map[Screen]tea.Model, initialScreen Screen) Model {
	if _, ok := screenEntries[initialScreen]; !ok {
		panic("compass: initialScreen cannot be nil")
	}

	stack := []Screen{initialScreen}
	return Model{
		screenEntries: screenEntries,
		stack:         stack,
	}
}

// Update forwards the message to the model at the top of the navigation stack.
// If the stack is empty, it is a no-op.
func (m Model) Update(message tea.Msg) (Model, tea.Cmd) {
	if m.topIndex == -1 {
		return m, nil
	}

	screen := m.stack[m.topIndex]
	model, cmd := m.screenEntries[screen].Update(message)
	m.screenEntries[screen] = model

	return m, cmd
}

// View renders the screen at the top of the navigation stack.
// Returns an empty view if the stack is empty.
func (m Model) View() tea.View {
	if m.topIndex == -1 {
		return tea.NewView("")
	}

	screen := m.stack[m.topIndex]
	return m.screenEntries[screen].View()
}

// Push navigates to the given screen.
func (m Model) Push(screen Screen, args tea.Msg) (Model, tea.Cmd) {
	registeredModel, ok := m.screenEntries[screen]
	if !ok {
		panic("compass: screen not found in the screen entries; considering registering it.")
	}

	m.topIndex++

	if m.topIndex < len(m.stack) {
		m.stack[m.topIndex] = screen
	} else {
		m.stack = append(m.stack, screen)
	}

	if args == nil {
		return m, nil
	}

	registeredModel, cmd := registeredModel.Update(args)
	m.screenEntries[screen] = registeredModel
	return m, cmd

}

// Pop navigates back to the last pushed screen.
func (m Model) Pop() Model {
	m, err := m.CanPop()
	if err != nil {
		panic(err)
	}

	return m
}

// CanPop removes the top screen from the stack.
// If pop is called on empty stack, an error is returned.
func (m Model) CanPop() (Model, error) {
	if m.topIndex == -1 {
		return m, errors.New("compass: pop called on an empty navigation stack")
	}

	m.topIndex--
	return m, nil
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
