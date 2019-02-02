package command

import (
	"fmt"
	"github.com/mitchellh/cli"
	"path/filepath"
	"reflect"
	"testing"
)

var testPath, _ = filepath.Abs("../fixtures")

func TestDotenvCommand(t *testing.T) {

	samples := []struct {
		info string
		tmpl string
		out  string
	}{
		{
			"Generate docker cli cmd.",
			"docker run -it -d{{range $idx, $ele := .}} -e {{$ele.Key}}={{$ele.Val}}{{end}} containerid",
			"docker run -it -d -e HELLO=world -e GOODNIGHT=moon containerid\n",
		},
		{
			"Generate chef data bag env object.",
			jsonTmpl,
			"{ \"HELLO\": \"world\", \"GOODNIGHT\": \"moon\" }\n",
		},
	}

	for _, s := range samples {
		ui := &cli.MockUi{}
		c := &DotenvCommand{
			Meta: Meta{
				Ui: ui,
			},
		}
		t.Log(s.info)
		args := []string{"-in", testPath + "/.env", "-tmpl", s.tmpl}
		if code := c.Run(args); code != 0 {
			t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
		} else {
			if s.out != ui.OutputWriter.String() {
				t.Fatalf("bad parse: \n want \n %s \n got \n %s", s.out, ui.OutputWriter.String())
			}
		}
	}
}

func TestDotenv_parseDotEnv(t *testing.T) {
	test := []struct {
		in  string
		out []string
	}{
		{
			"hello=world",
			[]string{"hello", "world"},
		},
		{
			"export something=new",
			[]string{"something", "new"},
		},
		{
			"export hiphip=hor#ray",
			[]string{"hiphip", "hor#ray"},
		},
		{
			"#commented=out",
			nil,
		},
		{
			"export #commented=out",
			nil,
		},
		{
			"#export commented=out",
			nil,
		},
		{
			"export =nokey",
			nil,
		},
		{
			"export fiz baz=fizbaz",
			nil,
		},
		{
			"=nokey",
			nil,
		},
		{
			"fiz baz=fizbaz",
			nil,
		},
	}

	errmsg := "does not match .env pattern"
	for _, try := range test {
		kv, err := parseDotEnv(try.in)
		if !reflect.DeepEqual(kv, try.out) {
			t.Fatalf("want %s but got %v:", try.out, kv)
		}
		if err != nil {
			if fmt.Sprintf("%v", err) != errmsg {
				t.Fatalf("want %v but got %v", errmsg, err)
			}
		}
	}
}

const jsonTmpl = `{ 
{{- $total := dec (len .)}}
{{- range $idx, $ele := .}} "{{$ele.Key}}": "{{$ele.Val}}"{{if ne $idx $total}},{{end}}{{end}} }`
