package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
	"github.com/go-yaml/yaml"
)

func FuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	// marshal
	f["toJson"] = toJson
	f["toYaml"] = toYaml
	f["toToml"] = toToml
	f["toBool"] = toBool
	// file
	f["fileSize"] = fileSize
	f["fileLastModified"] = fileLastModified
	f["fileGetBytes"] = fileGetBytes
	f["fileGetString"] = fileGetString
	return f
}

// toBool takes a string and converts it to a bool. It will
// always return a bool, even on parsing error (false).
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
//
// This is designed to be called from a template.
func toBool(value string) bool {
	result, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return result
}

// toJson takes an interface, marshals it to json, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

// toYaml takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toYaml(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

// toToml takes an interface, marshals it to toml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toToml(v interface{}) string {
	b := bytes.NewBuffer(nil)
	e := toml.NewEncoder(b)
	err := e.Encode(v)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

func fileSize(file string) int64 {
	info, err := os.Stat(file)
	if err != nil {
		// Swallow errors inside of a template.
		return -1
	}
	return info.Size()
}

func fileLastModified(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		// Swallow errors inside of a template.
		return time.Unix(0, 0)
	}
	return info.ModTime()
}

func fileGetBytes(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		// Swallow errors inside of a template.
		return nil
	}
	return data
}

func fileGetString(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}
