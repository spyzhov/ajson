package main

import (
	"flag"
	"fmt"
	"github.com/spyzhov/ajson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func usage() {
	text := ``
	if inArgs("-h", "-help", "--help", "help") {
		text = `Usage: ajson [-input=...] [-eval] "jsonpath"
  Read JSON and evaluate it with JSONPath.
Options:
  -input     Path to the JSON file. Leave it blank to use STDIN.
  -eval      Evaluate JSONPath as only value and print the result (Example: "avg($..price)").
Argument:
  jsonpath   Valid JSONPath or evaluate string (Examples: "$..[?(@.price)]", "$..price", "avg($..price)")
`
	}
	if text != "" {
		fmt.Println(text)
		os.Exit(2)
	}
}

func main() {
	log.SetFlags(0)
	usage()
	path := os.Args[len(os.Args)-1]
	if len(os.Args) < 2 || path == "" || strings.HasPrefix(path, "-") {
		log.Fatalf("JSONPath was not set")
	}
	isEval := inArgs("-eval")
	input := getInput()
	data, err := ioutil.ReadAll(input)
	_ = input.Close()
	if err != nil {
		log.Fatalf("error reading source: %s", err)
	}
	var result *ajson.Node

	root, err := ajson.Unmarshal(data)
	if err != nil {
		log.Fatalf("error parsing JSON: %s", err)
	}
	if isEval {
		result, err = ajson.Eval(root, path)
	} else {
		var nodes []*ajson.Node
		nodes, err = root.JSONPath(path)
		result = ajson.ArrayNode("", nodes)
	}
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	res, err := ajson.Marshal(result)
	if err != nil {
		log.Fatalf("error preparing JSON: %s", err)
	}
	fmt.Printf("%s\n", res)
}

func getInput() io.ReadCloser {
	input := ""
	flag.StringVar(&input, "input", "", "Path to the JSON file. Leave it blank to use STDIN")
	flag.Bool("eval", false, "Evaluate JSONPath as only value and print the result (Example: \"avg($..price)\")")
	flag.Parse()
	if input == "" {
		return os.Stdin
	}
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
