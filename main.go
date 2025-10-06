package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"runback/commands"
	"runback/utils/fs"

	"github.com/bwmarrin/discordgo"
)

// Define a struct to hold the CLI arguments
type CommandLineConfig struct {
	ConfigFile     string
	TokenFilePath  string
	TestGuildId    string
	RemoveCommands bool
}

var config = parseFlags()

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + fs.ReadFileWhole(config.TokenFilePath))
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
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands.AllCommands))
	for i, v := range commands.AllCommands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, config.TestGuildId, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if config.RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, config.TestGuildId, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
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
