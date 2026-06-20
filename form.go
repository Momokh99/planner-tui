package main

func (m *model) initAddForm() {
	m.titleInput.Reset()
	m.titleInput.Focus()
	cursorDate := m.weekStart.AddDate(0, 0, m.dayCursor)
	m.dateInput.SetValue(cursorDate.Format("2006-01-02"))
	m.dateInput.Blur()
	m.formStep = 0
	m.editingID = -1
	m.view = formView
}

func (m *model) initEditForm(id int) {
	todo := m.todos[id]
	m.titleInput.SetValue(todo.Title)
	m.titleInput.Focus()
	m.dateInput.SetValue(todo.DueDate.Format("2006-01-02"))
	m.dateInput.Blur()
	m.formStep = 0
	m.editingID = id
	m.view = formView
}

func (m model) formView() string {
	s := "Add Todo\n\n"
	if m.editingID >= 0 {
		s = "Edit Todo\n\n"
	}
	if m.formStep == 0 {
		s += "Title:\n" + m.titleInput.View()
	} else {
		s += "Date (YYYY-MM-DD):\n" + m.dateInput.View()
	}
	s += "\n\n  [Enter] next  [Esc] cancel"
	return s
}
