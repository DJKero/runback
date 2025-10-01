package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// This callback will be called every time one
// of our shards connects.
func OnConnect(s *discordgo.Session, evt *discordgo.Connect) {
	fmt.Printf("[INFO] Shard #%v connected.\n", s.ShardID)
}
