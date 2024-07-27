package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var configFile string

func main() {
	app := CreateApplication()

	rootCmd := &cobra.Command{}

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file path")
	err := rootCmd.MarkPersistentFlagRequired("config")
	if err != nil {
		panic(err)
	}

	if err = app.RunCommands(rootCmd); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}
