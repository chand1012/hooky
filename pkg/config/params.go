package config

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
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
