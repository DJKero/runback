package events

import (
	"fmt"
	"runback/bot"

	"github.com/bwmarrin/discordgo"
)

// This callback will be called every time a
// new  message is created on any channel that the authenticated bot has
// access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	// This isn't required in this specific example but it's a good
	// practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "ping":
		// If the message is "ping" reply with "Pong!"
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "pong":
		// If the message is "pong" reply with "Ping!"
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	case "restart":
		// If the message is "restart" restart the shard manager and rescale
		// if necessary, all with zero down-time.
		var err error
		s.ChannelMessageSend(m.ChannelID, "[INFO] Restarting shard manager...")
		fmt.Println("[INFO] Restarting shard manager...")
		bot.ShardsMgr, err = bot.ShardsMgr.Restart()
		if err != nil {
			fmt.Println("[ERROR] Error restarting manager,", err)
		} else {
			s.ChannelMessageSend(m.ChannelID, "[SUCCESS] Manager successfully restarted.")
			fmt.Println("[SUCCESS] Manager successfully restarted.")
		}
	}
}
