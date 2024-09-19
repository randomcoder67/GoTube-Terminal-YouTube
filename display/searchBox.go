package display

import (
	"github.com/gdamore/tcell/v2"
	"gotube/youtube"
	"strconv"
	"strings"
)

// This file handles the drawing and management of the search box, including parsing commands (e.g. /subs, /wl)

// Ctrl Backspace
func backspaceWord(s tcell.Screen) {
	var lastSpace int = 0
	if cursorLoc > 0 && cursorLoc < len(currentSearchTerm)+1 {
		var initialLen int = len(currentSearchTerm)
		var relevantString string = currentSearchTerm[:cursorLoc]
		lastSpace = strings.LastIndex(relevantString, " ")
		if currentSearchTerm[cursorLoc-1] == ' ' {
			lastSpace = strings.LastIndex(relevantString[:len(relevantString)-1], " ")
		}
		if lastSpace == -1 {
			currentSearchTerm = ""
			cursorLoc = 0
		} else {
			relevantString = relevantString[:lastSpace+1]
			currentSearchTerm = relevantString + currentSearchTerm[cursorLoc:]
			var finalLen int = len(currentSearchTerm)
			cursorLoc = cursorLoc - (initialLen - finalLen)
		}
	}
}

// Ctrl Delete - Not implemented yet
func deleteWord(s tcell.Screen) {

}

// Normal backspace
func backspaceChar(s tcell.Screen) {
	if cursorLoc > 0 {
		currentSearchTerm = currentSearchTerm[:cursorLoc-1] + currentSearchTerm[cursorLoc:]
		cursorLoc--
	}
}

// Normal delete
func deleteChar(s tcell.Screen) {
	if cursorLoc < len(currentSearchTerm) {
		currentSearchTerm = currentSearchTerm[:cursorLoc] + currentSearchTerm[cursorLoc+1:]
	}
}

// Insert char at cursor position
func insertChar(s tcell.Screen, char rune) {
	currentSearchTerm = currentSearchTerm[:cursorLoc] + string(char) + currentSearchTerm[cursorLoc:]
	cursorLoc++
}

// Renders the search box, either with or without cursor
func renderSearchBox(s tcell.Screen, showCursor bool) {
	//fmt.Println("THING" + strconv.Itoa(spareX) + "  ");
	var spareX int = curPageInfo.SpareX
	// Top line
	drawText(s, 1, 0, termWidth, 0, styles["white"], BOLD_TOP_LEFT_CORNER + strings.Repeat(BOLD_HORIZONTAL_BAR, 4) + "┯" + strings.Repeat(BOLD_HORIZONTAL_BAR, termWidth-6-spareX) + "┳" + strings.Repeat(BOLD_HORIZONTAL_BAR, spareX-4) + BOLD_TOP_RIGHT_CORNER)
	// Logo
	drawText(s, 1, 1, 7, 1, styles["white"], BOLD_VERTICAL_BAR + "   " + VERTICAL_BAR)
	// Search logo
	drawText(s, termWidth-spareX+1, 1, termWidth-2, 1, styles["white"], BOLD_VERTICAL_BAR + "   Search (/)")
	drawText(s, termWidth-2, 1, termWidth-2, 1, styles["white"], BOLD_VERTICAL_BAR)
	// Bottom line
	drawText(s, 1, 2, termWidth, 2, styles["white"], BOLD_BOTTOM_LEFT_CORNER + strings.Repeat(BOLD_HORIZONTAL_BAR, 4) + "┷" + strings.Repeat(BOLD_HORIZONTAL_BAR, termWidth-6-spareX) + "┻" + strings.Repeat(BOLD_HORIZONTAL_BAR, spareX-4) + BOLD_BOTTOM_RIGHT_CORNER)
	// Search term
	drawText(s, 8, 1, termWidth-spareX, 1, styles["white"], currentSearchTerm + strings.Repeat(" ", termWidth))
	// Cursor
	if showCursor {
		s.ShowCursor(8 + cursorLoc, 1)
	} else {
		s.HideCursor()
	}
}

// The TUI/event-loop function for the search box
func FocusSearchBox(content MainContent, lock bool, returnString bool) (int, []string) {
	screen := content.getScreen()
	renderSearchBox(screen, true)
	for {
		screen.Sync()
		termWidth, termHeight = screen.Size()
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			content.handleResize(REDRAW_IMAGES, SHOW_CURSOR)
		case *tcell.EventKey:
			// Exit
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return youtube.NONE, []string{""}
			// Enter chracter to search box
			} else if ev.Key() == tcell.KeyRune {
				insertChar(screen, ev.Rune())
				renderSearchBox(screen, true)
			// Backspace
			} else if ev.Key() == tcell.KeyBackspace2 {
				backspaceChar(screen)
				renderSearchBox(screen, true)
			// Delete
			} else if ev.Key() == tcell.KeyDelete {
				deleteChar(screen)
				renderSearchBox(screen, true)
			// Ctrl + Backspace
			} else if ev.Key() == tcell.KeyBackspace {
				backspaceWord(screen)
				renderSearchBox(screen, true)
			// Cursor Left
			} else if ev.Key() == tcell.KeyLeft {
				if cursorLoc > 0 {
					cursorLoc--
				}
				renderSearchBox(screen, true)
			// Cursor Right
			} else if ev.Key() == tcell.KeyRight {
				if cursorLoc < len(currentSearchTerm) {
					cursorLoc++
				}
				renderSearchBox(screen, true)
			// Enter
			} else if ev.Key() == tcell.KeyEnter && len(currentSearchTerm) > 0 {
				if currentSearchTerm[0] == '/' && !returnString {
					ret, data := parseCommand(currentSearchTerm)
					if ret == youtube.ERROR {
						drawStatusBar(screen, []string{"Error, invalid command"})
						continue
					}
					screen.HideCursor()
					return ret, data
				}
				screen.HideCursor()
				return youtube.PERFORM_SEARCH, []string{currentSearchTerm}
			// Tab
			} else if ev.Key() == tcell.KeyTab {
				if lock {
					drawStatusBar(screen, []string{"Error, no content displayed, enter a search term or type /q to quit"})
					continue
				}
				screen.HideCursor()
				screen.Sync()
				return youtube.NONE, []string{""}
			}
		}
	}
}

// Parse a given command, returns youtube.ERROR if not valid
func parseCommand(command string) (int, []string) {
	thing := []string{}
	switch command {
	case "/home", "/HOME":
		return youtube.GET_HOME, thing
	case "/subs", "/SUBS":
		return youtube.GET_SUBS, thing
	case "/wl", "/WL":
		return youtube.GET_WL, thing
	case "/his", "/HIS":
		return youtube.GET_HISTORY, thing
	case "/lik", "/LIK":
		return youtube.GET_LIKED, thing
	case "/p", "/P":
		return youtube.GET_LIBRARY, thing
	//case "/help":
		//return SHOW_HELP, ""
	case "/q", "/Q", "/quit", "/QUIT", "/exit", "/EXIT":
		return youtube.EXIT, thing
	default:
		if len(command) > 1 && len(command) < 3 {
			num, err := strconv.Atoi(string(command[1]))
			if err == nil && num > 0 && num < 9 {
				return youtube.GET_PLAYLIST, []string{RECENT_PLAYLIST_IDS[num], RECENT_PLAYLIST_NAMES[num]}
			}
		}
	}
	return youtube.ERROR, thing
}
