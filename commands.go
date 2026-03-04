package compass

import tea "charm.land/bubbletea/v2"

// Make sure compass.Update() is called in your root model for these commands to work.

// PushCmd returns a tea.Cmd that pushes the given screen onto the navigation stack.
// Use the argument parameter to forward a message to the screen on initialization.
func PushCmd(to Screen, argument tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return NavigateMsg{
			To:       to,
			Argument: argument,
		}
	}
}

// PopCmd returns a tea.Cmd that removes the current screen from the navigation stack.
func PopCmd() tea.Cmd {
	return func() tea.Msg {
		return NavigateBackMsg{}
	}
}
