package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aweist/ordle/parse"
)

const (
	DictionaryDir = "english3.txt"
)

type parser func(io.Reader) []parse.State

func main() {
	filteTypeUsage := fmt.Sprintf("Filetype must be one of [%s]", strings.Join(validFileTypes(), ","))
	fileType := flag.String("t", "", filteTypeUsage)

	filenameUsage := "Path to input file.  File should be the .html file generated by using File -> Save Page As... within chrome"
	filename := flag.String("f", "", filenameUsage)

	flag.Parse()

	switch *fileType {
	case "dordle":
		SolveDordle(*filename)
	case "quordle":
		SolveQuordle(*filename)
	case "octordle":
		SolveOctordle(*filename)
	default:
		flag.PrintDefaults()
	}

}

func validFileTypes() []string {
	return []string{
		// "wordle",
		"dordle",
		"quordle",
		"octordle",
	}
}

func SolveDordle(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	states := parse.ParseDordle(f)
	printResults(Solutions(states))
}

func printResults(results [][]string) {
	for i, result := range results {
		log.Println("Results for", i)
		for _, r := range result {
			log.Println("  ", r)
		}
	}
}

func SolveOctordle(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	states := parse.ParseOctordle(f)
	printResults(Solutions(states))
}

func SolveQuordle(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	states := parse.ParseQuordle(f)
	printResults(Solutions(states))
}

func Solutions(states []parse.State) (results [][]string) {
	for _, s := range states {
		results = append(results, Solution(s))
	}
	return
}

func Solution(state parse.State) (result []string) {
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
