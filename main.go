package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"runback/bot"
	"runback/events"
	"runback/utils/fs"

	"github.com/bwmarrin/discordgo"
	"github.com/servusdei2018/shards/v2"
)

// Define a struct to hold the CLI arguments
type CommandLineConfig struct {
	ConfigFile    string
	TokenFilePath string
	TestGuildId   string
}

var config = parseFlags()

func main() {
	var err error

	// Create a new shard manager using the provided bot token.
	bot.ShardsMgr, err = shards.New("Bot " + fs.ReadFileWhole(config.TokenFilePath))
	if err != nil {
		fmt.Println("[ERROR] Error creating manager,", err)
		return
	}

	// Register event handlers
	bot.ShardsMgr.AddHandler(events.OnConnect)
	bot.ShardsMgr.AddHandler(events.MessageCreate)

	// Register bot commands
	// ---------------------

	// In this example, we only care about receiving message events.
	bot.ShardsMgr.RegisterIntent(discordgo.IntentsGuildMessages)

	fmt.Println("[INFO] Starting shard manager...")

	// Start all of our shards and begin listening.
	err = bot.ShardsMgr.Start()
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
	bot.ShardsMgr.Shutdown()
	fmt.Println("[SUCCESS] Shard manager stopped. Bot is shut down.")
}

func parseFlags() CommandLineConfig {
	var cfg CommandLineConfig

	// Bind command-line flags to struct fields
	flag.StringVar(&cfg.ConfigFile, "config-file", "config.json", "Path to a JSON file to load config values from (overrides other flags if present: see examples). Defaults to `config.json`.")

	flag.StringVar(&cfg.TokenFilePath, "token", "token.txt", "Path to txt file containing the token. Defaults to `token.txt`.")

	flag.StringVar(&cfg.TestGuildId, "guild-id", "", "GuildId of the test server.")

	// Parse the flags
	log.Println("[CLI] Parsing arguments.")
	flag.Parse()
	if cfg.ConfigFile != "" {
		log.Println("[CLI] Loading config from JSON file:", cfg.ConfigFile)
		data, err := os.ReadFile(cfg.ConfigFile)
		if err != nil {
			log.Fatalf("[CLI] Failed to read config file: %v", err)
		}
		var jsonCfg CommandLineConfig
		if err := json.Unmarshal(data, &jsonCfg); err != nil {
			log.Fatalf("[CLI] Failed to parse config file: %v", err)
		}
		cfg = jsonCfg
	}

	log.Println("[CLI] Bot config file: " + cfg.ConfigFile)
	log.Println("[CLI] Bot token file: " + cfg.TokenFilePath)
	log.Println("[CLI] Bot test server's guild id: " + cfg.TestGuildId)

	return cfg
}
