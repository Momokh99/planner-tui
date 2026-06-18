package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
)

const dataFile = ".todo.json"

type Todo struct {
	ID        int
	Title     string
	Completed bool
	DueDate   time.Time
}

type model struct {
	todos     []Todo
	cursor    int
	dayCursor int
	weekStart time.Time
	view      viewState
}
type viewState int

const (
	weekView viewState = iota
	listView
)

// try to find whare to safe the todos and save them
func dataPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return dataFile
	}
	return filepath.Join(home, dataFile)
}
func loadTodos() ([]Todo, error) {
	data, err := os.ReadFile(dataPath())
	if err != nil || len(data) == 0 {
		return []Todo{}, nil
	}
	var todos []Todo
	err = json.Unmarshal(data, &todos)
	return todos, err
}
func saveTodos(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataPath(), data, 0644)
}

// the essentials for the tui
func (m model) Init() tea.Cmd { return nil }

func initialModel() model {
	// Map Go weekday (Sun=0..Sat=6) to column index (Mon=0..Sun=6)
	dayCursor := int(time.Now().Weekday())
	if dayCursor == 0 {
		dayCursor = 6
	} else {
		dayCursor--
	}
	return model{
		todos:     []Todo{},
		weekStart: mondayOfWeek(time.Now()),
		view:      weekView,
		cursor:    0,
		dayCursor: dayCursor,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Global keys
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "t":
			if m.view == weekView {
				m.view = listView
			} else {
				m.view = weekView
			}
		case "left":
			m.dayCursor--
			if m.dayCursor < 0 {
				m.dayCursor = 6
				m.weekStart = m.weekStart.AddDate(0, 0, -7) // slide to previous week
			}
		case "right":
			m.dayCursor++
			if m.dayCursor > 6 {
				m.dayCursor = 0
				m.weekStart = m.weekStart.AddDate(0, 0, 7) // slide to next week
			}
		case "up", "k":
			if m.view == listView && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.view == listView && m.cursor < len(m.todos)-1 {
				m.cursor++
			}
		case "enter":
			if m.view == listView && len(m.todos) > 0 {
				m.todos[m.cursor].Completed = !m.todos[m.cursor].Completed
				saveTodos(m.todos)
			}
		}
	}
	return m, nil
}
func (m model) View() tea.View {
	switch m.view {
	case weekView:
		return tea.NewView(m.weekView())
	case listView:
		return tea.NewView(m.listView()) // we'll build this in Step 5
	default:
		return tea.NewView("unknown view")
	}
}

func main() {
	todos, _ := loadTodos() // load from disk
	m := initialModel()     // build initial state
	m.todos = todos         // inject loaded data
	p := tea.NewProgram(m)  // start the TUI
	//           ^
	// The program calls m.Init() automatically
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
