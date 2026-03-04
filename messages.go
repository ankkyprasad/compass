package compass

import tea "charm.land/bubbletea/v2"

type NavigateMsg struct {
	To       Screen
	Argument tea.Msg
}

type NavigateBackMsg struct{}
