package handle

import "github.com/bwmarrin/discordgo"

type OptionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

func ParseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) (om OptionMap) {
	om = make(OptionMap)
	for _, opt := range options {
		om[opt.Name] = opt
	}
	return
}
