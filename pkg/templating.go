package pkg

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// Match struct - specifies what should be replaced.
type Match struct {
	original string
	new string
}

// Get action from a match.
func (match Match) getAction() (string, error) {
	stripped := strings.Trim(match.original, "{} ")
	fields := strings.Fields(stripped)
	if len(fields) == 0 {
		return "", GenericError{"Empty bars"}
	}
	return fields[0], nil
}

// Get payload from a match.
func (match Match) getPayload() string {
	stripped := strings.Trim(match.original, "{} ")
	fields := strings.SplitN(stripped, " ", 2)
	if len(fields) < 2 {
		return ""
	}
	return fields[1]
}

// Get all {{.+}}s, 
func getMatches(template string) []Match {
	// Match to all of the {{.+}}s
	re := regexp.MustCompile(`{{.+}}`)
	raw_matches := re.FindAll([]byte(template), -1)

	// Iterate through and make Match structs.
	matches := make([]Match, len(raw_matches))
	for i, m := range raw_matches {
		matches[i] = Match{original: string(m)}
	}
	return matches
}

// Populate a match "new" fields based on the action.
func populateMatches(matches []Match, data Data, globalData Data) ([]Match, error) {
	var err error
	replaced := make([]Match, len(matches))
	for i, match := range matches {
		replaced[i] = match
		action, a_err := match.getAction()
		Check(a_err)
		switch action {
		case "title":
			replaced[i].new, err = ExtractData("title", data)
			Check(err)
		case "data":
			replaced[i].new, err = ExtractData(match.getPayload(), data)
			Check(err)
		case "global":
			replaced[i].new, err = ExtractData(match.getPayload(), globalData)
			Check(err)
		case "loop":
			fmt.Println("Managing loop!")
			replaced[i].new = "LOOP PLACEHOLDER"
		case "content":
			fmt.Println("Managing content!")
			replaced[i].new = "CONTENT PLACEHOLDER"
		default:
			fmt.Println("Unknown action encountered")
			replaced[i].new = replaced[i].original
		}
	}
	return replaced, nil
}

// Replace a match's original with it's new in a template.
func replaceMatches(template string, matches []Match) (string, error) {
	ret := template
	for _, match := range matches {
		ret = strings.ReplaceAll(ret, match.original, match.new)
	}
	return ret, nil
}

// Apply data to a template.
func ProcessData(template string, data Data, globalData Data) string {
	matches := getMatches(template)
	matches, p_err := populateMatches(matches, data, globalData)
	Check(p_err)
	ret, r_err := replaceMatches(template, matches)
	Check(r_err)
	return ret
}

// Apply partials to a template.
func ProcessPartials(template string, partials map[string]string) string {
	// Get matches, then iterate through them.
	matches := getMatches(template)
	replaced := make([]Match, len(matches))
	for i, match := range matches {
		replaced[i] = match
		action, a_err := match.getAction()
		Check(a_err)
		switch action {
		case "partial":
			// Apply partial; replace with partial file contents.
			new, err := ioutil.ReadFile(partials[match.getPayload()])
			Check(err)
			replaced[i].new = string(new)
		default:
			// Default; do no replacement.
			replaced[i].new = match.original
		}
	}
	ret, r_err := replaceMatches(template, replaced)
	Check(r_err)
	return ret
}
