package pkg

import (
	"io/ioutil"
)

// Page struct. For each in pages.
type Page struct {
	Title    string
	Path     string
	Template string
	Data     map[string]string
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
// NOTE: Could probably be merged into the WritePage function.
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
		data.(DataNode).setGlobal(globalData)
		ret = ProcessData(ret, data)
	}
	return ret
}

// Write page to file.
func WritePage(page Page, config Config, globalData Data) {
	// Get data and path
	data := page.manifest(config, globalData)
	path, path_err := page.getPath(config)
	Check(path_err)

	// Write to file
	wr_err := ioutil.WriteFile(path, []byte(data), 0644)
	Check(wr_err)
}
