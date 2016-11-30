package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/go-yaml/yaml"
	"github.com/subchen/goutils/cli"
	"github.com/Masterminds/sprig"
)

const VERSION = "1.1.0"

var (
	BuildVersion   string
	BuildGitCommit string
	BuildDate      string
)

// create template context
func newTemplateVariables(ctx *cli.Context) map[string]interface{} {
	// ENV
	vars := make(map[string]interface{})
	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		vars[kv[0]] = kv[1]
	}

	// --json
	if jsonStr := ctx.String("--json"); jsonStr != "" {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
			cli.Fatalf("fatal: bad json format: %v", err)
		}
		for k, v := range obj {
			vars[k] = v
		}
	}

	// --load
	for _, file := range ctx.StringList("--load") {
		if bytes, err := ioutil.ReadFile(file); err != nil {
			cli.Fatalf("fatal: cannot load file: %s", file)
		} else {
			var obj map[string]interface{}
			if strings.HasSuffix(file, ".json") {
				if err := json.Unmarshal(bytes, &obj); err != nil {
					cli.Fatalf("fatal: bad json format: %v", err)
				}
			} else if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
				if err := yaml.Unmarshal(bytes, &obj); err != nil {
					cli.Fatalf("fatal: bad yaml format: %v", err)
				}
			} else {
				cli.Fatalf("fatal: bad file type: %s", file)
			}

			for k, v := range obj {
				vars[k] = v
			}
		}
	}

	// --env
	for _, env := range ctx.StringList("--env") {
		kv := strings.SplitN(env, "=", 2)
		vars[kv[0]] = kv[1]
	}

	return vars
}

func templateExecute(t *template.Template, file string, ctx interface{}, testing bool, overwrite bool) {
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
		cli.Fatalf("fatal: unable to parse template: %v", err)
	}

	dest := os.Stdout
	if !testing {
		if !overwrite {
			if _, err := os.Stat(destFile); err == nil {
				cli.Fatalf("fatal: unable overwrite destination file: %s", destFile)
			}
		}

		dest, err = os.Create(destFile)
		if err != nil {
			cli.Fatalf("fatal: unable to create file: %v", err)
		}
		defer dest.Close()
	}

	err = tmpl.ExecuteTemplate(dest, filepath.Base(srcFile), ctx)
	if err != nil {
		cli.Fatalf("fatal: transform template error: %v", err)
	}
}

func cliExecute(ctx *cli.Context) {
	t := template.New("noname").Funcs(sprig.TxtFuncMap())
	if delimsStr := ctx.String("--delims"); delimsStr != "" {
		delims := strings.Split(delimsStr, ":")
		if len(delims) != 2 {
			cli.Fatalf("fatal: bad delimiters argument: %s. expected \"left:right\"", delimsStr)
		}
		t = t.Delims(delims[0], delims[1])
	}

	vars := newTemplateVariables(ctx)

	for _, file := range ctx.Args() {
		testing := ctx.Bool("--testing")
		overwrite := ctx.Bool("--overwrite")
		templateExecute(t, file, vars, testing, overwrite)
	}
}

func main() {
	app := cli.NewApp("frep", "Transform template file using environment, arguments, json/yaml files")

	app.Flag("-e, --env", "set variable name=value, can be passed multiple times").Multiple()
	app.Flag("--json", "load variables from json object").Placeholder("string")
	app.Flag("--load", "load variables from json/yaml files").Multiple()
	app.Flag("--overwrite", "overwrite if destination file exists").Bool()
	app.Flag("--testing", "test mode, output transform result to console").Bool()
	app.Flag("--delims", `template tag delimiters`).Default("{{:}}")

	if BuildVersion == "" {
		app.Version = VERSION
	} else {
		app.Version = func() {
			fmt.Printf("Version: %s-%s\n", VERSION, BuildVersion)
			fmt.Printf("Go version: %s\n", runtime.Version())
			fmt.Printf("Git commit: %s\n", BuildGitCommit)
			fmt.Printf("Built: %s\n", BuildDate)
			fmt.Printf("OS/Arch: %s-%s\n", runtime.GOOS, runtime.GOARCH)
		}
	}

	app.Usage = func() {
		fmt.Println("Usage: frep [OPTIONS] input-file:[output-file] ...")
		fmt.Println("   or: frep [ --version | --help ]")
	}

	app.MoreHelp = func() {
		fmt.Println("Examples:")
		fmt.Println("  frep nginx.conf.in -e webroot=/usr/share/nginx/html -e port=8080")
		fmt.Println("  frep nginx.conf.in:/etc/nginx.conf -e webroot=/usr/share/nginx/html -e port=8080")
		fmt.Println("  frep nginx.conf.in --json '{\"webroot\": \"/usr/share/nginx/html\", \"port\": 8080}'")
		fmt.Println("  frep nginx.conf.in --load context.json --overwrite")
	}

	app.AllowArgumentCount(1, -1)

	app.Execute = cliExecute

	app.Run()
}
