package commands

import "github.com/bwmarrin/discordgo"

var AllCommands = []*discordgo.ApplicationCommand{
	&pingCommand,
}

var AllCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	pingCommand.Name: pingCommandHandler,
}
