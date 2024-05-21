package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/itchyny/gojq"
)

type Param struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Required    bool     `json:"required,omitempty"`
	Options     []string `json:"options,omitempty"`
	Default     any      `json:"default,omitempty"`
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

	var choices []*discordgo.ApplicationCommandOptionChoice
	if len(p.Options) > 0 {
		choices = make([]*discordgo.ApplicationCommandOptionChoice, 0)
		for _, opt := range p.Options {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  opt,
				Value: opt,
			})
		}
	}

	return &discordgo.ApplicationCommandOption{
		Type:        t,
		Name:        p.Name,
		Description: p.Description,
		Required:    p.Required,
		Choices:     choices,
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

type ParseJSONResponse map[string]string

func (p *ParseJSONResponse) Parse(body []byte) (map[string]string, error) {
	respData := make(map[string]string)
	var input map[string]interface{}
	if err := json.Unmarshal(body, &input); err != nil {
		return nil, err
	}
	for key, value := range *p {
		query, err := gojq.Parse(value)
		if err != nil {
			log.Errorf("could not parse jq query: %s", err)
			continue
		}
		iter := query.Run(input)
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
					break
				}
				log.Errorf("could not run jq query: %s", err)
			}
			respData[key] = fmt.Sprintf("%v", v)
		}
	}

	return respData, nil
}

type Command struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Method           string            `json:"method"`
	URL              string            `json:"url"`
	Body             Params            `json:"body,omitempty"`
	BodyTemplate     string            `json:"body_template,omitempty"`
	Query            Params            `json:"query,omitempty"`
	Form             Params            `json:"form,omitempty"`
	Headers          Headers           `json:"headers,omitempty"`
	ParseJSON        ParseJSONResponse `json:"parse_json,omitempty"`
	ResponseTemplate string            `json:"response_template,omitempty"`
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
