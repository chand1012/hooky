package config

import "github.com/bwmarrin/discordgo"

type Headers map[string]string

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
