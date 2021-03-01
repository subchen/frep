package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
	"github.com/go-yaml/yaml"
	"github.com/ismferd/ssm/package/parameterstore"
	"github.com/overdrive3000/secretsmanager"
)

func FuncMap(templateName string) template.FuncMap {
	f := sprig.TxtFuncMap()

	// marshal
	f["toBool"] = toBool
	f["toJson"] = toJson
	f["toToml"] = toToml
	f["toYaml"] = toYaml

	// file
	f["fileExists"] = fileExists
	f["fileSize"] = fileSize
	f["fileLastModified"] = fileLastModified
	f["fileGetBytes"] = fileGetBytes
	f["fileGetString"] = fileGetString

	// include
	f["include"] = include(templateName)

	// strings
	f["countRune"] = func(s string) int {
		return len([]rune(s))
	}

	// Fix sprig regex functions
	oRegexReplaceAll := f["regexReplaceAll"].(func(regex string, s string, repl string) string)
	oRegexReplaceAllLiteral := f["regexReplaceAllLiteral"].(func(regex string, s string, repl string) string)
	oRegexSplit := f["regexSplit"].(func(regex string, s string, n int) []string)
	f["reReplaceAll"] = func(regex string, replacement string, input string) string {
		return oRegexReplaceAll(regex, input, replacement)
	}
	f["reReplaceAllLiteral"] = func(regex string, replacement string, input string) string {
		return oRegexReplaceAllLiteral(regex, input, replacement)
	}
	f["reSplit"] = func(regex string, n int, input string) []string {
		return oRegexSplit(regex, input, n)
	}

	// Add function to get secrets from AWS Secrets Manager
	f["awsSecret"] = getAWSSecret
	f["awsParameterStore"] = getAWSParameterStore

	// base64 decode
	f["base64decode"] = base64Decode

	return f
}

// base64decode return a string encoded in base64
// On decode error will panic if in strict mode, otherwise returns empty string.
//
// This is designed to be called from a template.
func base64Decode(value string) string {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		if Strict {
			panic(err)
		}
		return ""
	}
	return string(decoded)
}

// getAWSSecret return a secret stored in AWS Secret Manager
// function accepts as parameter secret name and secret key.
// if secret key is not set then will return first key stored in
// secret.
func getAWSSecret(secret ...string) string {
	var name, key string
	c := secretsmanager.New(
		&secretsmanager.AWSConfig{},
	)

	name = secret[0]

	if len(secret) == 2 {
		key = secret[1]
	}

	spec := &secretsmanager.SecretSpec{
		Name: name,
		Key:  key,
	}

	s, err := c.GetSecret(spec)
	if err != nil {
		if Strict {
			panic(err)
		}
		return ""
	}

	return s
}

// getAWSParameterStore return a parameter stored in AWS SSM Parameter Store.
// function accepts as parameter a names.
func getAWSParameterStore(parameter string) string {

	c := parameterstore.New(
		&parameterstore.AWSConfig{},
	)

	spec := &parameterstore.ParemeterString{
		Name: parameter,
	}

	p, err := c.GetParam(spec)
	if err != nil {
		if Strict {
			panic(err)
		}
		return ""
	}

	return p
}

// toBool takes a string and converts it to a bool.
// On marshal error will panic if in strict mode, otherwise returns false.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
//
// This is designed to be called from a template.
func toBool(value string) bool {
	result, err := strconv.ParseBool(value)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return false
	}
	return result
}

// toJson takes an interface, marshals it to json, and returns a string.
// On marshal error will panic if in strict mode, otherwise returns empty string.
//
// This is designed to be called from a template.
func toJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return ""
	}
	return string(data)
}

// toYaml takes an interface, marshals it to yaml, and returns a string.
// On marshal error will panic if in strict mode, otherwise returns empty string.
//
// This is designed to be called from a template.
func toYaml(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return ""
	}
	return string(data)
}

// toToml takes an interface, marshals it to toml, and returns a string.
// On marshal error will panic if in strict mode, otherwise returns empty string.
//
// This is designed to be called from a template.
func toToml(v interface{}) string {
	b := bytes.NewBuffer(nil)
	e := toml.NewEncoder(b)
	err := e.Encode(v)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return ""
	}
	return b.String()
}

func fileExists(file string) bool {
	_, err := os.Stat(file)

	return err == nil
}

func fileSize(file string) int64 {
	info, err := os.Stat(file)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return 0
	}
	return info.Size()
}

func fileLastModified(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return time.Unix(0, 0)
	}
	return info.ModTime()
}

func fileGetBytes(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return []byte{}
	}
	return data
}

func fileGetString(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if Strict {
			panic(err.Error())
		}
		return ""
	}
	return string(data)
}

type relativeIncludeFunc func(include string) string

func include(callingFile string) relativeIncludeFunc {
	filePair := strings.SplitN(callingFile, ":", 2)
	callingFile = filePair[0]

	return func(includedFile string) string {
		if !path.IsAbs(includedFile) {
			includedFile = path.Join(path.Dir(callingFile), includedFile)
		}

		t := template.New(includedFile)
		t.Delims(delims[0], delims[1])
		t.Funcs(FuncMap(includedFile))

		var err error
		var templateBytes []byte

		templateBytes, err = ioutil.ReadFile(includedFile)
		if err != nil {
			if Strict {
				panic(fmt.Errorf("unable to read from %v, caused:\n\n   %v\n", includedFile, err))
			}
			return ""
		}

		tmpl, err := t.Parse(string(templateBytes))
		if err != nil {
			if Strict {
				panic(fmt.Errorf("unable to parse template file, caused:\n\n   %v\n", err))
			}
			return ""
		}

		var output bytes.Buffer
		err = tmpl.Execute(&output, ctx)

		if err != nil {
			if Strict {
				panic(fmt.Errorf("render template error, caused:\n\n   %v\n", err))
			}
			return ""
		}

		return output.String()
	}

}
