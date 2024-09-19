package display

import (
	"github.com/gdamore/tcell/v2"
	"gotube/ueberzug"
	"gotube/youtube"
	"math"
	"strconv"
	"strings"
)

// File containing methods/functions specific to the video page content type

// VideoPage contains an extra item, the videopage struct
type VideoPage struct {
	CurSelection
	UebChan
	youtube.VideoHolder
	youtube.VideoPage
	Screen
	MainFocused bool
}

// This function is what you call to re render everything. Note: Doesn't recalculate size, call calcSizing first if the terminal has resized
func (videoPage *VideoPage) redraw(redrawImages bool, renderCursor bool) {
	videoPage.Screen.Clear()
	// Draw the suggestions grid with videos
	drawMain(videoPage.Screen, videoPage.VideoHolder, videoPage.CurSelection, videoPage.UebChan, redrawImages, renderCursor)
	// Draw the main video
	drawMainVideoInfo(videoPage.Screen, videoPage.VideoPage, videoPage.UebChan, videoPage.MainFocused)
	videoPage.Screen.Show()
}

// Called when size changes (and on init)
func (videoPage *VideoPage) calcSizing() {
	numCellsY, spareY, numCellsX, spareX := calcStandardSizing(videoPage.Screen)

	// Main video box should be a maximum of 2 cells wide, and suggestions gird a minimum of 1 wide
	var mainX int = 0
	var mainY int = numCellsY
	var suggestionsY int = numCellsY
	if numCellsX > 3 {
		mainX = 2
	} else {
		mainX = 1
	}
	var suggestionsX int = numCellsX - mainX

	gridInfo := GridInfo{
		W:         suggestionsX,
		H:         suggestionsY,
		TotalVids: len(videoPage.VideoHolder.Videos),
		NumPages:  int(math.Ceil(float64(len(videoPage.VideoHolder.Videos)) / float64(suggestionsX*suggestionsY))),
	}

	curPageInfo = videoPageInfo{
		GridInfo: gridInfo,
		MainW:    mainX,
		MainH:    mainY,
		SpareX:   spareX,
		SpareY:   spareY,
	}
}

// Called to recalibrate CurSelection, for example when a resize is performed
func (videoPage *VideoPage) recalibrate() {
	var index int = videoPage.CurSelection.Index
	if index > curPageInfo.GridInfo.TotalVids - 1 {
		index--
		videoPage.CurSelection.Index--
	}
	var page int = index / (curPageInfo.GridInfo.H * curPageInfo.GridInfo.W)
	var pageIndex int = index - page*curPageInfo.GridInfo.H*curPageInfo.GridInfo.W
	var x int = pageIndex % curPageInfo.GridInfo.W
	var y int = pageIndex / curPageInfo.GridInfo.W
	videoPage.CurSelection.Page = page
	videoPage.CurSelection.X = x
	videoPage.CurSelection.Y = y
}

// Called to remove all images currently rendered with Ueberzug (for VideoPage this includes the main video thumbnail and channel thumbnail)
func (videoPage *VideoPage) removeImgs() {
	removeGridImages(videoPage.UebChan)
	imgCmd := ueberzug.CommandInfo{
		Action:     "remove",
		Identifier: "mainImage",
	}
	videoPage.UebChan <- imgCmd
	imgCmd = ueberzug.CommandInfo{
		Action:     "remove",
		Identifier: "mainChannelThumbnail",
	}
	videoPage.UebChan <- imgCmd
}

// Meta function to simplify resize
func (videoPage *VideoPage) handleResize(redrawImages bool, renderCursor bool) {
	videoPage.removeImgs()
	videoPage.calcSizing()
	videoPage.recalibrate()
	videoPage.redraw(redrawImages, renderCursor)
}

