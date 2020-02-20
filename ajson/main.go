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

func usage(force bool) {
	text := ``
	if force || inArgs("-h", "-help") {
		// 		if inArgs("update") {
		// 			text = `
		// Usage: ajson [-input=...] update "what" "how"
		//   Update some fields with evaluated data. NB! This method couldn't create new fields, only update existing.
		// Options:
		//   -input     Path to the JSON file. Leave it blank to use STDIN.
		// Examples:
		// 	ajson -input=input.json update "$..price" "@ * 1.13"
		// `
		// 		} else
		if inArgs("jsonpath") {
			text = `
Usage: ajson [-input=...] jsonpath "..."
  Evaluate JSONPath and print the result.
Options:
  -input     Path to the JSON file. Leave it blank to use STDIN.
Examples:
	ajson -input=input.json jsonpath "$..price"
	ajson -input=http://example.com/file.json jsonpath "$..success"
	cat "filename.json" | ajson jsonpath "$..[?(@.price)]"
`
		} else if inArgs("eval") {
			text = `
Usage: ajson [-input=...] eval "..."
  Evaluate JSONPath as only value and print the result. 
Options:
  -input     Path to the JSON file. Leave it blank to use STDIN.
Examples:
	ajson -input=input.json eval "avg($..price)"
	ajson -input=http://example.com/file.json eval "$.result.length"
	cat "filename.json" | ajson eval "round(avg($..price))"
`
		} else {
			// 			text = `
			// Usage: ajson [-input=...] action ...
			//   Read JSON and evaluate it with JSONPath.
			// Actions:
			//   jsonpath   Evaluate JSONPath and print the result (Example: "$..price").
			//   eval       Evaluate JSONPath as only value and print the result (Example: "avg($..price)").
			//   update     Update some fields with evaluated data (Example: "$..price" "@ * 1.13").
			// Options:
			//   -input     Path to the JSON file. Leave it blank to use STDIN.
			// `
			text = `
Usage: ajson [-input=...] action ...
  Read JSON and evaluate it with JSONPath.
Actions:
  jsonpath   Evaluate JSONPath and print the result (Example: "$..price"). 
  eval       Evaluate JSONPath as only value and print the result (Example: "avg($..price)"). 
Options:
  -input     Path to the JSON file. Leave it blank to use STDIN.
`
		}
	}
	if text != "" {
		fmt.Println(text)
		os.Exit(2)
	}
}

func main() {
	usage(false)
	action := getAction()
	input := getInput()
	data, err := ioutil.ReadAll(input)
	_ = input.Close()
	if err != nil {
		log.Fatalf("error reading source: %s", err)
	}
	var result *ajson.Node
	switch action {
	case "jsonpath":
		root, err := ajson.Unmarshal(data)
		if err != nil {
			log.Fatalf("error parsing JSON: %s", err)
		}

		paths := fromArgs("jsonpath", 1)
		if paths[0] == "" {
			usage(true)
		}

		nodes, err := root.JSONPath(paths[0])
		if err != nil {
			log.Fatalf("error: %s", err)
		}

		result = ajson.ArrayNode("", nodes)
	case "eval":
		root, err := ajson.Unmarshal(data)
		if err != nil {
			log.Fatalf("error parsing JSON: %s", err)
		}

		paths := fromArgs("eval", 1)
		if paths[0] == "" {
			usage(true)
		}

		result, err = ajson.Eval(root, paths[0])
		if err != nil {
			log.Fatalf("error: %s", err)
		}
	default:
		usage(true)
	}

	res, err := ajson.Marshal(result)
	if err != nil {
		log.Fatalf("error preparing JSON: %s", err)
	}
	fmt.Printf("%s", res)
}

func getAction() string {
	actions := map[string]bool{
		"jsonpath": inArgs("jsonpath"),
		"eval":     inArgs("eval"),
		// "update":   inArgs("update"),
	}
	action := ""
	for current, ok := range actions {
		if ok {
			if action == "" {
				action = current
			} else {
				log.Fatal("selected more than one action")
			}
		}
	}
	if action == "" {
		usage(true)
	}
	return action
}

func getInput() io.ReadCloser {
	input := ""
	flag.StringVar(&input, "input", "", "")
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

func fromArgs(action string, count int) []string {
	result := make([]string, count)
	for i, val := range os.Args {
		if action == val {
			for j := i + 1; j < len(os.Args) && j < i+1+count; j++ {
				result[j-(i+1)] = os.Args[j]
			}
		}
	}
	return result
}
