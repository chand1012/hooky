package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadConfigs loads all .json files in the given directory and returns them as a list of Command structs.
func LoadConfigs(dir string) (Commands, error) {
	var commands Commands
	names := make(map[string]Command)
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		var command Command
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &command)
		if err != nil {
			return nil, err
		}
		if _, ok := names[command.Name]; ok {
			return nil, fmt.Errorf("duplicate command name: %s", command.Name)
		}
		commands = append(commands, command)
		names[command.Name] = command
	}
	return commands, nil
}