func VideoPageTUI(screen Screen, videoHolder youtube.VideoHolder, mainVideo youtube.VideoPage, curSel CurSelection, uebChan chan ueberzug.CommandInfo) (int, []string, CurSelection) {

	curPage := &VideoPage{
		CurSelection: curSel,
		UebChan:      uebChan,
		VideoHolder:  videoHolder,
		VideoPage:    mainVideo,
		Screen:       screen,
	}

	cursorLoc = len(currentSearchTerm)

	// Initialise curSel if it doesn't exist (if it does and is at index 0, this basically does nothing)
	if curSel.Index == 0 {
		curSel = CurSelection{
			X:     0,
			Y:     0,
			Page:  0,
			Index: 0,
		}
	}

	// Perform initial render before entering REPL
	curPage.removeImgs()
	curPage.calcSizing()
	curPage.redraw(REDRAW_IMAGES, HIDE_CURSOR)
	screen.Sync()
	var ret int = 0
	var data []string

	// REPL
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			curPage.handleResize(REDRAW_IMAGES, HIDE_CURSOR)
		case *tcell.EventKey:
			// First, handle movement
			ret, curPage.CurSelection = gridHandleMovement(ev.Key(), ev.Rune(), curPage.CurSelection)
			if ret == REDRAW {
				curPage.redraw(DONT_REDRAW_IMAGES, HIDE_CURSOR)
				continue
			} else if ret == REDRAW_NEW_PAGE {
				curPage.redraw(REDRAW_IMAGES, HIDE_CURSOR)
				continue
			} else if ret == GET_MORE && curPage.VideoHolder.ContinuationToken != "" {
				curPage.VideoHolder = getExtension(curPage.Screen, curPage.VideoHolder)
				curPage.calcSizing()
				curPage.redraw(REDRAW_IMAGES, HIDE_CURSOR)
				continue
			}

			// Then general functions
			ret, data = handleGeneralFunctions(ev.Key(), ev.Rune(), ev.Modifiers(), curPage)
			if ret != youtube.NONE {
				return ret, data, curSel
			}

			// Then video functions
			ret, data = handleVideoFunctions(ev.Key(), ev.Rune(), ev.Modifiers(), curPage)
			if ret != youtube.NONE {
				return ret, data, curSel
			}
		}
	}
}

