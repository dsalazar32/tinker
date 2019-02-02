package main

import (
	"github.com/mitchellh/cli"
	"github.com/dsalazar32/tinker/command"
)

var Commands map[string]cli.CommandFactory
var Ui cli.Ui

const (
	ErrorPrefix  = "e:"
	OutputPrefix = "o:"
)

func initCommands() {
	meta := command.Meta {
		Ui: Ui,
	}

	Commands = map[string]cli.CommandFactory{
		"dotenv": func() (cli.Command, error) {
			return &command.DotenvCommand{
				Meta: meta,
			}, nil
		},
	}
}
