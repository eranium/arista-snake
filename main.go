package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/aristanetworks/goeapi"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"time"
)

var difficultySpeedMap = map[string]time.Duration{
	"easy":   time.Second,
	"medium": 700 * time.Millisecond,
	"hard":   400 * time.Millisecond,
}

func loadGridDefinitionFromFile(path string) ([][]string, error) {
	var grid [][]string

	v, err := os.ReadFile(path)
	if err != nil {
		return grid, err
	}

	if err := json.Unmarshal(v, &grid); err != nil {
		return grid, err
	}

	return grid, nil
}

func run(ctx context.Context) error {
	var configFile string
	flag.StringVar(&configFile, "config-file", "", "eAPI configuration file")

	var connection string
	flag.StringVar(&connection, "connection", "", "eAPI connection to use")

	var difficulty string
	flag.StringVar(&difficulty, "difficulty", "medium", "Game difficulty, available options are 'easy', 'medium' and 'hard'")

	var gridDefinitionFile string
	flag.StringVar(&gridDefinitionFile, "grid-definition-file", "", "JSON file containing the port grid definitions for the switch")

	var dryRun bool
	flag.BoolVar(&dryRun, "dry-run", false, "Enable dry-run mode, in which no switch communication is required. Defaults to false.")

	flag.Parse()

	speed, ok := difficultySpeedMap[difficulty]
	if !ok {
		return fmt.Errorf("invalid difficulty %s, available options are 'easy', 'medium' and 'hard'", difficulty)
	}

	grid, err := loadGridDefinitionFromFile(gridDefinitionFile)
	if err != nil {
		return err
	}

	var controller InterfaceLEDController
	if dryRun {
		controller = NewNoOpInterfaceLEDController()
	} else {
		if configFile == "" {
			return errors.New("specify a configuration file using --config-file")
		}

		if connection == "" {
			return errors.New("specify a connection using --connection")
		}

		goeapi.LoadConfig(configFile)
		node, err := goeapi.ConnectTo(connection)
		if err != nil {
			return err
		}

		var managedIntfs []string
		for _, row := range grid {
			for _, intf := range row {
				managedIntfs = append(managedIntfs, intf)
			}
		}

		controller = NewAristaInterfaceLEDController(managedIntfs, node)
	}

	start := coord{x: 0, y: 0}
	m := model{
		controller: controller,
		grid:       grid,
		direction:  DirectionRight,
		body:       []coord{start},
		speed:      speed,
	}

	if err := controller.SetAllStatuses(InterfaceLEDStatusOff); err != nil {
		return err
	}
	m.foodLocation = m.placeFood()

	p := tea.NewProgram(m, tea.WithContext(ctx))

	if _, err := p.Run(); err != nil {
		log.Fatal("error: ", err)
	}

	if err := controller.SetAllStatuses(InterfaceLEDStatusOrange); err != nil {
		return err
	}
	time.Sleep(time.Second)
	if err := controller.SetAllStatuses(InterfaceLEDStatusOff); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
