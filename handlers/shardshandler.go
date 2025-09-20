package handlers

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/servusdei2018/shards/v2"
)

var mgr *shards.Manager

func StartShards(m *shards.Manager) {
	mgr = m
	var err error

	// Register the messageCreate func as a callback for MessageCreate
	// events.
	mgr.AddHandler(messageCreate)
	// Register the onConnect func as a callback for Connect events.
	mgr.AddHandler(onConnect)

	// In this example, we only care about receiving message events.
	mgr.RegisterIntent(discordgo.IntentsGuildMessages)

	fmt.Println("[INFO] Starting shard manager...")

	// Start all of our shards and begin listening.
	err = mgr.Start()
	if err != nil {
		fmt.Println("[ERROR] Error starting manager,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("[SUCCESS] Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Manager.
	fmt.Println("[INFO] Stopping shard manager...")
	mgr.Shutdown()
	fmt.Println("[SUCCESS] Shard manager stopped. Bot is shut down.")
}
