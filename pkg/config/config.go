package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

type Param struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Required    bool   `json:"required,omitempty"`
}

func (p *Param) ToApplicationCommandOption() (*discordgo.ApplicationCommandOption, error) {
	var t discordgo.ApplicationCommandOptionType
	switch p.Type {
	case "string":
		t = discordgo.ApplicationCommandOptionString
	case "integer":
		t = discordgo.ApplicationCommandOptionInteger
	case "boolean":
		t = discordgo.ApplicationCommandOptionBoolean
	default:
		return nil, fmt.Errorf("unknown command option type: %s", p.Type)
	}
	return &discordgo.ApplicationCommandOption{
		Type:        t,
		Name:        p.Name,
		Description: p.Description,
		Required:    p.Required,
	}, nil
}

type Params []Param

func (p Params) ToApplicationCommandOptions() ([]*discordgo.ApplicationCommandOption, error) {
	options := make([]*discordgo.ApplicationCommandOption, 0)
	for _, param := range p {
		opt, err := param.ToApplicationCommandOption()
		if err != nil {
			return nil, err
		}
		options = append(options, opt)
	}
	return options, nil
}

type Headers map[string]string

type Command struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Method      string  `json:"method"`
	URL         string  `json:"url"`
	Body        Params  `json:"body,omitempty"`
	Query       Params  `json:"query,omitempty"`
	Form        Params  `json:"form,omitempty"`
	Headers     Headers `json:"headers,omitempty"`
}

func (c *Command) ToApplicationCommand() (*discordgo.ApplicationCommand, error) {
	options, err := c.Body.ToApplicationCommandOptions()
	if err != nil {
		return nil, err
	}
	queryOptions, err := c.Query.ToApplicationCommandOptions()
	if err != nil {
		return nil, err
	}
	formOptions, err := c.Form.ToApplicationCommandOptions()
	if err != nil {
		return nil, err
	}
	// if no errors, append query and form options to the body options
	options = append(options, queryOptions...)
	options = append(options, formOptions...)
	return &discordgo.ApplicationCommand{
		Name:        c.Name,
		Description: c.Description,
		Options:     options,
	}, nil
}

type Commands []Command

func (c Commands) ToApplicationCommands() ([]*discordgo.ApplicationCommand, error) {
	commands := make([]*discordgo.ApplicationCommand, 0)
	for _, command := range c {
		cmd, err := command.ToApplicationCommand()
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}
	return commands, nil
}

func (c Commands) FindByName(name string) *Command {
	for _, command := range c {
		if command.Name == name {
			return &command
		}
	}
	return nil
}

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
