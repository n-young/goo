package pkg

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type Match struct {
	original string
	action string
	payload string
	new string
}

func getMatches(template string) []Match {
	// Match to all of the {{}}s
	re := regexp.MustCompile(`{{.+}}`)
	raw_matches := re.FindAll([]byte(template), -1)

	// Iterate through and make Match structs.
	matches := make([]Match, len(raw_matches))
	for i, match := range raw_matches {
		str_match := string(match)
		curr_action, a_err := getAction(str_match)
		Check(a_err)
		curr_payload, p_err := getPayload(str_match)
		Check(p_err)
		matches[i] = Match{original: str_match, action: curr_action, payload: curr_payload}
	}
	return matches
}

func getAction(match string) (string, error) {
	stripped := strings.Trim(match, "{} ")
	fields := strings.Fields(stripped)
	return fields[0], nil
}

func getPayload(match string) (string, error) {
	stripped := strings.Trim(match, "{} ")
	fields := strings.SplitN(stripped, " ", 2)
	if len(fields) < 2 {
		return "", nil
	}
	return fields[1], nil
}

func populateMatches(matches []Match, data Data, globalData Data) ([]Match, error) {
	var err error
	replaced := make([]Match, len(matches))
	for i, match := range matches {
		replaced[i] = match
		switch match.action {
		case "title":
			replaced[i].new, err = ExtractData("title", data)
			Check(err)
		case "data":
			replaced[i].new, err = ExtractData(match.payload, data)
			Check(err)
		case "global":
			replaced[i].new, err = ExtractData(match.payload, globalData)
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

func replaceMatches(template string, matches []Match) (string, error) {
	ret := template
	for _, match := range matches {
		ret = strings.ReplaceAll(ret, match.original, match.new)
	}
	return ret, nil
}

func ProcessData(template string, data Data, globalData Data) string {
	matches := getMatches(template)

	matches, p_err := populateMatches(matches, data, globalData)
	Check(p_err)

	ret, r_err := replaceMatches(template, matches)
	Check(r_err)
	return ret
}

func ProcessPartials(template string, partials map[string]string) string {
	matches := getMatches(template)
	replaced := make([]Match, len(matches))
	for i, match := range matches {
		replaced[i] = match
		switch match.action {
		case "partial":
			new, err := ioutil.ReadFile(partials[match.payload])
			Check(err)
			replaced[i].new = string(new)
		default:
			replaced[i].new = match.original
		}
	}
	ret, r_err := replaceMatches(template, replaced)
	Check(r_err)
	return ret
}
