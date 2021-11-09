package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spyzhov/ajson/v1"
)

var version = "v0.6.0"

func usage() {
	text := ``
	if inArgs("-h", "-help", "--help", "help") || len(os.Args) > 3 {
		text = `Usage: ajson "jsonpath" ["input"]
  Read JSON and evaluate it with JSONPath.
Argument:
  jsonpath   Valid JSONPath or evaluate string (Examples: "$..[?(@.price)]", "$..price", "avg($..price)")
  input      Path to the JSON file. Leave it blank to use STDIN.
Examples:
  ajson "avg($..registered.age)" "https://randomuser.me/api/?results=5000"
  ajson "$.results.*.name" "https://randomuser.me/api/?results=10"
  curl -s "https://randomuser.me/api/?results=10" | ajson "$..coordinates"
  ajson "$" example.json
  echo "3" | ajson "2 * pi * $"`
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
	if len(os.Args) < 2 {
		log.Fatalf("JSONPath was not set")
	}
	path := os.Args[1]
	input := getInput()
	defer func() {
		_ = input.Close()
	}()
	data, err := ioutil.ReadAll(input)

	if err != nil {
		log.Fatalf("error reading source: %s", err)
	}
	var result *ajson.Node

	root, err := ajson.Unmarshal(data)
	if err != nil {
		log.Fatalf("error parsing JSON: %s", err)
	}

	var nodes []*ajson.Node
	nodes, err = root.JSONPath(path)
	result = ajson.NewArray(nodes)
	if err != nil {
		result, err = ajson.Eval(root, path)
	}
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	data, err = ajson.Marshal(result)
	if err != nil {
		log.Fatalf("error preparing JSON: %s", err)
	}
	fmt.Printf("%s\n", data)
}

func getInput() io.ReadCloser {
	if len(os.Args) < 3 {
		return os.Stdin
	}

	input := os.Args[2]
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		resp, err := http.DefaultClient.Get(input)
		if err != nil {
			log.Fatalf("Error on getting data from '%s': %s", input, err)
		}
		if resp.StatusCode >= 400 {
			log.Printf("WARNING: status code is '%s'", resp.Status)
		}
		return resp.Body
	}
	file, err := os.Open(input)
	if err != nil {
		log.Fatalf("Error on open file '%s': %s", input, err)
	}
	return file
}

func inArgs(value ...string) bool {
	index := make(map[string]bool, len(value))
	for _, val := range value {
		index[val] = true
	}
	for _, val := range os.Args {
		if index[val] {
			return true
		}
	}
	return false
}
