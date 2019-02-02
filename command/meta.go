package command

import (
	"flag"
	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

type Meta struct {
	input     bool
	variables map[string]interface{}

	Ui cli.Ui
}

func (m *Meta) flagSet(n string) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)
	return f
}

func (m *Meta) mergeArgsFromConfigs(key string, args []string) []string {
	for elem, val := range viper.GetStringMap(key) {
		args = append(args, "-"+elem, val.(string))
	}
	args = append(args, args...)
	return args
}
