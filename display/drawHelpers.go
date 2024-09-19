package display

import (
	"github.com/gdamore/tcell/v2"
	"gotube/ueberzug"
	"gotube/youtube"
)

// This file is for helper functions interacting with tcell, but not directly related to the grid, and the MainContent interface

const HORIZONTAL_BAR string = "─"
const VERTICAL_BAR string = "│"
const TOP_LEFT_CORNER string = "┌"
const TOP_RIGHT_CORNER string = "┐"
const BOTTOM_LEFT_CORNER string = "└"
const BOTTOM_RIGHT_CORNER string = "┘"

const BOLD_HORIZONTAL_BAR string = "━"
const BOLD_VERTICAL_BAR string = "┃"
const BOLD_TOP_LEFT_CORNER string = "┏"
const BOLD_TOP_RIGHT_CORNER string = "┓"
const BOLD_BOTTOM_LEFT_CORNER string = "┗"
const BOLD_BOTTOM_RIGHT_CORNER string = "┛"

const REDRAW int = 0
const REDRAW_NEW_PAGE int = 1
const DO_NOTHING int = 2
const GET_MORE int = 3

const REDRAW_IMAGES bool = true
const DONT_REDRAW_IMAGES bool = false
const SHOW_CURSOR bool = true
const HIDE_CURSOR bool = false

var styles = map[string]tcell.Style{}
var termHeight int
var termWidth int

var cursorLoc int = 0
var currentSearchTerm string = ""

// Allows defining custom methods for type not in this module
type UebChan chan ueberzug.CommandInfo
type Screen struct {
	tcell.Screen
}

// Interface outlining methods all content types should have
type MainContent interface {
	redraw(bool, bool)
	calcSizing()
	recalibrate()
	removeImgs()
	handleResize(bool, bool)
	getScreen() Screen
	getCurSel() CurSelection
	setCurSel(CurSelection)
	getUebChan() UebChan
	GetVidHolder() youtube.VideoHolder
	SetVideosList([]youtube.Video)
}

// Creators(?), Getters and Setters for items common to all content types

func GetNewScreen(screen tcell.Screen) Screen {
	return Screen{screen}
}

func (curSel CurSelection) getCurSel() CurSelection {
	return curSel
}

func (curSel CurSelection) setCurSel(newA CurSelection) {
	curSel = newA
}

func (uebChan UebChan) getUebChan() UebChan {
	return uebChan
}

func (screen Screen) getScreen() Screen {
	return screen
}

// Shutdown the screen, returning terminal to a usable state
func DisplayShutdown(screen tcell.Screen) {
	screen.Fini()
}

// Initialise the screen and perform setup
func InitScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = screen.Init()
	if err != nil {
		panic(err)
	}

	setupColours()

	screen.SetStyle(styles["white"])
	screen.Clear()

	termWidth, termHeight = screen.Size()

	return screen
}

// Setup the colours (uses terminal defaults)
func setupColours() {
	// Create colours
	redColour := tcell.PaletteColor(1)
	greenColour := tcell.PaletteColor(2)
	yellowColour := tcell.PaletteColor(3)
	magentaColour := tcell.PaletteColor(4)
	pinkColour := tcell.PaletteColor(5)
	blueColour := tcell.PaletteColor(6)
	//backgroundColour := tcell.ColorWhite

	// Use the colours to make styles and add to "styles" map
	styles["white"] = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault)
	styles["yellow"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(yellowColour)
	styles["red"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(redColour)
	styles["blue"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(blueColour)
	styles["magenta"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(magentaColour)
	styles["green"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(greenColour)
	styles["pink"] = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(pinkColour)
	styles["inverse"] = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.PaletteColor(0))
	styles["inverse"] = tcell.StyleDefault.Italic(true).Background(tcell.ColorReset).Foreground(tcell.ColorReset)
}

// Function to allow adding strings to screen instead of just individual runes (chars) (minor modifications of one in tcell examples)
func drawText(screen tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		if r == '\n' {
			row++
			if row > y2 {
				break
			}
			col = x1
			continue
		}
		screen.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
