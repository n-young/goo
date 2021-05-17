package pkg

import (
	"fmt"
	"io/ioutil"
    "os"

	"gopkg.in/yaml.v2"
)

// Config struct. Built from site.yaml.
type Config struct {
	Name  string
	Output string
	Static string
	Pages  []Page
}


// Page struct. For each in pages.
type Page struct {
	Title    string
	Path     string
	Template string
	Data  	 string
}

// Gets the real path.
func (p Page) getPath(config Config) (string, error) {
	switch {
	case len(p.Path) == 0 || string(p.Path[0]) != "/":
		return "", GenericError{"Path did not begin with \"/\"."}
	case p.Path == "/":
		return config.Output + "/index.html", nil
	case string(p.Path[len(p.Path)-1]) == "/":
		return config.Output + p.Path[:len(p.Path)-1] + ".html", nil
	default:
		return config.Output + p.Path + ".html", nil
	}
}

// Gets data as a map.
func (p Page) getData() map[string]interface{} {
	// Reads data out into a string.
	data, io_err := ioutil.ReadFile(p.Data)
	Check(io_err)

	// Unmarshal into m.
	m := make(map[string]interface{})
	y_err := yaml.Unmarshal(data, &m)
	Check(y_err)
	return m
}

// Manifests a page. Fills fields, converts to string, and outputs to file.
func (p Page) manifest(config Config) string {
	// Read out template file.
	bytes, io_err := ioutil.ReadFile(p.Template)
	Check(io_err)
	ret := string(bytes)

	// TODO: Fill out partials, if necessary.
	// ret = ProcessPartials(ret, config)

	// TODO: Fill out data, if necessary.
	// TODO: Make it such that if there are data 
	if p.Data != "" {
		data := p.getData()
		ret = ProcessData(ret, data)
	}

	// TODO: Copy static folder over

	// Return.
	return ret
}


// Parse them site.yaml file to get site Config.
func parseConfig(file string) Config {
	// Read config out to a string
	data, io_err := ioutil.ReadFile(file)
	Check(io_err)

	// Unmarshal YAML into c.
	config := Config{}
	y_err := yaml.Unmarshal([]byte(data), &config)
	Check(y_err)
	return config
}

// Write to file.
func writePage(page Page, config Config) {
	// Get data and path
	data := page.manifest(config)
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

	// Parse config
	config := parseConfig(file)

	// Setup dir
	os.RemoveAll(config.Output)
	os.MkdirAll(config.Output, 0744)

	// Write pages
	for _, page := range config.Pages {
		writePage(page, config)
	}
}
