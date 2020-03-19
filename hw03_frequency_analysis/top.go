package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(text string) []string {
	// If string is empty, return 0
	if text == "" {
		return nil
	}

	// Separation regex
	sep := regexp.MustCompile(`[[:blank:]\,\.\?\!\~\/\(\)\*\+\;\:\\\n']+`)
	// Split string to words
	words := sep.Split(text, -1)

	// Count frequences of words
	freqs := make(map[string]uint)
	for _, w := range words {
		matched, _ := regexp.Compile(`^[[[:punct:]]|\s|\n]+$`, w)
		if (!matched) && (w != "") {
			w = strings.ToLower(w)
			_, ok := freqs[w]
			if ok {
				freqs[w]++
			} else {
				freqs[w] = 1
			}
		}
	}

	// Convert map with frequences into slice of structs
	type mapStruct struct {
		Key   string
		Value uint
	}
	top := make([]mapStruct, len(words))
	for key, val := range freqs {
		top = append(top, mapStruct{key, val})
	}

	// Sort words by frequences
	sort.Slice(top, func(i, j int) bool {
		if top[i].Value == top[j].Value {
			return top[i].Key < top[j].Key
		}
		return top[i].Value > top[j].Value
	})

	// Pick first 10
	res := make([]string, 10)
	for i := range res {
		res[i] = top[i].Key
	}
	return res
}
