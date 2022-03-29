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
	new      string
}

// Get action from a match.
func (match Match) getToken(n int) (string, error) {
	stripped := strings.Trim(match.original, "{}")
	stripped = strings.Trim(stripped, " ")
	fields := strings.Fields(stripped)
	if len(fields) < n {
		return "", GenericError{"Invalid token call. " + match.original}
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
	matches, m_err := getMatches(template, `\${`, `}`)
	Check(m_err)
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
		// First, check if the payload has any matches.
		curr_payload := payload
		matches, m_err := getMatches(payload, `{{`, `}}`)
		Check(m_err)
		if matches != nil {
			curr_payload = ProcessData(payload, data)
		}
		//
		entry, r_err := processTemplate(curr_payload, data)
		Check(r_err)
		ret += entry
	}
	return ret, nil
}

// Get all matches to a string,
func getMatches(template string, ldel string, rdel string) ([]Match, error) {
	// Match to all left and right delims
	re := regexp.MustCompile("(" + ldel + "|" + rdel + ")")
	match_indices := re.FindAllStringIndex(template, -1)
	if match_indices == nil {
		return nil, nil
	}

	// Create raw_matches from indices
	raw_matches := make([]string, 0)
	count := 0
	curr_li := 0
	for _, mi := range match_indices {
		// In the beginning, set the current running left to the current left if it is a TLD
		if count == 0 {
			curr_li = mi[0]
		}
		// Count up or down depending on the current delimiter.
		// TODO: Handle other delim edge cases.
		trimmed_ldel := strings.Trim(ldel, "\\")
		trimmed_rdel := strings.Trim(rdel, "\\")
		switch curr_delim := template[mi[0]:mi[1]]; curr_delim {
		case trimmed_ldel:
			count++
		case trimmed_rdel:
			count--
		default:
			return nil, GenericError{"Bad delimiter."}
		}

		// If we reach a negative delimiter, we love. Otherwise, if it's 0, append the running string.
		if count < 0 {
			return nil, GenericError{"Unmatched paren 1."}
		} else if count == 0 {
			raw_matches = append(raw_matches, template[curr_li:mi[1]])
		}
	}
	// If we don't end up at zero, there must have been an unmatched paren.
	if count != 0 {
		return nil, GenericError{"Unmatched paren 2."}
	}

	// Iterate through and make Match structs.
	matches := make([]Match, len(raw_matches))
	for i, m := range raw_matches {
		matches[i] = Match{original: m}
	}
	return matches, nil
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
			list, l_err := ExtractList(subaction, data)
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
	matches, m_err := getMatches(template, `{{`, `}}`)
	Check(m_err)
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
	matches, m_err := getMatches(template, `{{`, `}}`)
	Check(m_err)
	matches, p_err := populateMatches(matches, data)
	Check(p_err)
	ret, r_err := replaceMatches(template, matches)
	Check(r_err)
	return ret
}

// Apply content to a template.
func ProcessContent(template string, content string) string {
	// Get matches, then iterate through them.
	matches, m_err := getMatches(template, `{{`, `}}`)
	Check(m_err)
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
