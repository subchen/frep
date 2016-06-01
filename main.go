package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const VERSION = "1.0.0"

// template context
func newContext() map[string]string {
	ctx := make(map[string]string)
	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		ctx[kv[0]] = kv[1]
	}
	return ctx
}

// template function
func defaultValue(a, b interface{}) string {
	if a != nil {
		return a.(string)
	}
	return b.(string)
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
	tmpl, err := t.ParseFiles(file)
	if err != nil {
		log.Fatalf("unable to parse template: %s", err)
	}

	dest := os.Stdout
	if !stdout {
		pos := strings.LastIndex(file, ".")
		destFile := file[0:pos]

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

	err = tmpl.ExecuteTemplate(dest, filepath.Base(file), ctx)
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
		versionFlag   bool
	)

	flag.Var(&templatesFlag, "t", "Template (/template:/dest). can be passed multiple times")
	flag.StringVar(&delimsFlag, "delims", "", `Template tag delimiters. default "{{:}}" `)
	flag.BoolVar(&stdoutFlag, "stdout", false, "Output to console instead of file")
	flag.BoolVar(&overwriteFlag, "overwrite", false, "Overwrite file without errors if dest file exists")
	flag.Var(&envsFlag, "e", "Environment name=value pair, can be passed multiple times")
	flag.BoolVar(&versionFlag, "version", false, "Show version")
	flag.Parse()

	if versionFlag {
		fmt.Println(VERSION)
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

	for _, file := range templatesFlag {
		executeTemplate(t, file, ctx, stdoutFlag, overwriteFlag)
	}
}
