package urlGen

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

const (
	arg            = "(/[^{}/ ]+)" // Single "name"
	preffix        = arg + "*"     // Any combination of `/name)` without spaces
	mandatoryVar   = "/{[^{}/ ]+}"
	optionalVar    = mandatoryVar + "[?]?"
	varLong        = arg + optionalVar
	varsLong       = "(" + varLong + ")*"
	varsShort      = "(" + optionalVar + ")*"
	validLongUrl   = "^" + preffix + varsLong + "$"
	validShortUrl  = "^" + preffix + varsShort + "$"
	unclearVarName = "/(?:[{])[^{}/ :]+(?:[:}])"
)

func urlErrors(url string, isLong bool) error {
	var pattern *regexp.Regexp

	if isLong {
		pattern = regexp.MustCompile(validLongUrl)
	} else {
		pattern = regexp.MustCompile(validShortUrl)
	}
	if !pattern.MatchString(url) {
		return errors.New("Invalid Mux Patern")
	} else {
		var exists = make(map[string]bool)
		var varNames = regexp.MustCompile(unclearVarName)
		matches := varNames.FindAllString(url, -1)
		for _, s := range matches {
			s = strings.Replace(s, ":", "}", 1)
			if exists[s] {
				return errors.New("Error: Invalid Mux Path: Duplicated var name in url pattern")
			} else {
				exists[s] = true
			}
		}
		return nil
	}
}

func containsAll(url string, words []string) bool {
	for _, w := range words {
		if !strings.Contains(url, w) {
			return false
		}
	}
	return true
}

func shortVarFix(w string) string {
	var (
		wn          string
		unclearName = regexp.MustCompile(unclearVarName)
		getNameOnly = regexp.MustCompile("[^{}/ :]+")
	)
	wn = unclearName.FindString(w)
	wn = getNameOnly.FindString(w)
	return "/" + wn + w
}

func urlBuilder(url string, isLong bool) []string {
	var permutations []string
	var varsPattern *regexp.Regexp
	if isLong {
		varsPattern = regexp.MustCompile(varLong)
	} else {
		varsPattern = regexp.MustCompile(optionalVar)
	}

	firstMatchIndex := varsPattern.FindStringIndex(url)
	if firstMatchIndex != nil {
		var mandatoryVars []string = nil

		vars := varsPattern.FindAllString(url, -1)
		if !isLong {
			for i, w := range vars {
				vars[i] = shortVarFix(w)
			}
		}

		// Finding Mandatory Vars
		for i, w := range vars {
			cw := strings.TrimSuffix(w, "?")

			if !strings.HasSuffix(w, "?") {
				mandatoryVars = append(mandatoryVars, cw)
			}
			vars[i] = cw
		}

		permute(url[:firstMatchIndex[0]], vars, mandatoryVars, &permutations)

		// If there aren't mandatory vars, the "preffix" is a valid url
		if len(mandatoryVars) == 0 {
			permutations = append(permutations, url[:firstMatchIndex[0]])
		}

	} else {
		return []string{url}
	}
	return permutations
}

func permute(preffix string, alphabet []string, mustContain []string, permutations *[]string) {
	var newAlphabet []string = nil

	if len(alphabet) == 0 {
		return
	}
	for i, s := range alphabet {
		word := preffix + s

		newAlphabet = nil
		newAlphabet = append(newAlphabet, alphabet[:i]...)
		newAlphabet = append(newAlphabet, alphabet[i+1:]...)

		if containsAll(word, mustContain) {
			*permutations = append(*permutations, word)
		}

		permute(word, newAlphabet, mustContain, permutations)
	}
}

func GetUrlVarsPermutations(url string, isLong bool) []string {

	err := urlErrors(url, isLong)
	if err != nil {
		log.Panic(err)
	}
	return urlBuilder(url, isLong)
}
