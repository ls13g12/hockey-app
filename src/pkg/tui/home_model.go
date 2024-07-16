package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type homeModel struct {
	choices 						[]choice
	resourceCursor 			int
	resourceSelected		bool
	actionCursor				int
}

type choice struct {
	resource string
	actions []string
}


func NewHomeModel() homeModel {
	return homeModel{
		choices: []choice{
			{
				resource: "players",
				actions: []string{
					"list",
					"delete",
					"update",
				},
			},
			{
				resource: "matches",
				actions: []string{
					"list",
					"delete",
					"update",
				},
			},
		},
		resourceCursor: 0,
		actionCursor: 0,
	}
}

func (m homeModel) Init() tea.Cmd {
	return nil
}


func (m homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.resourceSelected {
		switch msg := msg.(type) {
		case tea.KeyMsg:
				switch msg.Type {
						case tea.KeyCtrlC:
								m.resourceSelected = false
	
						case tea.KeyUp:
								if m.actionCursor > 0 {
										m.actionCursor--
								}
	
						case tea.KeyDown:
								if m.actionCursor < len(m.choices[m.resourceCursor].actions)-1 {
										m.actionCursor++
								}
	
						case tea.KeyEnter, tea.KeyBackspace:
								switch m.actionCursor {
									case 0:
										playerScreenModel := NewListModel(m.choices[m.resourceCursor].resource)
										return RootScreen().SwitchScreen(playerScreenModel)
								}
				}
		}
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
				switch msg.Type {
						case tea.KeyCtrlC:
								return m, tea.Quit
	
						case tea.KeyUp:
								if m.resourceCursor > 0 {
										m.resourceCursor--
								}
	
						case tea.KeyDown:
								if m.resourceCursor < len(m.choices)-1 {
										m.resourceCursor++
								}
	
						case tea.KeyEnter, tea.KeyBackspace:
								m.resourceSelected = true
				}
		}
	}

	return m, nil
}

func (m homeModel) View() string {
	var s string

	var resourceColumn string
	for i, choice := range m.choices {
		resourceCursor := " "
		if m.resourceCursor == i {
				resourceCursor = ">" 
		}
		resourceColumn += fmt.Sprintf("%s %s\n", resourceCursor, choice.resource)
	}

	var actionColumn string
	for i, action := range m.choices[m.resourceCursor].actions {
		actionCursor := " "
		if m.actionCursor == i {
			actionCursor = ">" 
		}
		actionColumn += fmt.Sprintf("%s %s\n", actionCursor, action)
	}

	if m.resourceSelected {
		s += lipgloss.JoinHorizontal(lipgloss.Top, listBox.Render(resourceColumn), listBox.Render(actionColumn))
	} else {
		s += listBox.Render(resourceColumn)
	}

	return s
}
