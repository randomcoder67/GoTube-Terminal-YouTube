package download

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// This file contains functions to manage the playlist frecency

// Date is stored as the number of days since the epoch, to make calulations easier
func getDate() string {
	return strconv.Itoa(int(time.Now().Unix() / 86400))
}

func readFile(fileName string) string {
	dat, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(dat)
}

func fileLen(contents string) int {
	return strings.Count(contents, "\n")
}

// https://stackoverflow.com/questions/73925647/how-to-split-a-string-based-on-second-occurence-of-delimiter-in-go (241223)
func split(s string, sep rune, n int) string {
	for i, sep2 := range s {
		if sep2 == sep {
			n--
			if n == 0 {
				return s[:i]
			}
		}
	}
	return s
}

func saveFile(contents, fileName string) {
	err := os.WriteFile(fileName, []byte(contents), 0666)
	if err != nil {
		panic(err)
	}
}

func collapseSlice(input []string, sep string) string {
	var toReturn string = ""
	for _, part := range input {
		toReturn = toReturn + part + sep
	}
	return toReturn[:len(toReturn)-len(sep)]
}

func GetFrecencyData(fileName string) [][]string {
	var contents string = readFile(fileName)
	var toReturn [][]string
	var lines []string = strings.Split(contents, "\n")
	for _, line := range lines {
		var parts []string = strings.Split(line, " ")
		if len(parts) < 3 {
			break
		}
		toReturn = append(toReturn, []string{parts[0], parts[1], collapseSlice(parts[2:], " ")})
	}
	return toReturn
}

func AddToFile(playlistId, playlistName string, fileName string) {
	var existingContents string = readFile(fileName)
	if fileLen(existingContents) >= 999 {
		existingContents = split(existingContents, '\n', 999) + "\n"
	}

	var editedContents string = fmt.Sprintf("%s %s %s\n", getDate(), playlistId, playlistName) + existingContents
	saveFile(editedContents, fileName)
}

func getScore(date string, curDate string) int {
	dateNum, _ := strconv.Atoi(date)
	curDateNum, _ := strconv.Atoi(curDate)
	var gap int = curDateNum - dateNum

	switch {
	case gap < 6:
		return 100
	case gap < 11:
		return 30
	case gap < 16:
		return 10
	case gap < 21:
		return 5
	default:
		return 1
	}
}

func GetTopN(fileName string, n int) [][]string {
	var frecencyData [][]string = GetFrecencyData(fileName)
	var curDate string = getDate()
	_, _ = curDate, frecencyData

	var frecency map[string]int = make(map[string]int)
	var names map[string]string = make(map[string]string)

	for _, entry := range frecencyData {
		if _, ok := frecency[entry[1]]; ok {
			frecency[entry[1]] += getScore(entry[0], curDate)
		} else {
			frecency[entry[1]] = getScore(entry[0], curDate)
			names[entry[1]] = entry[2]
		}
	}

	// https://stackoverflow.com/questions/18695346/how-can-i-sort-a-mapstringint-by-its-values (241223)

	keys := make([]string, 0, len(frecency))
	for key := range frecency {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return frecency[keys[i]] > frecency[keys[j]] })

	var results [][]string

	for i, key := range keys {
		if i >= n {
			return results
		}
		results = append(results, []string{key, names[key]})
	}
	return results
}
