package main

import (
	"fmt"
	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
	"os"
)

func init() {
	Ui = &cli.PrefixedUi{
		AskPrefix:    OutputPrefix,
		OutputPrefix: OutputPrefix,
		InfoPrefix:   OutputPrefix,
		ErrorPrefix:  ErrorPrefix,
		Ui:           &cli.BasicUi{Writer: os.Stdout},
	}

	_, err := os.Stat("tinker.toml")
	if err == nil {
		viper.SetConfigName("tinker")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
}

func main() {
	// Initialize executable commands.
	initCommands()

	args := os.Args[1:]
	cliRunner := &cli.CLI{
		Args:     args,
		Commands: Commands,
	}
	cliRunner.Run()
}
