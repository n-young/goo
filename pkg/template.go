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
func (match Match) getToken(n int) (string, error) {
	stripped := strings.Trim(match.original, "{}")
	stripped = strings.Trim(stripped, " ")
	fields := strings.Fields(stripped)
	if len(fields) < n {
		return "", GenericError{"Invalid token call."}
	}
	return fields[n-1], nil
}

// Get payload from a match.
func (match Match) getPayload(n int) string {
	stripped := strings.Trim(match.original, "{} ")
	fields := strings.SplitN(stripped, " ", n+1)
	if len(fields) < n+1 {
		return ""
	}
	return fields[n]
}

// Process a single template.
func processTemplate(template string, data Data) (string, error) {
	matches := getMatches(template, `\${(.|\n)+?}`)
		replaced := make([]Match, len(matches))
		for i, match := range matches {
			replaced[i] = match
			stripped := strings.Trim(match.original, "${}")
			new, e_err := ExtractData(stripped, data)
			Check(e_err)
			replaced[i].new = new
		}
		return replaceMatches(template, replaced)
}

// Get and replace all ${.+}s
func populateTemplates(payload string, list []Data) (string, error) {
	ret := ""
	for _, data := range list {
		entry, r_err := processTemplate(payload, data)
		Check(r_err)
		ret += entry
	}
	return ret, nil
}

// Get all matches to a string, 
func getMatches(template string, match string) []Match {
	// Match to all.
	re := regexp.MustCompile(match)
	raw_matches := re.FindAll([]byte(template), -1)

	// Iterate through and make Match structs.
	matches := make([]Match, len(raw_matches))
	for i, m := range raw_matches {
		matches[i] = Match{original: string(m)}
	}
	return matches
}

// Populate a match "new" fields based on the action.
func populateMatches(matches []Match, data Data) ([]Match, error) {
	var err error
	replaced := make([]Match, len(matches))
	for i, match := range matches {
		
		replaced[i] = match
		action, a_err := match.getToken(1)
		Check(a_err)
		switch action {
		case "title":
			replaced[i].new, err = ExtractData("title", data)
			Check(err)
		case "data":
			replaced[i].new, err = ExtractData(match.getPayload(1), data)
			Check(err)
		case "template":
			replaced[i].new, err = processTemplate(match.getPayload(1), data)
			Check(err)
		case "loop":
			subaction, s_err := match.getToken(2)
			Check(s_err)
			list, l_err :=  ExtractList(subaction, data)
			Check(l_err)
			replaced[i].new, err = populateTemplates(match.getPayload(2), list)
			Check(err)
		default:
			fmt.Println("Invalid action encountered")
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

// Apply partials to a template.
func ProcessPartials(template string, partials map[string]string) string {
	// Get matches, then iterate through them.
	matches := getMatches(template, `{{(.|\n)+?}}`)
	replaced := make([]Match, len(matches))
	for i, match := range matches {
		replaced[i] = match
		action, a_err := match.getToken(1)
		Check(a_err)
		switch action {
		case "partial":
			// Apply partial; replace with partial file contents.
			new, err := ioutil.ReadFile(partials[match.getPayload(1)])
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

// Apply data to a template.
func ProcessData(template string, data Data) string {
	matches := getMatches(template, `{{(.|\n)+?}}`)
	matches, p_err := populateMatches(matches, data)
	Check(p_err)
	ret, r_err := replaceMatches(template, matches)
	Check(r_err)
	return ret
}

// Apply content to a template.
func ProcessContent(template string, content string) string {
	// Get matches, then iterate through them.
	matches := getMatches(template, `{{(.|\n)+?}}`)
	replaced := make([]Match, len(matches))
	for i, match := range matches {
		replaced[i] = match
		action, a_err := match.getToken(1)
		Check(a_err)
		switch action {
		case "content":
			// Apply partial; replace with partial file contents.
			replaced[i].new = content
		default:
			// Default; do no replacement.
			replaced[i].new = match.original
		}
	}
	ret, r_err := replaceMatches(template, replaced)
	Check(r_err)
	return ret
}
