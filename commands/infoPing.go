package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var PingCommand = discordgo.ApplicationCommand{
	Name:        "infoPing",
	Description: "Shows the bots ping.",
}

func pingCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var latency time.Duration = s.HeartbeatLatency().Round(time.Millisecond)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Ping: %s", latency.String()),
		},
	})
}
