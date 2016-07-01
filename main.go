package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/go-yaml/yaml"
	"github.com/subchen/goutils/cli"
)

const VERSION = "1.0.0"

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
			log.Fatalf("bad json format: %s", jsonStr)
		}
		for k, v := range obj {
			vars[k] = v
		}
	}

	// --load
	for _, file := range ctx.StringList("--load") {
		if bytes, err := ioutil.ReadFile(file); err != nil {
			log.Fatalf("cannot load file: %s", file)
		} else {
			var obj map[string]interface{}
			if strings.HasSuffix(file, ".json") {
				if err := json.Unmarshal(bytes, &obj); err != nil {
					log.Fatalf("bad json format: %s", string(bytes))
				}
			} else if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
				if err := yaml.Unmarshal(bytes, &obj); err != nil {
					log.Fatalf("bad yaml format: %s", string(bytes))
				}
			} else {
				log.Fatalf("bad file type: %s", file)
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

// builtin template function
func defaultValue(a, b interface{}) interface{} {
	if a == nil {
		return b
	}
	if s, ok := a.(string); ok && s == "" {
		return b
	}
	return a
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
		log.Fatalf("unable to parse template: %s", err)
	}

	dest := os.Stdout
	if !testing {
		if !overwrite {
			if _, err := os.Stat(destFile); err == nil {
				log.Fatalf("unable overwrite destination file: %s", destFile)
			}
		}

		dest, err = os.Create(destFile)
		if err != nil {
			log.Fatalf("unable to create file: %s", err)
		}
		defer dest.Close()
	}

	err = tmpl.ExecuteTemplate(dest, filepath.Base(srcFile), ctx)
	if err != nil {
		log.Fatalf("transform template error: %s\n", err)
	}
}

func cliExecute(ctx *cli.Context) {
	funcMap := template.FuncMap{
		"split":   strings.Split,
		"default": defaultValue,
	}

	t := template.New("noname").Funcs(funcMap)
	if delimsStr := ctx.String("--delims"); delimsStr != "" {
		delims := strings.Split(delimsStr, ":")
		if len(delims) != 2 {
			log.Fatalf("bad delimiters argument: %s. expected \"left:right\"", delimsStr)
		}
		t = t.Delims(delims[0], delims[1])
	}

	vars := newTemplateVariables(ctx)

	for _, file := range ctx.Args() {
		testing := ctx.Bool("--test")
		overwrite := ctx.Bool("--overwrite")
		templateExecute(t, file, vars, testing, overwrite)
	}
}

func main() {
	app := cli.NewApp("frep", "transform template file using environment, arguments, json/yaml files")

	app.Flag("-e, --env", "Set variable name=value, can be passed multiple times").Multiple()
	app.Flag("--test", "Test mode, output transform result to console").Bool()
	app.Flag("--overwrite", "Overwrite if destination file exists").Bool()
	app.Flag("--delims", `Template tag delimiters`).Default("{{:}}")
	app.Flag("--json", "Load variables from json object").Placeholder("string")
	app.Flag("--load", "Load variables from json/yaml files").Placeholder("file").Multiple()

	app.Version = func() {
		fmt.Printf("Version: %s-%s\n", VERSION, BuildVersion)
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("Git commit: %s\n", BuildGitCommit)
		fmt.Printf("Built: %s\n", BuildDate)
		fmt.Printf("OS/Arch: %s-%s\n", runtime.GOOS, runtime.GOARCH)
	}

	app.Execute = cliExecute

	app.Run()
}
