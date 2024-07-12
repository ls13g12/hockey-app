package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState uint

const (
	sidebarView   sessionState = iota
	createView
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	modelStyle = lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type mainModel struct {
	state					sessionState
	sidebarModel 	SidebarModel
	createModel 	CreateModel
}

type SidebarModel interface {
	tea.Model
}

type CreateModel interface {
	tea.Model
}


func MainScreen(sidebarModel SidebarModel, createModel CreateModel) mainModel {
	return mainModel{
		state: sidebarView,
		sidebarModel: sidebarModel,
		createModel: createModel,
	}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}


func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.state {
		case sidebarView:
			switch msg := msg.(type) {
			case tea.KeyMsg:
					switch msg.Type {
						case tea.KeyCtrlC, tea.KeyEsc:
							home_screen := ScreenOne(globalDB)
							return RootScreen(globalDB).SwitchScreen(&home_screen)
						case tea.KeyRight:
							m.state = createView
						default:
							m.sidebarModel, cmd = m.sidebarModel.Update(msg)
							cmds = append(cmds, cmd)
						}
			}

		case createView:
			switch msg := msg.(type) {
			case tea.KeyMsg:
					switch msg.Type {
						case tea.KeyCtrlC, tea.KeyEsc:
							m.state = sidebarView
						default:
							m.createModel, cmd = m.createModel.Update(msg)
							cmds = append(cmds, cmd)
						}
			}
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	if m.state == sidebarView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(m.sidebarModel.View()), modelStyle.Render(m.createModel.View()))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(m.sidebarModel.View()), focusedModelStyle.Render(m.createModel.View()))
	}

	return s
}

