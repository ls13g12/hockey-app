package tui

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ls13g12/hockey-app/src/pkg/common"
	"go.mongodb.org/mongo-driver/mongo"
)


func NewTuiApp(cfg common.TuiAppConfig, logger *slog.Logger, db *mongo.Database) {
	p := tea.NewProgram(RootScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type rootScreenModel struct {
	model  tea.Model    // this will hold the current screen model
}

func RootScreen() rootScreenModel {
	var rootModel tea.Model

	screen_one := ScreenOne()
	rootModel = &screen_one

	return rootScreenModel{
			model: rootModel,
	}
}

func (m rootScreenModel) Init() tea.Cmd {
	return m.model.Init()    // rest methods are just wrappers for the model's methods
}

func (m rootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m rootScreenModel) View() string {
	return m.model.View()
}

// this is the switcher which will switch between screens
func (m rootScreenModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	m.model = model
	return m.model, m.model.Init()    // must return .Init() to initialize the screen (and here the magic happens)
}
