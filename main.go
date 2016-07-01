package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/go-yaml/yaml"
	//"github.com/subchen/goutils/cli"
)

const VERSION = "1.0.0"

var (
	BuildVersion   string
	BuildGitCommit string
	BuildDate      string
)

// template context
func newContext() map[string]interface{} {
	ctx := make(map[string]interface{})
	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		ctx[kv[0]] = kv[1]
	}
	return ctx
}

// template function
func defaultValue(a, b interface{}) interface{} {
	if a != nil {
		return a
	}
	return b
}

// flag Value
type StringList []string

func (s *StringList) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s *StringList) String() string {
	return strings.Join(*s, ",")
}

func executeTemplate(t *template.Template, file string, ctx interface{}, stdout bool, overwrite bool) {
	filePair := strings.SplitN(file, ":", 2)
	srcFile := filePair[0]
	destFile := ""
	if len(filePair) == 2 {
		destFile = filePair[1]
	} else {
		pos := strings.LastIndex(srcFile, ".")
		destFile = srcFile[0:pos]
	}

	tmpl, err := t.ParseFiles(srcFile)
	if err != nil {
		log.Fatalf("unable to parse template: %s", err)
	}

	dest := os.Stdout
	if !stdout {
		if !overwrite {
			if _, err := os.Stat(destFile); err == nil {
				log.Fatalf("cannot overwrite dest file: %s", destFile)
			}
		}

		dest, err = os.Create(destFile)
		if err != nil {
			log.Fatalf("unable to create %s", err)
		}
		defer dest.Close()
	}

	err = tmpl.ExecuteTemplate(dest, filepath.Base(srcFile), ctx)
	if err != nil {
		log.Fatalf("template error: %s\n", err)
	}
}

func main() {
	var (
		templatesFlag StringList
		delimsFlag    string
		stdoutFlag    bool
		overwriteFlag bool
		envsFlag      StringList
		jsonenvFlag   string
		loadenvFlag   string
		versionFlag   bool
	)

	flag.Var(&templatesFlag, "template", "Template (/template:/dest). can be passed multiple times")
	flag.StringVar(&delimsFlag, "delims", "", `Template tag delimiters. default "{{:}}" `)
	flag.BoolVar(&stdoutFlag, "stdout", false, "Output to console instead of file")
	flag.BoolVar(&overwriteFlag, "overwrite", false, "Overwrite file without errors if dest file exists")
	flag.Var(&envsFlag, "e", "Environment name=value pair, can be passed multiple times")
	flag.StringVar(&jsonenvFlag, "json", "", "load environment from json object")
	flag.StringVar(&loadenvFlag, "load", "", "load environment from json file")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
	flag.Parse()

	if versionFlag {
		fmt.Printf("Version: %s-%s\n", VERSION, BuildVersion)
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("Git commit: %s\n", BuildGitCommit)
		fmt.Printf("Built: %s\n", BuildDate)
		fmt.Printf("OS/Arch: %s-%s\n", runtime.GOOS, runtime.GOARCH)
		return
	}

	funcMap := template.FuncMap{
		"split":   strings.Split,
		"default": defaultValue,
	}

	t := template.New("noname").Funcs(funcMap)
	if delimsFlag != "" {
		delims := strings.Split(delimsFlag, ":")
		if len(delims) != 2 {
			log.Fatalf("bad delimiters argument: %s. expected \"left:right\"", delimsFlag)
		}
		t = t.Delims(delims[0], delims[1])
	}

	ctx := newContext()
	for _, env := range envsFlag {
		kv := strings.SplitN(env, "=", 2)
		ctx[kv[0]] = kv[1]
	}

	if jsonenvFlag != "" {
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(jsonenvFlag), &obj); err != nil {
			log.Fatalf("bad json format: %s", jsonenvFlag)
		}
		for name, value := range obj {
			ctx[name] = value
		}
	}

	if loadenvFlag != "" {
		if bytes, err := ioutil.ReadFile(loadenvFlag); err != nil {
			log.Fatalf("cannot load file: %s", loadenvFlag)
		} else {
			var obj map[string]interface{}
			if strings.HasSuffix(loadenvFlag, ".json") {
				if err := json.Unmarshal(bytes, &obj); err != nil {
					log.Fatalf("bad json format: %s", string(bytes))
				}
			} else if strings.HasSuffix(loadenvFlag, ".yaml") || strings.HasSuffix(loadenvFlag, ".yml") {
				if err := yaml.Unmarshal(bytes, &obj); err != nil {
					log.Fatalf("bad yaml format: %s", string(bytes))
				}
			}
			for name, value := range obj {
				ctx[name] = value
			}
		}
	}

	for _, file := range templatesFlag {
		executeTemplate(t, file, ctx, stdoutFlag, overwriteFlag)
	}
}
