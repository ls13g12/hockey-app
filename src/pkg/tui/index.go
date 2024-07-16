package tui

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ls13g12/hockey-app/src/pkg/common"
	"go.mongodb.org/mongo-driver/mongo"
)

var globalDB *mongo.Database

func NewTuiApp(cfg common.TuiAppConfig, logger *slog.Logger, db *mongo.Database) {
	err := db.Client().Ping(context.TODO(), nil)
	if err != nil {
		fmt.Printf("Error pinging db: %v", err)
		os.Exit(1)
	}
	globalDB = db

	p := tea.NewProgram(RootScreen(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type rootScreenModel struct {
	model  tea.Model
}

func RootScreen() rootScreenModel {
	var rootModel tea.Model

	homeModel := NewHomeModel()
	rootModel = &homeModel

	return rootScreenModel{
			model: rootModel,
	}
}

func (m rootScreenModel) Init() tea.Cmd {
	return m.model.Init()
}

func (m rootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m rootScreenModel) View() string {
	return m.model.View()
}

func (m rootScreenModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	m.model = model
	return m.model, m.model.Init()
}
