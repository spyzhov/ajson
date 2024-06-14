package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spyzhov/ajson"
)

var version = "v0.9.3"

func usage() {
	text := ``
	if inArgs("-h", "-help", "--help", "help") {
		text = `Usage: ajson [-mq] "jsonpath" ["input"]
  Read JSON and evaluate it with JSONPath.
Parameters:
  -m, --multiline  Input file/stream will be read as a multiline JSON. Each line should have a full valid JSON.
  -q, --quiet      Do not print errors into the STDERR.
  -man             Display man page with "man ajson -man"
Argument:
  jsonpath         Valid JSONPath or evaluate string (Examples: "$..[?(@.price)]", "$..price", "avg($..price)")
  input            Path to the JSON file. Leave it blank to use STDIN.
Examples:
  ajson "avg($..registered.age)" "https://randomuser.me/api/?results=5000"
  ajson "$.results.*.name" "https://randomuser.me/api/?results=10"
  curl -s "https://randomuser.me/api/?results=10" | ajson "$..coordinates"
  ajson "$" example.json
  echo "3" | ajson "2 * pi * $"
  docker logs image-name -f | ajson -m 'root($[?(@=="ERROR" && key(@)=="severity")])'`
	} else if inArgs("version", "-version", "--version") {
		text = fmt.Sprintf(`ajson: Version %s
Copyright (c) 2020 Pyzhov Stepan
MIT License <https://github.com/spyzhov/ajson/blob/master/LICENSE>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.`, version)
	}
	if text != "" {
		fmt.Println(text)
		os.Exit(2)
	}
}

func main() {
	log.SetFlags(0)
	usage()
	cfg := getConfig()
	if cfg.jsonpath == "" {
		pFatal("JSONPath was not set")
	}
	input := getInput(cfg)
	defer func() {
		_ = input.Close()
	}()

	if cfg.multiline {
		reader := bufio.NewReader(input)
		line := 1
		for {
			data, err := reader.ReadBytes('\n')
			if err != nil {
				if !errors.Is(err, io.EOF) {
					mlError(cfg, "unable to read input at line %d: %s", line, err)
				}
				return
			}
			apply(cfg, data, line)
			line++
		}
	} else {
		data, err := io.ReadAll(input)
		if err != nil {
			pFatal("error reading source: %s", err)
		}
		apply(cfg, data, 0)
	}
}

func apply(cfg config, data []byte, line int) {
	var result *ajson.Node
	msg := "error"
	if cfg.multiline {
		msg = fmt.Sprintf("error at line %d", line)
	}

	root, err := ajson.Unmarshal(data)
	if err != nil {
		mlFatal(cfg, "%s parsing JSON: %s", msg, err)
		return
	}

	var nodes []*ajson.Node
	nodes, err = root.JSONPath(cfg.jsonpath)
	if err != nil { // try to eval
		result, err = ajson.Eval(root, cfg.jsonpath)
	} else {
		result = ajson.ArrayNode("", nodes)
	}
	if err != nil {
		mlFatal(cfg, "jsonpath %s: %s", msg, err)
		return
	}

	if cfg.multiline {
		if (result.IsArray() || result.IsObject()) && result.Empty() {
			return
		}
		if result.IsString() && result.MustString() == "" {
			return
		}
		if result.IsNull() {
			return
		}
	}

	data, err = ajson.Marshal(result)
	if err != nil {
		mlFatal(cfg, "%s preparing JSON: %s", msg, err)
	}
	pPrint("%s\n", data)
}

func getInput(cfg config) io.ReadCloser {
	if cfg.input == "" {
		return os.Stdin
	}

	if strings.HasPrefix(cfg.input, "http://") || strings.HasPrefix(cfg.input, "https://") {
		resp, err := http.DefaultClient.Get(cfg.input)
		if err != nil {
			pFatal("Error on getting data from '%s': %s", cfg.input, err)
		}
		if resp.StatusCode >= 300 {
			if !cfg.quiet {
				pError("WARNING: status code is '%s'", resp.Status)
			}
		}
		return resp.Body
	}

	file, err := os.Open(cfg.input)
	if err != nil {
		pFatal("Error on open file '%s': %s", cfg.input, err)
	}

	return file
}

func inArgs(value ...string) bool {
	index := make(map[string]bool, len(value))
	for _, val := range value {
		index[val] = true
	}
	args := os.Args
	for _, val := range args {
		if index[val] {
			return true
		}
	}
	return false
}

type config struct {
	jsonpath  string
	input     string
	multiline bool
	quiet     bool
}

func getConfig() (cfg config) {
	for _, val := range os.Args[1:] {
		switch val {
		case "-m", "--multiline":
			cfg.multiline = true
		case "-q", "--quiet":
			cfg.quiet = true
		case "-mq", "-qm":
			cfg.multiline = true
			cfg.quiet = true
		default:
			if cfg.jsonpath == "" {
				cfg.jsonpath = val
			} else if cfg.input == "" {
				cfg.input = val
			} else {
				pFatal("Wrong arguments count, unknown flag %q", val)
			}
		}
	}
	return
}

func pPrint(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

func pError(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func pFatal(format string, args ...interface{}) {
	pError(format, args...)
	os.Exit(1)
}

func mlFatal(cfg config, format string, args ...interface{}) {
	if cfg.multiline {
		if !cfg.quiet {
			pError(format, args...)
		}
	} else {
		pFatal(format, args...)
	}
}

func mlError(cfg config, format string, args ...interface{}) {
	if !cfg.quiet {
		pError(format, args...)
	}
}
