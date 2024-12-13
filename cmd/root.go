/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/aisalamdag23/promptme-cli/internal/infrastructure/caching"
	"github.com/aisalamdag23/promptme-cli/internal/infrastructure/config"
	"github.com/aisalamdag23/promptme-cli/internal/infrastructure/logger"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "promptme-cli",
		Short: "A command-line tool that leverages a language model API to generate text based on user prompts",
		Long: `A command-line tool that leverages a language model API to generate text based on user prompts. 
	
Requirements:

- Accept user prompts
- Stream real-time responses (output as the LM generates text)
- Process and display the incoming data quickly, and provide measurement of response time
- Handle errors gracefully, and provide informative error messages
- Log errors, warnings, informative messages for debugging and monitoring
- Include structural and behavioral design patterns
- Make integrations to LM APIs swappable`,
	}
	log        *logrus.Entry
	cfg        *config.Config
	inMemCache caching.Cache
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Fatal("config.load.failed:", err)
	}

	log = logger.NewLogger(cfg.General.LogLevel)

	inMemCache = caching.NewInMemory(cfg)
}
