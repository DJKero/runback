package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runback/handlers"
	"runback/utils/fs"
	"syscall"

	"github.com/servusdei2018/shards/v2"
)

// Define a struct to hold the CLI arguments
type CommandLineConfig struct {
	ConfigFile       string
	BotTokenFilePath string

	StartProfiler bool
}

var (
	config = parseFlags()
)

func main() {
	go startBot()

	if config.StartProfiler { // Launch the profiler if enabled
		go startProfiler()
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Println("[MAIN] Waiting for exit signal.")
	<-sc
}

func startBot() {
	var mgr *shards.Manager
	var err error

	// Create a new shard manager using the provided bot token.
	mgr, err = shards.New("Bot " + fs.ReadFileWhole(config.BotTokenFilePath))
	if err != nil {
		fmt.Println("[ERROR] Error creating manager,", err)
		return
	}

	handlers.StartShards(mgr)
}

func startProfiler() {
	log.Println("[PROFILER] Starting pprof server at :6060")
	log.Println("[PROFILER]", http.ListenAndServe("0.0.0.0:6060", nil))
}

func parseFlags() CommandLineConfig {
	var cfg CommandLineConfig

	// Bind command-line flags to struct fields
	flag.StringVar(&cfg.ConfigFile, "config-file", "", "Path to a JSON file to load config values from (overrides other flags if present: see examples)")

	flag.StringVar(&cfg.BotTokenFilePath, "token", "token.txt", "Path to txt file containing the token. Defaults to `token.txt`.")

	flag.BoolVar(&cfg.StartProfiler, "profiler", false, "Flag that enables the profiler")

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
	log.Println("[CLI] Bot Token File: " + cfg.BotTokenFilePath)

	return cfg
}
