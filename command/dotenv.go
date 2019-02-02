package command

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type DotenvCommand struct {
	Meta
}

type env struct {
	Key string
	Val string
}

func (c *DotenvCommand) Run(args []string) int {
	var denv, tmpl, outf string

	cmdName := "dotenv"
	cmdFlags := c.Meta.flagSet(cmdName)
	cmdFlags.StringVar(&denv, "in", ".env", "location of .env")
	cmdFlags.StringVar(&tmpl, "tmpl", "", "template to apply .env")
	cmdFlags.StringVar(&outf, "out", "stdout", "output file")
	if err := cmdFlags.Parse(c.Meta.mergeArgsFromConfigs(cmdName, args)); err != nil {
		return 1
	}

	if _, err := os.Stat(denv); err != nil {
		c.Ui.Error(fmt.Sprintf("[%s]: %v\n", cmdName, err))
		return 1
	}

	f, err := os.Open(denv)
	defer f.Close()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("[%s]: %v\n", cmdName, err))
		return 1
	}

	envs := []env{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		kv, err := parseDotEnv(scanner.Text())
		if err != nil {
			continue
		}

		envs = append(envs, env{kv[0], kv[1]})
	}

	// TODO: Break out template generation logic.
	if tmpl != "" {
		r, w := io.Pipe()
		go func() {
			defer w.Close()

			funcMap := template.FuncMap{
				"dec": func(i int) int {
					return i - 1
				},
			}
			outTmpl, err := template.New("parsedOut").Funcs(funcMap).Parse(tmpl)
			if err != nil {
				c.Ui.Error(fmt.Sprintf("parsing error: %v", err))
				os.Exit(1)
			}
			outTmpl.Execute(w, envs)
		}()

		b := &bytes.Buffer{}
		b.ReadFrom(r)
		c.Ui.Output(strings.TrimSpace(b.String()))
	} else {
		c.Ui.Output(fmt.Sprintf("%v", envs))
	}
	return 0
}

func (c *DotenvCommand) Help() string {
	return "Nothing to see here"
}

func (c *DotenvCommand) Synopsis() string {
	return "Not much to say"
}

func parseDotEnv(s string) ([]string, error) {
	// make sure line item matches the environmental variable declaration pattern:
	// export k=v
	// k=v
	valid := regexp.MustCompile(`^(export )?\w+=.+`)
	if !valid.MatchString(s) {
		return nil, errors.New("does not match .env pattern")
	}

	// extract key=value
	re := regexp.MustCompile(`\w+=.+`)
	s = re.FindAllString(s, -1)[0]

	// map[key]=value pair
	idx := strings.Index(s, "=")
	key, value := s[:idx], s[idx+1:]
	key = strings.TrimSpace(key)

	return []string{key, value}, nil
}
