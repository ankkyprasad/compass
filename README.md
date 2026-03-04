# Compass

A navigation stack for [Bubble Tea v2](https://charm.land/bubbletea) TUI applications.

Compass is inspired by flutter navigator and lets you manage multiple screens in your TUI app — push to navigate forward, pop to go back.

## Installation

```bash
go get github.com/ankkyprasad/compass/v2
```

## Usage

### 1. Define your screens

Use `compass.Screen` (an `int` alias) to identify each screen in your app:

```go
const (
    HomeScreen compass.Screen = iota
    DetailScreen
)
```

### 2. Create the navigator

Pass a map of screen identifiers to their `tea.Model` implementations and the inital starting screen.

```go
nav := compass.New(map[compass.Screen]tea.Model{
    HomeScreen:   NewHomeModel(),
    DetailScreen: NewDetailModel(),
}, HomeScreen)
```

### 3. Integrate it in your root model

```go
type RootModel struct {
    compass compass.Model
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.compass, cmd = m.compass.Update(msg)
	return m, cmd
}

func (m RootModel) View() tea.View {
    return m.compass.View()
}
```

### 4. Navigate between screens

Send messages from within any screen model to trigger navigation:

```go
func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Navigate forward to another screen. You can pass args in the second parameter.
			return m, compass.Push(DetailScreen, nil)

		case "esc":
			// Navigate back to the previous screen
			return m, compass.PopCmd()

		case "q":
			return m, tea.Quit
		}
	}
}
```

## License

Compass is licensed under the [MIT License](./LICENSE).
