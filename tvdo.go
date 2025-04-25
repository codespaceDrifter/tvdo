/*
TO DO:
read and understand every line

clean it up. make it functional? 
3 types of functions
no class. yes return : purely functional
yes class. no return: modify state in place
no class. no return: side effects

fix the save load
add the auto decrease day system
for the ROOT Day. do DAY until 2030. 
*/


package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Task struct {
	Name     string   `json:"name"`
	Subtasks []*Task  `json:"subtasks"`
	Due      int      `json:"due"`
	Parent   *Task    `json:"-"`
}

var root = &Task{Name: "ROOT"}

type model struct {
	tasks      []*Task
	cursor     int
	input      string
	adding     bool
	editingDue bool
}

func getSavePath() string {
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "tasks.json")
}

func saveToFile() {
	f, err := os.Create(getSavePath())
	if err != nil {
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(root)
}


func (t *Task) setParents(parent *Task) {
    t.Parent = parent
    for _, sub := range t.Subtasks {
        sub.setParents(t)
    }
}

func (t *Task) loadFromFile() {
    f, err := os.Open(getSavePath())
    if err != nil {
        return
    }
    defer f.Close()
    var loaded Task
    if err := json.NewDecoder(f).Decode(&loaded); err == nil {
        *t = loaded
        t.setParents(nil)
    }
}

func (m *model) updateTasks() {
	m.tasks = flatten(root)
}

func flatten(task *Task) []*Task {
	result := []*Task{task}

	sorted := make([]*Task, len(task.Subtasks))
	copy(sorted, task.Subtasks)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Due < sorted[j].Due
	})

	for _, sub := range sorted {
		result = append(result, flatten(sub)...)
	}

	return result
}



func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if m.adding || m.editingDue{
			switch key {
			case "enter":
				if strings.TrimSpace(m.input) != "" {
					selected := m.tasks[m.cursor]

					if (m.adding){
						t := &Task{Name: m.input, Parent: selected, Due: 0}
						selected.Subtasks = append(selected.Subtasks, t)
						m.adding = false
					} else if (m.editingDue) {
						if due, err := strconv.Atoi(strings.TrimSpace(m.input)); err == nil {
							m.tasks[m.cursor].Due = due
						}
						m.editingDue = false
					}

					m.input = ""
					m.updateTasks()
					saveToFile()
				}
			case "backspace":
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			default:
				m.input += key
			}
		}else {
			switch key {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "j":
				if m.cursor < len(m.tasks)-1 {
					m.cursor++
				}
			case "o":
				m.adding = true
				m.input = ""
			case "a":
				m.editingDue = true
				m.input = ""
			case "d":
				if m.tasks[m.cursor] != root {
					parent := m.tasks[m.cursor].Parent
					if parent != nil {
						for i, sub := range parent.Subtasks {
							if sub == m.tasks[m.cursor] {
								parent.Subtasks = append(parent.Subtasks[:i], parent.Subtasks[i+1:]...)
								break
							}
						}
						m.cursor--
						if m.cursor < 0 {
							m.cursor = 0
						}
						m.updateTasks()
						saveToFile()
					}
				}
			}
		}
	}
	return m, nil
}

func depth(t *Task) int {
	d := 0
	for p := t.Parent; p != nil; p = p.Parent {
		d++
	}
	return d
}




func (m model) View() string {
	var b strings.Builder
	maxNameLen := 0
	for _, t := range m.tasks {
		depth := depth(t)
		nameLen := len(strings.Repeat("  ", depth) + t.Name)
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}
	}

	for i, t := range m.tasks {
		indent := strings.Repeat("  ", depth(t))
		taskText := fmt.Sprintf("%s%s", indent, t.Name)
		padLen := maxNameLen - len(taskText) + 2
		padding := strings.Repeat(" ", padLen)
		line := fmt.Sprintf("%s%s%d", taskText, padding, t.Due)
		if i == m.cursor {
			// bold, cyan
			line = fmt.Sprintf("\033[1;36m> %s\033[0m", line)
		} else {
			line = "  " + line
		}
		b.WriteString(line + "\n")
	}
	if m.adding {
		b.WriteString("\nEnter task name: " + m.input)
	} else if m.editingDue {
		b.WriteString("\nEnter due value: " + m.input)
	}
	return b.String()
}


func main() {
	root.loadFromFile()
	initialTasks := flatten(root)
	var bubbleModel = model{
		tasks: initialTasks,
		cursor: 0,
	}
	p := tea.NewProgram(bubbleModel)
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}


