package display

import (
	"github.com/gdamore/tcell/v2"
	"strings"
	"time"
)

// This file contains functions related to the status bar

// Channel to start/stop the loading animation
var cmdChan chan int

// The characters making up the loading spinner
var loadingChars = [4]string{"-", "\\", "|", "/"}

// Draw the status bar with given data (does not sync screen)
func drawStatusBar(screen tcell.Screen, data []string) {
	var dataString string = ""
	for _, text := range data {
		dataString += text + " - "
	}
	dataString = dataString[:len(dataString)-3]
	drawText(screen, 1, termHeight-1, termWidth-1, termHeight-1, styles["inverse"], dataString + strings.Repeat(" ", termWidth))
}

// Draw the loading animation
func drawLoading(screen tcell.Screen, cmdChan chan int) {
	var i int = 0
	for {
		select {
		case <-cmdChan:
			drawStatusBar(screen, []string{"            "})
			cmdChan <- 1
			return
		default:
			drawText(screen, 1, termHeight-1, termWidth-1, termHeight-1, styles["inverse"], "Loading " + loadingChars[i%4] + strings.Repeat(" ", termWidth))
			screen.Show()
			i++
			time.Sleep(time.Millisecond * 200)
		}
	}
}

// Call to start loading animation
func StartLoading(screen tcell.Screen) {
	cmdChan = make(chan int)
	go drawLoading(screen, cmdChan)
}

// Call to stop loading animation
func EndLoading() {
	cmdChan <- 1
	_ = <-cmdChan
}
