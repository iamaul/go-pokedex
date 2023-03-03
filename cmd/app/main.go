package main

import (
	"log"
	"os"

	"github.com/iamaul/go-pokedex/config"
	"github.com/iamaul/go-pokedex/internal/server"
	"github.com/iamaul/go-pokedex/pkg/db/mongodb"
	"github.com/iamaul/go-pokedex/pkg/logger"
	"github.com/iamaul/go-pokedex/pkg/utils"
)

func main() {
	log.Println("Starting API server")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	mongoClient, err := mongodb.NewClient(cfg)
	if err != nil {
		appLogger.Fatalf("MongoDB init: %s", err)
	} else {
		appLogger.Infof("MongoDB connected.")
	}
	db := mongoClient.Database("pokedex")

	s := server.NewServer(cfg, db, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
