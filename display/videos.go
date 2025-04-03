package display

import (
	"github.com/gdamore/tcell/v2"
	"gotube/ueberzug"
	"gotube/youtube"
	//"strconv"
	"math"
)

// File containing methods/functions specific to the grid only content type

// GridOnly contains only the essentially for a content type
type GridOnly struct {
	CurSelection
	UebChan
	youtube.VideoHolder
	Screen
}

// This function is what you call to re render everything. Note: Doesn't recalculate size, call calcSizing first if the terminal has resized
func (gridOnly *GridOnly) redraw(redrawImages bool, renderCursor bool) {
	gridOnly.Screen.Clear()
	drawMain(gridOnly.Screen, gridOnly.VideoHolder, gridOnly.CurSelection, gridOnly.UebChan, redrawImages, renderCursor)
	gridOnly.Screen.Show()
}

// Called when size changes (and on init)
func (gridOnly *GridOnly) calcSizing() {
	numCellsY, spareY, numCellsX, spareX := calcStandardSizing(gridOnly.Screen)

	gridInfo := GridInfo{
		W:         numCellsX,
		H:         numCellsY,
		TotalVids: len(gridOnly.VideoHolder.Videos),
		NumPages:  int(math.Ceil(float64(len(gridOnly.VideoHolder.Videos)) / float64(numCellsX*numCellsY))),
	}
	curPageInfo = videoPageInfo{
		GridInfo: gridInfo,
		MainW:    0,
		MainH:    0,
		SpareX:   spareX,
		SpareY:   spareY,
	}
}

// Called to recalibrate CurSelection, for example when a resize is performed
func (gridOnly *GridOnly) recalibrate() {
	if len(gridOnly.VideoHolder.Videos) == 0 {
		return
	}
	var index int = gridOnly.CurSelection.Index
	if index > curPageInfo.GridInfo.TotalVids-1 && index != 0 {
		index--
		gridOnly.CurSelection.Index--
	}
	var page int = index / (curPageInfo.GridInfo.H * curPageInfo.GridInfo.W)
	var pageIndex int = index - page*curPageInfo.GridInfo.H*curPageInfo.GridInfo.W
	var x int = pageIndex % curPageInfo.GridInfo.W
	var y int = pageIndex / curPageInfo.GridInfo.W
	gridOnly.CurSelection.Page = page
	gridOnly.CurSelection.X = x
	gridOnly.CurSelection.Y = y
}

// Called to remove all images currently rendered with Ueberzug (no special cases for GridOnly content type, so just call removeGridImages)
func (gridOnly *GridOnly) removeImgs() {
	removeGridImages(gridOnly.UebChan)
	imgCmd := ueberzug.CommandInfo{
		Action:     "remove",
		Identifier: "mainImage",
	}
	gridOnly.UebChan <- imgCmd
	imgCmd = ueberzug.CommandInfo{
		Action:     "remove",
		Identifier: "mainChannelThumbnail",
	}
	gridOnly.UebChan <- imgCmd
}

// Meta function to simplify resize
func (gridOnly *GridOnly) handleResize(redrawImages bool, renderCursor bool) {
	gridOnly.removeImgs()
	gridOnly.calcSizing()
	if curPageInfo.GridInfo.W < 1 || curPageInfo.GridInfo.H < 1 {
		gridOnly.Screen.Clear()
		drawText(gridOnly.Screen, 1, 0, 10, 0, styles["white"], "Too Small")
		gridOnly.Screen.Sync()
		return
	}
	gridOnly.recalibrate()
	gridOnly.redraw(redrawImages, renderCursor)
}

// TUI/Event-loop function for GridOnly, this function is active whenever the grid has focus
func TUIWithVideos(screen Screen, videoHolder youtube.VideoHolder, curSel CurSelection, uebChan chan ueberzug.CommandInfo) (int, []string, CurSelection) {

	curGrid := &GridOnly{
		CurSelection: curSel,
		UebChan:      uebChan,
		VideoHolder:  videoHolder,
		Screen:       screen,
	}

	cursorLoc = len(currentSearchTerm)

	// Initialise curSel if it doesn't exist (if it does and is at index 0, this basically does nothing)
	if len(videoHolder.Videos) == 0 {
		curGrid.CurSelection = CurSelection{
			X:     -1,
			Y:     0,
			Page:  0,
			Index: 0,
		}
	}

	// Perform initial render before entering REPL
	curGrid.removeImgs()
	curGrid.calcSizing()
	if curPageInfo.GridInfo.W == 0 || curPageInfo.GridInfo.H == 0 {
		var doLoop bool = true
		for doLoop {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventResize:
				curGrid.calcSizing()
				if curPageInfo.GridInfo.W > 0 && curPageInfo.GridInfo.H > 0 {
					doLoop = false
				}
			case *tcell.EventKey:
				r := ev.Rune()
				k := ev.Key()
				if k == tcell.KeyEscape || k == tcell.KeyCtrlW || r == 'q' || r == 'Q' {
					return youtube.EXIT, []string{""}, curSel
				} else if k == tcell.KeyCtrlC {
					return youtube.EXIT, []string{"exit"}, curSel
				}
			}
		}
	}
	curGrid.redraw(REDRAW_IMAGES, HIDE_CURSOR)
	screen.Sync()
	var ret int = 0
	var data []string

	// REPL
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			curGrid.handleResize(REDRAW_IMAGES, HIDE_CURSOR)
		case *tcell.EventKey:
			// First, handle movement
			ret, curGrid.CurSelection = gridHandleMovement(ev.Key(), ev.Rune(), curGrid.CurSelection)
			if ret == REDRAW {
				curGrid.redraw(DONT_REDRAW_IMAGES, HIDE_CURSOR)
				continue
			} else if ret == REDRAW_NEW_PAGE {
				curGrid.redraw(REDRAW_IMAGES, HIDE_CURSOR)
				continue
			} else if ret == GET_MORE && curGrid.VideoHolder.ContinuationToken != "" {
				curGrid.VideoHolder = getExtension(curGrid.Screen, curGrid.VideoHolder)
				curGrid.calcSizing()
				curGrid.redraw(REDRAW_IMAGES, HIDE_CURSOR)
				continue
			}
			
			// Then general functions
			ret, data = handleGeneralFunctions(ev.Key(), ev.Rune(), ev.Modifiers(), curGrid)
			if ret != youtube.NONE {
				return ret, data, curSel
			} else if len(curGrid.VideoHolder.Videos) == 0 {
				continue
			}

			// Then either video or playlist functions depending on what is currently highlighted
			if curGrid.VideoHolder.Videos[curGrid.CurSelection.Index].Type == youtube.VIDEO {
				ret, data = handleVideoFunctions(ev.Key(), ev.Rune(), ev.Modifiers(), curGrid)
				if ret != youtube.NONE {
					return ret, data, curSel
				}
			} else {
				ret, data = handlePlaylistFunctions(ev.Key(), ev.Rune(), ev.Modifiers(), curGrid)
				if ret != youtube.NONE {
					return ret, data, curSel
				}
			}
		}
	}
}
