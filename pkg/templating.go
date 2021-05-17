package pkg

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Match struct {
	original string
	action string
	payload string
	new string
}

func getAction(match string) (string, error) {
	stripped := strings.Trim(match, "{} ")
	fields := strings.Fields(stripped)
	if len(fields) == 0 {
		return "", GenericError{"Empty payload."}
	}
	return fields[0], nil
}

func getPayload(match string) (string, error) {
	stripped := strings.Trim(match, "{} ")
	fields := strings.SplitN(stripped, " ", 2)
	if len(fields) < 2{
		return "", GenericError{"Empty payload."}
	}
	return fields[1], nil
}

func getMatches(template string) []Match {
	// Match to all of the {{}}s
	re := regexp.MustCompile(`{{*}}`)
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

	// Return
	return matches
}

func getMatchesByAction(template string, action string) []Match {
	// Match to all of the {{}}s
	re := regexp.MustCompile(`{{*}}`)
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
		if curr_action == action {
			matches[i] = Match{original: str_match, action: curr_action}
		}
	}

	// Return
	return matches
}

// TODO: This is definitely sketchy
// TODO: Define a new interface called "DataPoint", which can either be a
// string or a map? The whole thing should be a map, at most nested two?
// Gosh idk...
func getData(payload string, data map[string]interface{}) string {
	fields := strings.Split(payload, ".")
	var curr interface{}
	curr = data
	for _, field := range fields {
		switch curr.(type) {
		case map[string]interface{}:
			curr = map[string]interface{}(curr)[field]
		}
	}
	return string(curr)
}

func populateMatches(matches []Match, data map[string]interface{}) ([]Match, error) {
	replaced := make([]Match, len(matches))
	for i, match := range matches {
		replaced[i] = match
		replaced[i].new = getData(match.payload, data)
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

func ProcessData(template string, data map[string]interface{}) string {
	fmt.Println("Proccessing Data")
	matches := getMatchesByAction(template, "data")
	matches, p_err := populateMatches(matches, data)
	Check(p_err)
	ret, r_err := replaceMatches(template, matches)
	Check(r_err)
	return ret
}