func drawMainVideoInfo(screen tcell.Screen, mainVideo youtube.VideoPage, uebChan chan ueberzug.CommandInfo, isFocused bool) {

	//var width int = 35 + (39 * (curPageInfo.MainW-1))
	//var height int = 9 + (9 * (curPageInfo.MainW-1))

	var subscribeIcon string = ""
	var subscribeString string = "Subscribe"
	var subscribeLen int = 12
	//Print(mainVideo.SubStatus)
	if mainVideo.SubStatus == "Subbed" {
		subscribeLen = 14
		subscribeIcon = ""
		subscribeString = "Unsubscribe"
	}
	var likeIcon string = ""
	var dislikeIcon string = ""
	if mainVideo.LikeStatus == "LIKE" {
		likeIcon = ""
	} else if mainVideo.LikeStatus == "DISLIKE" {
		dislikeIcon = ""
	}

	_ = subscribeLen

	// Small page
	if curPageInfo.MainW == 1 {
		width := 24
		height := 6
		drawCell(screen, 0, 0, curPageInfo.SpareY, isFocused, 1, curPageInfo.MainH)

		drawText(screen, 28, 4, 38, 4, styles["white"], "Length:")
		drawText(screen, 28, 5, 38, 5, styles["white"], shortLength(mainVideo.Length))
		drawText(screen, 28, 6, 38, 6, styles["white"], "Views:")
		drawText(screen, 28, 7, 38, 7, styles["white"], mainVideo.ViewsShort)
		drawText(screen, 28, 8, 38, 8, styles["white"], "Release:")
		drawText(screen, 28, 9, 38, 9, styles["white"], mainVideo.ReleaseDateShort)

		drawText(screen, 3, 4 + height, 38, 5 + height, styles["white"], trimTitle(mainVideo.Title, mainVideo.Channel, 1))
		drawText(screen, 3, 6 + height, 38, 6 + height, styles["white"], likeIcon + "  " + mainVideo.Likes + " " + dislikeIcon)
		drawText(screen, 15, 6 + height, 38, 6 + height, styles["white"], subscribeIcon + "  " + subscribeString)
		drawText(screen, 3, 7 + height, 38, 3 + (curPageInfo.MainH*9), styles["white"], mainVideo.Description)

		imgCmd := ueberzug.CommandInfo{
			Action:     "add",
			Identifier: "mainImage",
			Path:       mainVideo.ThumbnailFile,
			X:          "3",
			Y:          "4",
			W:          strconv.Itoa(width),
			H:          strconv.Itoa(height),
		}

		uebChan <- imgCmd
		// Large page
	} else {
		width := 49
		height := 12
		drawCell(screen, 0, 0, curPageInfo.SpareY, isFocused, 2, curPageInfo.MainH)

		drawText(screen, 53, 4, 77, 4, styles["white"], "Length:")
		drawText(screen, 53, 5, 77, 5, styles["white"], longLength(mainVideo.Length))
		drawText(screen, 53, 6, 77, 6, styles["white"], "Views:")
		drawText(screen, 53, 7, 77, 7, styles["white"], mainVideo.Views)
		drawText(screen, 53, 8, 77, 8, styles["white"], "Release:")
		drawText(screen, 53, 9, 77, 9, styles["white"], mainVideo.ReleaseDate + " - " + mainVideo.ReleaseDateShort)

		var trimmedTitle string = trimTitle(mainVideo.Title, mainVideo.Channel, 2)
		var titleOffset int = 0
		if len(trimmedTitle) > 74 {
			titleOffset = 1
		}

		drawText(screen, 3, 4 + height, 77, 4 + height + titleOffset, styles["white"], trimmedTitle)

		const SMALL_BUTTON_TOP string = "┌────────────────┬────────────┐   ┌──────────────┐"
		const SMALL_BUTTON_BOTTOM string = "└────────────────┴────────────┘   └──────────────┘"

		drawText(screen, 3, 5 + height + titleOffset, 77, 5 + height + titleOffset, styles["white"], TOP_LEFT_CORNER + strings.Repeat(HORIZONTAL_BAR, 12 + len(mainVideo.Likes)) + "┬────────────┐   ┌" + strings.Repeat(HORIZONTAL_BAR, 5 + len(subscribeString)) + TOP_RIGHT_CORNER)
		drawText(screen, 3, 6 + height + titleOffset, 77, 6 + height + titleOffset, styles["white"], "│ Like (" + mainVideo.Likes + ") " + likeIcon + "  │ Dislike " + dislikeIcon + "  │   │ " + subscribeIcon + "  " + subscribeString + " │")
		drawText(screen, 3, 7 + height + titleOffset, 77, 7 + height + titleOffset, styles["white"], BOTTOM_LEFT_CORNER + strings.Repeat(HORIZONTAL_BAR, 12 + len(mainVideo.Likes)) + "┴────────────┘   └" + strings.Repeat(HORIZONTAL_BAR, 5 + len(subscribeString)) + BOTTOM_RIGHT_CORNER)

		//drawText(screen, 3, 6 + height, 77, 6 + height, styles["white"], likeIcon + "  " + mainVideo.Likes + " " + dislikeIcon)
		//drawText(screen, 15, 6 + height, 77, 6 + height, styles["white"], subscribeIcon + "  " + subscribeString)
		drawText(screen, 3, 8 + height + titleOffset, 77, 3 + (curPageInfo.MainH*10) - 2, styles["white"], mainVideo.Description)

		//Print(strconv.Itoa(3 + (curPageInfo.MainH*9)))

		imgCmd := ueberzug.CommandInfo{
			Action:     "add",
			Identifier: "mainImage",
			Path:       mainVideo.ThumbnailFile,
			X:          "3",
			Y:          "4",
			W:          strconv.Itoa(width),
			H:          strconv.Itoa(height),
		}
		uebChan <- imgCmd
	}

	// 70 has same spacing
	// 49 has same spacing
}
