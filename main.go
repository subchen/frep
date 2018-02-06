package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/go-yaml/yaml"
	"github.com/subchen/go-cli"
)

// version
var (
	BuildVersion   string
	BuildGitRev    string
	BuildGitCommit string
	BuildDate      string
)

// flags
var (
	EnvironList  []string
	JsonStr      string
	LoadFileList []string
	Overwrite    bool
	Dryrun       bool
	Delims       string
)

// create template context
func newTemplateVariables() map[string]interface{} {
	vars := make(map[string]interface{})

	// Env
	envs := make(map[string]interface{})
	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		envs[kv[0]] = kv[1]
		vars[kv[0]] = kv[1] // legacy: use env in root scope
	}
	vars["Env"] = envs

	// --json
	if JsonStr != "" {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(JsonStr), &obj); err != nil {
			panic(fmt.Errorf("bad json format: %v", err))
		}
		for k, v := range obj {
			vars[k] = v
		}
	}

	// --load
	for _, file := range LoadFileList {
		if bytes, err := ioutil.ReadFile(file); err != nil {
			panic(fmt.Errorf("cannot load file, caused:\n\n   %v\n", err))
		} else {
			var obj map[string]interface{}
			if strings.HasSuffix(file, ".json") {
				if err := json.Unmarshal(bytes, &obj); err != nil {
					panic(fmt.Errorf("bad json format, caused:\n\n   %v\n", err))
				}
			} else if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
				if err := yaml.Unmarshal(bytes, &obj); err != nil {
					panic(fmt.Errorf("bad yaml format, caused:\n\n   %v\n", err))
				}
			} else if strings.HasSuffix(file, ".toml") {
				if err := toml.Unmarshal(bytes, &obj); err != nil {
					panic(fmt.Errorf("bad toml format, caused:\n\n   %v\n", err))
				}
			} else {
				panic(fmt.Errorf("bad file type: %s", file))
			}

			for k, v := range obj {
				vars[k] = v
			}
		}
	}

	// --env
	for _, env := range EnvironList {
		kv := strings.SplitN(env, "=", 2)

		// remove quotes for key="value"
		v := kv[1]
		if strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"") {
			v = v[1 : len(v)-1]
		} else if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			v = v[1 : len(v)-1]
		}
		vars[kv[0]] = v
	}

	return vars
}

func templateExecute(t *template.Template, file string, ctx interface{}) {
	filePair := strings.SplitN(file, ":", 2)
	srcFile := filePair[0]
	destFile := ""

	if len(filePair) == 2 {
		destFile = filePair[1]
	} else {
		if pos := strings.LastIndex(srcFile, "."); pos == -1 {
			destFile = srcFile
		} else {
			destFile = srcFile[0:pos]
		}
	}

	tmpl, err := t.ParseFiles(srcFile)
	if err != nil {
		panic(fmt.Errorf("unable to parse template file, caused:\n\n   %v\n", err))
	}

	dest := os.Stdout
	if !Dryrun {
		if !Overwrite {
			if _, err := os.Stat(destFile); err == nil {
				panic(fmt.Errorf("unable overwrite destination file: %s", destFile))
			}
		}

		dest, err = os.Create(destFile)
		if err != nil {
			panic(fmt.Errorf("unable to create file, caused:\n\n   %v\n", err))
		}
		defer dest.Close()
	}

	err = tmpl.ExecuteTemplate(dest, filepath.Base(srcFile), ctx)
	if err != nil {
		panic(fmt.Errorf("render template error, caused:\n\n   %v\n", err))
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "frep"
	app.Usage = "Generate file using template"
	app.UsageText = "[options] input-file[:output-file] ..."
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"

	app.Flags = []*cli.Flag{
		{
			Name:        "e, env",
			Usage:       "set variable name=value, can be passed multiple times",
			Placeholder: "name=value",
			Value:       &EnvironList,
		},
		{
			Name:        "json",
			Usage:       "load variables from json object string",
			Placeholder: "jsonstring",
			Value:       &JsonStr,
		},
		{
			Name:        "load",
			Usage:       "load variables from json/yaml/toml file",
			Placeholder: "file",
			Value:       &LoadFileList,
		},
		{
			Name:  "overwrite",
			Usage: "overwrite if destination file exists",
			Value: &Overwrite,
		},
		{
			Name:  "dryrun",
			Usage: "just output result to console instead of file",
			Value: &Dryrun,
		},
		{
			Name:     "delims",
			Usage:    "template tag delimiters",
			DefValue: "{{:}}",
			Value:    &Delims,
		},
	}

	app.Examples = strings.TrimSpace(`
frep nginx.conf.in -e webroot=/usr/share/nginx/html -e port=8080
frep nginx.conf.in:/etc/nginx.conf -e webroot=/usr/share/nginx/html -e port=8080
frep nginx.conf.in --json '{"webroot": "/usr/share/nginx/html", "port": 8080}'
frep nginx.conf.in --load config.json --overwrite
`)

	if BuildVersion != "" {
		app.Version = BuildVersion + "-" + BuildGitRev
	}
	app.BuildGitCommit = BuildGitCommit
	app.BuildDate = BuildDate

	app.Action = func(c *cli.Context) {
		if c.NArg() == 0 {
			c.ShowHelp()
			return
		}

		defer func() {
			if err := recover(); err != nil {
				os.Stderr.WriteString(fmt.Sprintf("fatal: %v\n", err))
				os.Exit(1)
			}
		}()

		t := template.New("noname").Funcs(FuncMap())
		if Delims != "" {
			pairs := strings.Split(Delims, ":")
			if len(pairs) != 2 {
				panic(fmt.Errorf("bad delimiters argument: %s. expected \"left:right\"", Delims))
			}
			t = t.Delims(pairs[0], pairs[1])
		}

		vars := newTemplateVariables()

		for _, file := range c.Args() {
			templateExecute(t, file, vars)
		}
	}

	app.Run(os.Args)
}
