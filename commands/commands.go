package commands

import "github.com/bwmarrin/discordgo"

var AllCommands = []*discordgo.ApplicationCommand{
	&PingCommand,
}

var AllCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	PingCommand.Name: pingCommand,
}
