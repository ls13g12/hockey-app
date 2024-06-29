package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type screenOneModel struct {
	choices []string
}

func ScreenOne() screenOneModel {
	return screenOneModel{
		choices: []string{"home1", "home2"},
	}
}

func (m screenOneModel) Init() tea.Cmd {
	return nil
}


func (m screenOneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC:
					return m, tea.Quit
			default:
					screen_two := ScreenTwo()
					return RootScreen().SwitchScreen(&screen_two)
			}
	default:
			return m, nil
	}
}

func (m screenOneModel) View() string {
	str := "This is the first screen. Press any key to switch to the second screen."
	return str
}
