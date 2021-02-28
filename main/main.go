package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gophercises/urlshort"
)

var ymlFile string
var jsnFile string

func init() {
	flag.StringVar(&ymlFile, "yamlFile", "", "YAML file name having 'path' and respective redirection 'url'.")
	flag.StringVar(&ymlFile, "y", "", "YAML file name having 'path' and respective redirection 'url' - Short Form.")

	flag.StringVar(&jsnFile, "jsonFile", "", "JSON file name having array of 'path' and respective redirection 'url'.")
	flag.StringVar(&jsnFile, "j", "", "JSON file name having array of 'path' and respective redirection 'url' - Short Form.")
}

func main() {
	flag.Parse()
	fmt.Printf("ymlFile: %v\n", ymlFile)
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	var yaml string

	if ymlFile != "" {
		fileData, err := os.ReadFile(ymlFile)
		if err != nil {
			log.Fatalf("Unable to read file: %v", ymlFile)
		}

		yaml = string(fileData)
	} else {
		yaml = `
			- path: /urlshort
			url: https://github.com/gophercises/urlshort
			- path: /urlshort-final
			url: https://github.com/gophercises/urlshort/tree/solution
			`
	}

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the YAMLHandler as the
	// fallback
	var jsn string

	if jsnFile != "" {
		fileData, err := os.ReadFile(jsnFile)
		if err != nil {
			log.Fatalf("Unable to read file: %v", jsnFile)
		}

		jsn = string(fileData)
	} else {
		jsn = `
		[
			{
			   "path":"/yahoo",
			   "url":"https://yahoo.com/"
			},
			{
			   "path":"/stackoverflow",
			   "url":"https://stackoverflow.com/"
			}
		 ]
			`
	}

	jsonHandler, err := urlshort.JSONHandler([]byte(jsn), yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
