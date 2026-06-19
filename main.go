package main

import (
	"encoding/json"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"path/filepath"
	"time"
)

const dataFile = ".todo.json"

type Todo struct {
	ID        int
	Title     string
	Completed bool
	DueDate   time.Time
}

type model struct {
	todos      []Todo
	cursor     int
	dayCursor  int
	weekStart  time.Time
	view       viewState
	titleInput textinput.Model
	dateInput  textinput.Model
	formStep   int
	editingID  int
	width      int
	height     int
}
type viewState int

const (
	weekView viewState = iota
	listView
	formView
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
func (m model) Init() tea.Cmd { return textinput.Blink }

func initialModel() model {
	// Map Go weekday (Sun=0..Sat=6) to column index (Mon=0..Sun=6)
	dayCursor := int(time.Now().Weekday())
	if dayCursor == 0 {
		dayCursor = 6
	} else {
		dayCursor--
	}

	ti := textinput.New()
	ti.Placeholder = "Task title"
	ti.CharLimit = 100
	ti.Width = 40

	di := textinput.New()
	di.Placeholder = "2026-06-19"
	di.CharLimit = 10
	di.Width = 12

	return model{
		todos:      []Todo{},
		weekStart:  mondayOfWeek(time.Now()),
		view:       weekView,
		cursor:     0,
		dayCursor:  dayCursor,
		titleInput: ti,
		dateInput:  di,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.titleInput.Width = msg.Width / 3
	case tea.KeyMsg:
		// If in form view, handle form keys first
		if m.view == formView {
			switch msg.String() {
			case "esc":
				m.view = listView
			case "enter":
				if m.formStep == 0 {
					if m.titleInput.Value() == "" {
						return m, nil
					}
					m.formStep = 1
					m.dateInput.Focus()
					m.titleInput.Blur()
				} else {
					dueDate, err := time.Parse("2006-01-02", m.dateInput.Value())
					if err != nil {
						return m, nil
					}
					todo := Todo{Title: m.titleInput.Value(), DueDate: dueDate}
					if m.editingID >= 0 {
						todo.ID = m.todos[m.editingID].ID
						todo.Completed = m.todos[m.editingID].Completed
						m.todos[m.editingID] = todo
					} else {
						maxID := 0
						for _, t := range m.todos {
							if t.ID > maxID {
								maxID = t.ID
							}
						}
						todo.ID = maxID + 1
						m.todos = append(m.todos, todo)
					}
					saveTodos(m.todos)
					m.view = listView
				}
			default:
				var cmd tea.Cmd
				if m.formStep == 0 {
					m.titleInput, cmd = m.titleInput.Update(msg)
				} else {
					m.dateInput, cmd = m.dateInput.Update(msg)
				}
				return m, cmd
			}
			return m, nil
		}

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
		case "a":
			if m.view == listView {
				m.initAddForm()
			}
		case "e":
			if m.view == listView && len(m.todos) > 0 {
				m.initEditForm(m.cursor)
			}
		case "d":
			if m.view == listView && len(m.todos) > 0 {
				m.todos = append(m.todos[:m.cursor], m.todos[m.cursor+1:]...)
				if m.cursor >= len(m.todos) && m.cursor > 0 {
					m.cursor--
				}
				saveTodos(m.todos)
			}

		}
	}
	return m, nil
}
func (m model) View() string {
	switch m.view {
	case weekView:
		return m.weekView()
	case listView:
		return m.listView()
	case formView:
		return m.formView()
	default:
		return "unknown view"
	}
}

func main() {
	todos, _ := loadTodos()
	m := initialModel()
	m.todos = todos
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
