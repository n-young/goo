package pkg

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	hashtag "github.com/abhinav/goldmark-hashtag"
	wikilink "github.com/abhinav/goldmark-wikilink"
	mermaid "github.com/abhinav/goldmark-mermaid"
)

// Collection struct. For each in collections.
type Collection struct {
	Title    string
	Path     string
	Template string
	Data     map[string]string
	Posts    string
}

func (c Collection) getBasePath(config Config) (string, error) {
	switch {
	case len(c.Path) == 0 || string(c.Path[0]) != "/":
		// If path doesn't start with "/", it's malformed.
		return "", GenericError{"Path did not begin with \"/\"."}
	case string(c.Path[len(c.Path)-1]) != "/":
		// Add trailing slash.
		return config.Output + c.Path + "/", nil
	default:
		// Regular case.
		return config.Output + c.Path, nil
	}
}

// Convert a given Markdown file to HTML
func markdownToHtml(filename string) (string, map[interface{}]interface{}) {
	// Define Markdown parsing options.
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			emoji.Emoji,
			mathjax.MathJax,
			meta.Meta,
			highlighting.NewHighlighting(
				// TODO: Make this pickable. https://github.com/alecthomas/chroma/tree/master/styles.
				highlighting.WithStyle("monokai"),
			),
			&hashtag.Extender{
				Variant: hashtag.ObsidianVariant,
			},
			&wikilink.Extender{},
			&mermaid.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	// Get file source.
	source, r_err := ioutil.ReadFile(filename)
	Check(r_err)

	// Read source into the buffer.
	context := parser.NewContext()
	var buf bytes.Buffer
	Check(md.Convert(source, &buf, parser.WithContext(context)))

	// Get and convert Metadata
	metadata_raw := meta.Get(context)
	metadata := make(map[interface{}]interface{})
	for k, v := range metadata_raw {
		metadata[k] = v
	}
	return buf.String(), metadata
}

// Write collection to dir.
func WriteCollection(c Collection, config Config, globalData Data) error {
	// Read out template file.
	bytes, io_err := ioutil.ReadFile(c.Template)
	Check(io_err)
	ret := string(bytes)

	// Fill out partials, if necessary.
	ret = ProcessPartials(ret, config.Partials)

	// Fill out data, if necessary.
	if c.Data != nil {
		data := GetData(c.Data)
		data.(DataNode).setTitle(c.Title)
		data.(DataNode).setGlobal(globalData)
		ret = ProcessData(ret, data)
	}

	// Ensure base path exists.
	base_path, bp_err := c.getBasePath(config)
	Check(bp_err)
	os.MkdirAll(base_path, 0744)

	// Iterate through the markdown files in the posts directory, writing each one.
	files, err := filepath.Glob(c.Posts + "/*")
	Check(err)
	for _, file := range files {
		// Convert to HTML, inject into template.
		content, metadata_raw := markdownToHtml(file)
		metadata_nested := make(map[interface{}]interface{})
		metadata_nested["post"] = metadata_raw
		metadata, err := CastData(metadata_nested)
		Check(err)

		// If its not a draft, and the post is a draft, then skip it.
		if !config.Draft {
			post_metadata, pm_err := metadata.getChild("post")
			Check(pm_err)
			is_draft_leaf, idl_err := post_metadata.getChild("draft")
			if idl_err == nil {
				is_draft, id_err := is_draft_leaf.getValue()
				Check(id_err)
				if is_draft == "true" {
					continue
				}
			}
		}

		// Process the content and metadata.
		processed := ProcessContent(ret, content)
		processed = ProcessData(processed, metadata)

		// Get correct path.
		tokens := strings.Split(file, "/")
		trailing_path := strings.TrimSuffix(tokens[len(tokens)-1], ".md")

		// Write
		wr_err := ioutil.WriteFile(base_path+trailing_path+".html", []byte(processed), 0644)
		Check(wr_err)
	}
	return nil
}
