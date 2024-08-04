package tui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/ls13g12/hockey-app/src/pkg/db"
)

const (
	firstName = iota
	lastName
	nickname
	homeShirtNumber
	awayShirtNumber
	dateOfBirthString
)

var (
	playerBox      = lipgloss.NewStyle().Width(80).Padding(0, 2)
)

func generatePlayerCardString(player db.Player) string {
	displayedPlayer := fmt.Sprintf(
		`%s %s
%s %s

%s %s
`,
		headingStyle.Width(15).Render("First Name"),
		headingStyle.Width(15).Render("Last Name"),
		textStyle.Width(15).Render(player.FirstName),
		textStyle.Width(15).Render(player.LastName),
		headingStyle.Width(15).Render("Home Number"),
		headingStyle.Width(15).Render("Away Number"),
	)
	if player.HomeShirtNumber > 0 {
		displayedPlayer += textStyle.Width(15).Render(strconv.Itoa(player.HomeShirtNumber))
	} else {
		displayedPlayer += textStyle.Width(15).Render(" ")
	}

	if player.AwayShirtNumber > 0 {
		displayedPlayer += textStyle.Width(15).Render(strconv.Itoa(player.AwayShirtNumber))
	} else {
		displayedPlayer += textStyle.Width(15).Render(" ")
	}
	return displayedPlayer
}

func generatePlayerInputs() ([]textinput.Model, error) {
	var inputs []textinput.Model = make([]textinput.Model, 6)
	inputs[firstName] = textinput.New()
	inputs[firstName].Focus()
	inputs[firstName].Prompt = ""
	inputs[firstName].CharLimit = 19
	inputs[firstName].Width = 20
	inputs[firstName].Validate = nameValidator

	inputs[lastName] = textinput.New()
	inputs[lastName].Prompt = ""
	inputs[lastName].CharLimit = 30
	inputs[lastName].Width = 30
	inputs[lastName].Validate = nameValidator

	inputs[nickname] = textinput.New()
	inputs[nickname].Prompt = ""
	inputs[nickname].CharLimit = 30
	inputs[nickname].Width = 30
	inputs[nickname].Validate = nameValidator

	inputs[homeShirtNumber] = textinput.New()
	inputs[homeShirtNumber].Prompt = ""
	inputs[homeShirtNumber].CharLimit = 3
	inputs[homeShirtNumber].Width = 20
	inputs[homeShirtNumber].Validate = numberValidator

	inputs[awayShirtNumber] = textinput.New()
	inputs[awayShirtNumber].Prompt = ""
	inputs[awayShirtNumber].CharLimit = 3
	inputs[awayShirtNumber].Width = 30
	inputs[awayShirtNumber].Validate = numberValidator

	inputs[dateOfBirthString] = textinput.New()
	inputs[dateOfBirthString].Prompt = ""
	inputs[dateOfBirthString].CharLimit = 10
	inputs[dateOfBirthString].Width = 30
	inputs[dateOfBirthString].Validate = dateValidator

	return inputs, nil
}

func nameValidator(s string) error {
	s = strings.ReplaceAll(s, " ", "")

	isAlpha := regexp.MustCompile(`[A-Za-z]+$`).MatchString(s)
	if !isAlpha {
		return fmt.Errorf("names may only contain letters")
	}

	if len(s) < 1 {
		return fmt.Errorf("names cannot be empty")
	}

	return nil
}

func numberValidator(s string) error {
	s = strings.ReplaceAll(s, " ", "")

	num, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	if num < 1 {
		return fmt.Errorf("number be greater than 0")
	}

	if num > 999 {
		return fmt.Errorf("number must be less than 1000")
	}

	return nil
}

func dateValidator(s string) error {
	s = strings.ReplaceAll(s, " ", "")

	parts := strings.Split(s, "/")
	
	switch len(parts) {
	case 1:
		dd, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid day")
		}
		if dd < 0 || dd > 31 {
			return fmt.Errorf("invalid day")
		}
	case 2:
		if len(parts[1]) == 0 {
			return nil
		}
		mm, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid month")
		}
		if mm < 0 || mm > 12 {
			return fmt.Errorf("invalid month")
		}
	case 3:
		if len(parts[2]) == 0 {
			return nil
		}
		yyyy, err := strconv.Atoi(parts[2])
		if err != nil {
			return fmt.Errorf("invalid year")
		}
		if yyyy < 0 || yyyy > time.Now().Year() - 10 {
			return fmt.Errorf("invalid year")
		}
	default:
		return fmt.Errorf("invalid format")
	}
	return nil
}
