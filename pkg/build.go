package pkg

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
	"github.com/otiai10/copy"
)

// Config struct. Built from site.yaml.
type Config struct {
	Name        string
	Output      string
	Draft		bool
	Static_src      string
	Static_dest      string
	Global      map[string]string
	Partials    map[string]string
	Pages       []Page
	Collections []Collection
}

// Parse them site.yaml file to get site Config.
func parseConfig(file string) Config {
	// Read config out to a string.
	data, io_err := ioutil.ReadFile(file)
	Check(io_err)

	// Unmarshal YAML into config.
	config := Config{}
	y_err := yaml.Unmarshal([]byte(data), &config)
	Check(y_err)
	return config
}

// Main Build function.
func Build(file string) {
	// Sanity print.
	fmt.Println("Building Goo site.")

	// Parse config, setup Dir, get globalData.
	config := parseConfig(file)
	os.RemoveAll(config.Output)
	os.MkdirAll(config.Output, 0744)
	globalData := GetData(config.Global)

	// Write pages.
	for _, page := range config.Pages {
		WritePage(page, config, globalData)
	}

	// Write collections.
	for _, collection := range config.Collections {
		WriteCollection(collection, config, globalData)
	}

	// Copy static folder over.
	copy.Copy(config.Static_src, config.Output+"/"+config.Static_dest)

	// Sanity print
	fmt.Println("Goo site built at '" + config.Output + "'")
}
