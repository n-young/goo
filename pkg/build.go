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
	Name   string
	Output string
	Static string
	Global map[string]string
	Partials map[string]string
	Pages  []Page
}

// Page struct. For each in pages.
type Page struct {
	Title    string
	Path     string
	Template string
	Data  	 map[string]string
}

// Gets the real path.
func (p Page) getPath(config Config) (string, error) {
	switch {
	case len(p.Path) == 0 || string(p.Path[0]) != "/":
		// If path doesn't start with "/", it's malformed.
		return "", GenericError{"Path did not begin with \"/\"."}
	case p.Path == "/":
		// If root path, throw in the index.html.
		return config.Output + "/index.html", nil
	case string(p.Path[len(p.Path)-1]) == "/":
		// Removing trailing slash.
		return config.Output + p.Path[:len(p.Path)-1] + ".html", nil
	default:
		// Regular case.
		return config.Output + p.Path + ".html", nil
	}
}

// Manifests a page. Fills fields, converts to string, and outputs to file.
func (p Page) manifest(config Config, globalData Data) string {
	// Read out template file.
	bytes, io_err := ioutil.ReadFile(p.Template)
	Check(io_err)
	ret := string(bytes)

	// Fill out partials, if necessary.
	ret = ProcessPartials(ret, config.Partials)

	// Fill out data, if necessary.
	if p.Data != nil {
		data := GetData(p.Data)
		data.(DataNode).setTitle(p.Title)
		ret = ProcessData(ret, data, globalData)
	}

	// Copy static folder over, return.
	copy.Copy(config.Static, config.Output + "/" + config.Static)
	return ret
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

// Write to file.
func writePage(page Page, config Config, globalData Data) {
	// Get data and path
	data := page.manifest(config, globalData)
	path, path_err := page.getPath(config)
	Check(path_err)

	// Write to file
	wr_err := ioutil.WriteFile(path, []byte(data), 0644)
	Check(wr_err)
}

// Main Build function.
func Build(file string) {
	// Sanity print
	fmt.Println("Building Goo site.")

	// Parse config, setup Dir, get globalData
	config := parseConfig(file)
	os.RemoveAll(config.Output)
	os.MkdirAll(config.Output, 0744)
	globalData := GetData(config.Global)

	// Write pages
	for _, page := range config.Pages {
		writePage(page, config, globalData)
	}

	// Sanity print
	fmt.Println("Goo site built at " + config.Output)
}
