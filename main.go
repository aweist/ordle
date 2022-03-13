package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

const (
	DictionaryDir = "english3.txt"
)

// func SolutionFromFile(filename string) (result []string, err error) {
// 	return Solution(parseInput(filename))
// }

func QuordleSolutions(states []State) (results [][]string) {
	for _, s := range states {
		results = append(results, Solution(s))
	}
	return
}

func Solution(state State) (result []string) {
	d, err := buildDict(DictionaryDir)
	if err != nil {
		log.Fatal(err)
	}

	allMisplaced := state.AllMisplaced()

	result = []string{}
	var check func(int, []byte)
	check = func(i int, answer []byte) {
		if i == 5 {
			s := string(answer)
			if d[s] {
				required := copyMap(allMisplaced)
				for _, b := range answer {
					delete(required, b)
				}
				if len(required) == 0 {
					result = append(result, s)
				}
			}
			return
		}

		if c, ok := state.IsKnown(i); ok {
			answer[i] = c
			check(i+1, answer)
		} else {
			for c := byte('a'); c <= byte('z'); c++ {
				if !state.IsWrong(c) && !state.IsMisplaced(i, c) {
					answer[i] = c
					check(i+1, answer)
				}
			}
		}
	}
	a := make([]byte, 5)
	check(0, a)
	return
}

func copyMap(a map[byte]bool) map[byte]bool {
	result := map[byte]bool{}
	for i := range a {
		result[i] = a[i]
	}
	return result
}

func buildDict(dictFile string) (dict map[string]bool, err error) {
	dict = map[string]bool{}

	// open file
	f, err := os.Open(dictFile)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()
	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// do something with read line
		w := rec[0]
		w = strings.ToLower(w)
		if len(w) == 5 {
			dict[w] = true
		}
	}

	return
}
