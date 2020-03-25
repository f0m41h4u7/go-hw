package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"bufio"
	"regexp"
	"sort"
	"strings"
)

var validWord = regexp.MustCompile(`([A-Za-zА-Яа-я]+[^\s,.)(:;'"!?+=/\\]*)`)

func Top10(text string) []string {
	// If string is empty, return 0
	if text == "" {
		return nil
	}

	// Split string to words
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)
	var words []string
	for scanner.Scan() {
		valid := validWord.FindStringSubmatch(scanner.Text())
		if len(valid) != 0 {
			words = append(words, valid[0])
		}
	}

	// Count frequences of words
	freqs := make(map[string]int)
	for _, w := range words {
		w = strings.ToLower(w)
		freqs[w]++
	}

	// Convert map with frequences into slice of structs
	type wordStat struct {
		word string
		freq int
	}

	top := make([]wordStat, len(words))
	for word, freq := range freqs {
		top = append(top, wordStat{word, freq})
	}

	// Sort words by frequences
	sort.Slice(top, func(i, j int) bool {
		if top[i].freq == top[j].freq {
			return top[i].word < top[j].word
		}
		return top[i].freq > top[j].freq
	})

	// Pick first 10
	res := make([]string, 10)
	for i := range res {
		res[i] = top[i].word
	}

	return res
}
