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

type Params = []Param

type Headers = map[string]string

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
	options := make([]*discordgo.ApplicationCommandOption, 0)
	for _, p := range c.Body {
		opt, err := p.ToApplicationCommandOption()
		if err != nil {
			return nil, err
		}
		options = append(options, opt)
	}
	for _, p := range c.Query {
		opt, err := p.ToApplicationCommandOption()
		if err != nil {
			return nil, err
		}
		options = append(options, opt)
	}
	for _, p := range c.Form {
		opt, err := p.ToApplicationCommandOption()
		if err != nil {
			return nil, err
		}
		options = append(options, opt)
	}
	return &discordgo.ApplicationCommand{
		Name:        c.Name,
		Description: c.Description,
		Options:     options,
	}, nil
}

// LoadConfigs loads all .json files in the given directory and returns them as a list of Command structs.
func LoadConfigs(dir string) ([]Command, error) {
	var commands []Command
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
		commands = append(commands, command)
	}
	return commands, nil
}
