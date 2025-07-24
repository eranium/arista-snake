package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
	"strconv"
	"time"
)

type Direction int

const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

type coord struct {
	x, y int
}

type model struct {
	controller   InterfaceLEDController
	grid         [][]string
	direction    Direction
	body         []coord
	growNext     bool
	foodLocation coord
	speed        time.Duration
}

func (m model) placeFood() coord {
	for {
		x := rand.Intn(len(m.grid[0]))
		y := rand.Intn(len(m.grid))
		c := coord{x, y}

		occupied := false
		for _, b := range m.body {
			if b == c {
				occupied = true
				break
			}
		}

		if !occupied {
			if err := m.controller.SetStatus(m.grid[y][x], InterfaceLEDStatusGreen); err != nil {
				panic(err)
			}
			return c
		}
	}
}

func (m model) Init() tea.Cmd {
	return m.tick()
}

func (m model) tick() tea.Cmd {
	return tea.Tick(m.speed, func(_ time.Time) tea.Msg {
		return moveMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case moveMsg:
		return m.move()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.direction = DirectionUp
		case "down", "j":
			m.direction = DirectionDown
		case "left", "h":
			m.direction = DirectionLeft
		case "right", "l":
			m.direction = DirectionRight
		}
	}

	return m, nil
}

func (m model) move() (tea.Model, tea.Cmd) {
	head := m.body[0]
	var newHead coord

	switch m.direction {
	case DirectionUp:
		if head.y <= 0 {
			return m, tea.Quit
		}
		newHead = coord{head.x, head.y - 1}
	case DirectionDown:
		if head.y >= len(m.grid)-1 {
			return m, tea.Quit
		}

		newHead = coord{head.x, head.y + 1}
	case DirectionLeft:
		if head.x <= 0 {
			return m, tea.Quit
		}

		newHead = coord{head.x - 1, head.y}
	case DirectionRight:
		if head.x >= len(m.grid[0])-1 {
			return m, tea.Quit
		}

		newHead = coord{head.x + 1, head.y}
	}

	ateFood := newHead == m.foodLocation
	if ateFood {
		m.growNext = true
	}

	if err := m.controller.SetStatus(m.grid[newHead.y][newHead.x], InterfaceLEDStatusOrange); err != nil {
		panic(err)
	}

	if !m.growNext {
		tail := m.body[len(m.body)-1]
		if err := m.controller.SetStatus(m.grid[tail.y][tail.x], InterfaceLEDStatusOff); err != nil {
			panic(err)
		}
		m.body = m.body[:len(m.body)-1]
	} else {
		m.growNext = false
	}

	for _, c := range m.body {
		if c == newHead {
			return m, tea.Quit
		}
	}

	m.body = append([]coord{newHead}, m.body...)

	if ateFood {
		m.foodLocation = m.placeFood()
	}

	return m, m.tick()
}

func (m model) View() string {
	bodyMap := make(map[coord]bool)
	for _, c := range m.body {
		bodyMap[c] = true
	}

	var s string
	for y, row := range m.grid {
		for x, name := range row {
			c := coord{x, y}

			var char string
			switch {
			case m.body[0] == c:
				switch m.direction {
				case DirectionUp:
					char = "^"
				case DirectionDown:
					char = "v"
				case DirectionLeft:
					char = "<"
				case DirectionRight:
					char = ">"
				}
			case bodyMap[c]:
				char = "x"
			case c == m.foodLocation:
				char = "o"
			default:
				char = "-"
			}

			s += char + ": " + name + "\t"
		}
		s += "\n"
	}

	s += "\nScore: " + strconv.Itoa(len(m.body)-1)

	return s
}

type moveMsg struct{}
