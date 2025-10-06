package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runback/commands"
	"runback/utils/fs"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/servusdei2018/shards/v2"
)

// Define a struct to hold the CLI arguments
type CommandLineConfig struct {
	ConfigFile     string
	TokenFilePath  string
	TestGuildId    string
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
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("[INFO] Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Gateway.Open()
	if err != nil {
		log.Fatalf("[ERROR] Cannot open the session: %v", err)
	}
	defer s.Gateway.Close()

	log.Println("[INFO] Starting shard manager...")
	err = s.Start()
	if err != nil {
		fmt.Println("[ERROR] Error starting manager,", err)
		return
	}

	log.Println("[INFO] Adding commands...")
	for _, v := range commands.AllCommands {
		errs := s.ApplicationCommandCreate(config.TestGuildId, v)
		for _, err := range errs {
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
		}
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("[SUCCESS] Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if config.RemoveCommands {
		log.Println("[INFO] Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		registeredCommands, err := s.Gateway.ApplicationCommands(s.Gateway.State.Application.ID, config.TestGuildId)
		if err != nil {
			log.Fatalf("Could not fetch registered commands: %v", err)
		}
		for _, v := range commands.AllCommands {
			for _, rc := range registeredCommands {
				if rc.Name == v.Name {
					errs := s.ApplicationCommandDelete(config.TestGuildId, rc)
					for _, err := range errs {
						if err != nil {
							log.Panicf("Cannot delete '%v' command: %v", rc.Name, err)
						}
					}
				}
			}
		}
	}

	// Cleanly close down the Manager.
	fmt.Println("[INFO] Stopping shard manager...")
	s.Shutdown()
	fmt.Println("[SUCCESS] Shard manager stopped. Bot is shut down.")
}

func parseFlags() CommandLineConfig {
	var cfg CommandLineConfig

	// Bind command-line flags to struct fields
	flag.StringVar(&cfg.ConfigFile, "config-file", "config.json", "Path to a JSON file to load config values from (overrides other flags if present: see examples). Defaults to `config.json`.")

	flag.StringVar(&cfg.TokenFilePath, "token", "token.txt", "Path to txt file containing the token. Defaults to `token.txt`.")

	flag.StringVar(&cfg.TestGuildId, "guild-id", "", "GuildId of the test server.")

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
	log.Println("[CLI] Bot test server's guild id: " + cfg.TestGuildId)

	var boolStr string
	if cfg.RemoveCommands {
		boolStr = "true"
	} else {
		boolStr = "false"
	}

	log.Println("[CLI] Bot test server's guild id: " + boolStr)

	return cfg
}
