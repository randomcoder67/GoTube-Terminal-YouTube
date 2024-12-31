package ueberzug

import (
	"io"
	"os/exec"
	"os"
)

const PERFORM_SEARCH = 1
const GET_SUBS = 2
const GET_HISTORY = 3
const GET_LIBRARY = 4
const GET_PLAYLIST = 5
const GET_WL = 6
const PLAY_VIDEO = 7
const IMAGE_WIDTH = "24"
const IMAGE_HEIGHT = "7"

type CommandInfo struct {
	Action     string
	Identifier string
	Path       string
	X          string
	Y          string
	H          string
	W          string
}

func Print(str string) {
	cmd := exec.Command("notify-send", str)
	cmd.Run()
}

func MainUeberzug(commands chan CommandInfo) {
	ueb := exec.Command("ueberzug", "layer", "--no-cache")
	stdin, err := ueb.StdinPipe()
	if err != nil {
		panic(err)
	}
	
	ueb.Stdout = os.Stdout
	
	defer stdin.Close()

	if err = ueb.Start(); err != nil {
		panic(err)
	}

	var curCommand CommandInfo
	for {
		curCommand = <-commands
		//Print("Command Recieved: " + curCommand.Identifier + " " + curCommand.Path + " " + curCommand.Action)
		//fmt.Fprintln(os.Stderr, "recieved command: " + curCommand.Action)
		if curCommand.Action == "exit" {
			break
		} else if curCommand.Action == "add" {
			var stdinCommand string = "{\"action\":\"add\",\"identifier\":\"" + curCommand.Identifier + "\",\"path\":\"" + curCommand.Path + "\",\"x\":\"" + curCommand.X + "\",\"y\":\"" + curCommand.Y + "\",\"width\":\"" + curCommand.W + "\",\"height\":\"" + curCommand.H + "\",\"scaler\":\"fit_contain\"}\n"
			io.WriteString(stdin, stdinCommand)
		} else if curCommand.Action == "remove" {
			var stdinCommand string = "{\"action\":\"remove\",\"identifier\":\"" + curCommand.Identifier + "\"}\n"
			io.WriteString(stdin, stdinCommand)
		} else {
			panic("Error, unrecognised command, this is likely a compile time error not runtime")
		}
	}
}

func InitUeberzug() chan CommandInfo {
	var uebChan chan CommandInfo = make(chan CommandInfo, 100)
	go MainUeberzug(uebChan)
	return uebChan
}
