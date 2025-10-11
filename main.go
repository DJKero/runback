package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"runback/bot"
	"runback/commands"
	"runback/db"
	"runback/utils/fs"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/servusdei2018/shards/v2"
)

// Define a struct to hold the CLI arguments
type CommandLineConfig struct {
	ConfigFile     string
	TokenFilePath  string
	GuildId        string
	RemoveCommands bool
}

var config = parseFlags()

var s *shards.Manager

func init() { flag.Parse() }

func init() {
	var err error
	s, err = shards.New("Bot " + fs.ReadFileWhole(config.TokenFilePath))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commands.AllCommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	db.New()
	defer db.DBPool.Close()

	var err error
	var errs []error

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("[INFO] Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	log.Println("[INFO] Starting gateway session...")
	err = s.Gateway.Open()
	if err != nil {
		log.Fatalln("[ERROR] Cannot open gateway session:", err)
	}

	log.Println("[INFO] Starting shard manager...")
	err = s.Start()
	if err != nil {
		log.Fatalln("[ERROR] Error starting manager,", err)
		return
	}

	log.Println("[INFO] Adding commands...")
	for _, v := range commands.AllCommands {
		errs = s.ApplicationCommandCreate(config.GuildId, v)
		for _, err := range errs {
			if err != nil {
				log.Printf("[ERROR] Cannot create '%v' command: %v", v.Name, err)
			}
		}
	}

	bot.Client = bot.Bot{
		ShardsMgr: s,
	}

	log.Println("[SUCCESS] Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if config.RemoveCommands {
		log.Println("[INFO] Removing commands...")
		registeredCommands, err := s.Gateway.ApplicationCommands(s.Gateway.State.Application.ID, config.GuildId)
		if err != nil {
			log.Fatalln("Could not fetch registered commands:", err)
		}
		for _, v := range commands.AllCommands {
			for _, rc := range registeredCommands {
				if rc.Name == v.Name {
					errs := s.ApplicationCommandDelete(config.GuildId, rc)
					for _, err := range errs {
						if err != nil {
							log.Fatalf("Cannot delete '%v' command: %v", rc.Name, err)
						}
					}
				}
			}
		}
	}

	log.Println("[INFO] Stopping gateway session...")
	err = s.Gateway.Close()
	if err != nil {
		log.Fatalln("[ERROR] Failed to stop gateway session:", err)
	}

	log.Println("[INFO] Stopping shard manager...")
	err = s.Shutdown()
	if err != nil {
		log.Println("[ERROR] Failed to stop shard manager:", err)
		log.Fatalln("[ERROR] Bot is not shut down properly.")
	} else {
		log.Println("[SUCCESS] Shard manager stopped.")
		log.Println("[SUCCESS] Bot is shut down properly.")
	}
}

func parseFlags() CommandLineConfig {
	var cfg CommandLineConfig

	// Bind command-line flags to struct fields
	flag.StringVar(&cfg.ConfigFile, "config-file", "config.json", "Path to a JSON file to load config values from (overrides other flags if present: see examples). Defaults to `config.json`.")

	flag.StringVar(&cfg.TokenFilePath, "token", "token.txt", "Path to txt file containing the token. Defaults to `token.txt`.")

	flag.StringVar(&cfg.GuildId, "guild-id", "", "GuildId of the test server.")

	flag.BoolVar(&cfg.RemoveCommands, "remove-commands", true, "Whether to remove all commands added by the bot on shutdown. Defaults to true.")

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
	log.Println("[CLI] Bot test server's guild id: " + cfg.GuildId)
	log.Println("[CLI] Bot test server's guild id: " + boolToString(cfg.RemoveCommands))

	return cfg
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
