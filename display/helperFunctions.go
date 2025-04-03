package display

import (
	"fmt"
	"gotube/youtube"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"gotube/config"
)

// This file is for helper functions which do not interact directly with tcell. For functions which do, seem drawHelpers.go for general tcell functions and thing related to the MainContent interface, and gridHelper.go for tcell functions specific to drawing/managing the grids

// Debug function
func Print(str string) {
	cmd := exec.Command("notify-send", str)
	cmd.Run()
}

func SetCurrentSearchTerm(input string) {
	currentSearchTerm = input
	cursorLoc = len(input)
}

// Get the rune at num position (can't just use string[i], because of > 1 byte unicode characters
func getNChar(input string, num int) rune {
	return []rune(input)[num]
}

// Return slice of string up to num position, multi-byte-rune friendly
func toNSlice(input string, num int) string {
	var sb strings.Builder
	var runes []rune = []rune(input)

	if num > len(runes) {
		return ""
	}

	for i := 0; i < num; i++ {
		sb.WriteRune(runes[i])
	}

	return sb.String()
}

// Return slice of string from num position, multi-byte-rune friendly
func fromNSlice(input string, num int) string {
	var sb strings.Builder
	var runes []rune = []rune(input)
	var totalLen int = len(runes)

	if num > len(runes) {
		return ""
	}

	for i := num; i < totalLen; i++ {
		sb.WriteRune(runes[i])
	}

	return sb.String()
}

// Get length of string in runes, correctly accounting for multi byte characters
func properLen(input string) int {
	return (utf8.RuneCountInString(input))
}

// Return a slice containing the keys in a map
func sliceFromMap[T any](input map[string]T) []string {
	var toReturn []string
	for key := range input {
		toReturn = append(toReturn, key)
	}
	return toReturn
}

// Trim title to fit in 2 lines according to box width
func trimTitle(title string, channel string, cellWidth int) string {
	var lineLen int = 35 + ((cellWidth - 1) * 39)
	var totalLen int = lineLen * 2
	var finalString string = ""

	// If no trimming is needed
	if properLen(title) + properLen(channel) + 3 <= totalLen {
		finalString = title + " - " + channel
		// Remove a space character if it woud be the first character on second line
		if properLen(finalString) > lineLen && getNChar(finalString, lineLen) == ' ' {
			finalString = toNSlice(finalString, lineLen) + fromNSlice(finalString, lineLen+1)
		}
		return finalString
	// If trimming is needed
	} else {
		var maxChannelLen int = lineLen / 3
		var targetChannelLen int = maxChannelLen
		if properLen(channel) < targetChannelLen {
			targetChannelLen = properLen(channel)
		}

		var maxTitleLen int = totalLen - targetChannelLen - 3

		// Remove space if necessary
		if properLen(title) > lineLen && getNChar(title, lineLen) == ' ' {
			title = toNSlice(title, lineLen) + fromNSlice(title, lineLen+1)
		}

		// Trim title if needed
		if properLen(title) > maxTitleLen {
			title = toNSlice(title, maxTitleLen-2) + ".."
		}

		if properLen(channel) > targetChannelLen {
			channel = toNSlice(channel, targetChannelLen-2) + ".."
		}

		finalString = title + " - " + channel
		if properLen(finalString) > lineLen && getNChar(finalString, lineLen) == ' ' {
			finalString = toNSlice(finalString, lineLen) + fromNSlice(finalString, lineLen+1)
		}
		return finalString
	}
}

// Get filename for video page save
func getTimestampFilename() string {
	return fmt.Sprintf("%s%s%d.txt", youtube.HOME_DIR, youtube.CACHE_FOLDER, time.Now().Unix())
}

// For video page view, short version of video length
func shortLength(numerical string) string {
	num, err := strconv.Atoi(numerical)
	if err != nil {
		panic(err)
	}

	var hours int = num / 3600
	var minutes int = (num - (hours * 3600)) / 60
	var seconds int = num - (hours * 3600) - (minutes * 60)

	var toReturn string = ""
	if hours > 0 {
		toReturn += strconv.Itoa(hours) + ":"
	}
	toReturn += strconv.Itoa(minutes) + ":"
	toReturn += strconv.Itoa(seconds)
	return toReturn
}

// For video page view, long version of video length
func longLength(numerical string) string {
	num, err := strconv.Atoi(numerical)
	if err != nil {
		panic(err)
	}

	var hours int = num / 3600
	var minutes int = (num - (hours * 3600)) / 60
	var seconds int = num - (hours * 3600) - (minutes * 60)

	var toReturn string = ""
	if hours > 1 {
		toReturn += strconv.Itoa(hours) + " hours "
	} else if hours == 1 {
		toReturn += strconv.Itoa(hours) + " hour "
	}
	if minutes > 1 {
		toReturn += strconv.Itoa(minutes) + " minutes "
	} else if minutes == 1 {
		toReturn += strconv.Itoa(minutes) + " minute "
	}
	if seconds > 1 {
		toReturn += strconv.Itoa(seconds) + " seconds"
	} else if seconds == 1 {
		toReturn += strconv.Itoa(seconds) + " second"
	}
	return toReturn
}

// Copies string to system clipboard, currently only supports X11, will add in Wayland support later (need to install it)
func copyToClipboard(textToCopy string) {
	var cmd *exec.Cmd
	if config.ActiveConfig.SessionType == "wayland" {
		cmd = exec.Command("wl-copy")
	} else {
		cmd = exec.Command("xsel", "--clipboard", "--input")
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	io.WriteString(stdin, textToCopy)
	stdin.Close()
	cmd.Wait()
}

func getWindowSizeAndPosition() (int, int, int, int) {
	var ppid string = getPPID()
	return getWindowGeometry(ppid)
}

// Gets the pid of the current terminal window
func getPPID() string {
	cmd := exec.Command("ps", "-Tf")
	output, _ := cmd.Output()
	lines := strings.Split(string(output), "\n")
	var ppid string = ""
	ppid = strings.Fields(lines[1])[3]

	return ppid
}

// strconv.Atoi error handled
func getIntFromString(input string) int {
	num, err := strconv.Atoi(input)
	if err != nil {
		panic(err)
	}
	return num
}

// Gets a window ID from a pid
func getWindowID(pid string) string {
	cmd := exec.Command("xdotool", "search", "--pid", pid)
	output, _ := cmd.Output()
	return string(output)
}

// Gets the current window geometry
func getWindowGeometry(pid string) (w int, h int, x int, y int) {
	var wid string = getWindowID(pid)

	cmd := exec.Command("xdotool", "getwindowgeometry", wid)
	output, _ := cmd.Output()

	split := strings.Split(string(output), "\n")

	position := strings.Fields(split[1])[1]
	positionSplit := strings.Split(position, ",")
	x = getIntFromString(positionSplit[0])
	y = getIntFromString(positionSplit[1])

	geometry := strings.Fields(split[2])[1]
	geometrySplit := strings.Split(geometry, "x")

	w = getIntFromString(geometrySplit[0])
	h = getIntFromString(geometrySplit[1])

	return w, h, x, y
}
