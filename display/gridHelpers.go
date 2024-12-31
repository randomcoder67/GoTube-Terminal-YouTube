package display

import (
	"github.com/gdamore/tcell/v2"
	"gotube/ueberzug"
	"gotube/youtube"
	"gotube/config"
	"strconv"
	"strings"
)

// This file is for helper functions specific to drawing/managing the grid

// Information about the grid
type GridInfo struct {
	H         int
	W         int
	TotalVids int
	NumPages  int
}

// Information about the whole screen, contains GridInfo
type videoPageInfo struct {
	GridInfo GridInfo
	MainH    int
	MainW    int
	SpareX   int
	SpareY   int
}

// Struct to represent what you currently have selected in the visible grid
type CurSelection struct {
	X     int
	Y     int
	Page  int
	Index int
}

// The video page info is a global variable, as it is needed basically everywhere
var curPageInfo videoPageInfo

// Draws the main parts of the screen (status bar, sidebar, search box and the grid, either partial or whole)
func drawMain(screen tcell.Screen, videosHolder youtube.VideoHolder, curSel CurSelection, uebChan chan ueberzug.CommandInfo, redrawImages bool, renderCursor bool) {
	// Add search box
	renderSearchBox(screen, renderCursor)
	// Add status bar only if there are actually videos in the grid
	if len(videosHolder.Videos) > 0 {
		drawSidebar(screen, videosHolder.PageType, videosHolder.Videos[curSel.Index].Type)
		drawStatusBar(screen, []string{videosHolder.Videos[curSel.Index].Title, videosHolder.Videos[curSel.Index].Channel})
	} else {
		drawSidebar(screen, videosHolder.PageType, 0)
	}

	// Draw the main grid with videos
	var pageOffset int = curSel.Page * (curPageInfo.GridInfo.W * curPageInfo.GridInfo.H)
	drawGridVideos(screen, videosHolder.Videos[pageOffset:], curSel, uebChan, redrawImages)
}

// Calculate the standard sizing aspects. These are the same no matter the content
func calcStandardSizing(screen tcell.Screen) (int, int, int, int) {
	termWidth, termHeight = screen.Size()

	// One cell is 10 high, search is 3 high
	var numCellsY int = (termHeight - 4) / 10
	// Minus 1 for status bar at the bottom
	var spareY int = termHeight - 3 - numCellsY*10 - 1

	// Sidebar should be 12 minimum, and the cells are 39 wide, 2 padding at the edges (one at each end)
	var numCellsX int = (termWidth - 21) / 39
	var spareX int = termWidth - numCellsX*39 - 2

	return numCellsY, spareY, numCellsX, spareX
}

// This functions draws the grid outline (boxes themselves) by calling drawCell, and then fills in the video details and thumbnails
func drawGridVideos(screen tcell.Screen, videos []youtube.Video, curSel CurSelection, uebChan chan ueberzug.CommandInfo, redrawImages bool) {
	// Get necessary data from the current grid info
	var numCellsX int = curPageInfo.GridInfo.W + curPageInfo.MainW
	var numCellsY int = curPageInfo.GridInfo.H

	i := 0
	for posY := 0; posY < numCellsY; posY++ {
		for posX := curPageInfo.MainW; posX < numCellsX; posX++ {
			//Print("curSel.X: " + strconv.Itoa(curSel.X) + " , curSel.Y: " + strconv.Itoa(curSel.Y))
			drawCell(screen, posX, posY, curPageInfo.SpareY, (curSel.X + curPageInfo.MainW == posX && curSel.Y == posY), 1, 1)
			var leftSide int = posX*39 + 1
			var top int = posY*10 + 3
			var bottom int = top + 9

			// No need to remove the old image, as ueberzug will do that automatically if they have the same identifier
			if len(videos) > i {
				if redrawImages {
					if config.ActiveConfig.Thumbnails {
						imgCmd := ueberzug.CommandInfo{
							Action:     "add",
							Identifier: "img" + strconv.Itoa(i),
							Path:       videos[i].ThumbnailFile,
							X:          strconv.Itoa(3 + posX*39),
							Y:          strconv.Itoa(4 + posY*10),
							W:          "24",
							H:          "7",
						}

						uebChan <- imgCmd
					}
				}
				var viewTitle string = "Views:    "
				var colourString string = ""
				if videos[i].VidType == "Video" {
					viewTitle = "Views:    "
					colourString = "white"
				} else if videos[i].VidType == "Livestream" {
					viewTitle = "Watching: "
					colourString = "red"
				}
				if videos[i].Type == youtube.OTHER_PLAYLIST || videos[i].Type == youtube.MY_PLAYLIST {
					// Length
					drawText(screen, leftSide+27, top+1, leftSide+37, top+2, styles["white"], "Videos:   " + strconv.Itoa(videos[i].NumVideos))
					// Views
					drawText(screen, leftSide+27, top+3, leftSide+37, top+4, styles["white"], "Updated:  " + videos[i].LastUpdated)
					// Visibility
					drawText(screen, leftSide+27, top+5, leftSide+37, top+6, styles["white"], "Status:   " + videos[i].Visibility)
				} else {
					// Length
					drawText(screen, leftSide+27, top+1, leftSide+37, top+1, styles["white"], "Length:")
					drawText(screen, leftSide+27, top+2, leftSide+37, top+2, styles[colourString], videos[i].Length)
					// Views
					drawText(screen, leftSide+27, top+3, leftSide+37, top+4, styles["white"], viewTitle + videos[i].Views)
					// Publish Time
					drawText(screen, leftSide+27, top+5, leftSide+37, top+5, styles["white"], "Release:")
					drawText(screen, leftSide+27, top+6, leftSide+37, top+6, styles[colourString], videos[i].ReleaseDate)
				}
				// Title + Channel
				drawText(screen, leftSide+2, bottom-2, leftSide+37, bottom-1, styles["white"], trimTitle(videos[i].Title, videos[i].Channel, 1))
				// However they do need to be removed if no image should be in that position (i.e. if you have a grid of 2x2 and only 3 images)
			} else {
				if config.ActiveConfig.Thumbnails {
					imgCmd := ueberzug.CommandInfo{
						Action:     "remove",
						Identifier: "img" + strconv.Itoa(i),
					}
					uebChan <- imgCmd
				}
			}
			i++
		}
	}
	//screen.Sync()
	screen.Show()
}

// Draw the outline of a cell at given position with given attributes
func drawCell(screen tcell.Screen, posX int, posY int, offsetY int, colour bool, width int, height int) {
	var leftSide int = posX*39 + 1
	var top int = posY*10 + 3
	var bottom int = top + 9 + (10 * (height - 1))
	var widthCell int = 38 + (39 * (width - 1))

	var lenY int = 8 + (10 * (height - 1))

	style := styles["white"]
	if colour {
		style = styles["red"]
	}

	// When drawing the grid, draw left and right first, including corners, then top and bottom without corners to complete the grip
	// Draw right edge
	drawText(screen, leftSide + widthCell, top, leftSide + widthCell, bottom, style, TOP_RIGHT_CORNER + strings.Repeat(VERTICAL_BAR, lenY) + BOTTOM_RIGHT_CORNER)
	// Draw left edge
	drawText(screen, leftSide, top, leftSide, bottom, style, TOP_LEFT_CORNER + strings.Repeat(VERTICAL_BAR, lenY) + BOTTOM_LEFT_CORNER)
	// Draw top
	drawText(screen, leftSide+1, top, leftSide+widthCell, top, style, strings.Repeat(HORIZONTAL_BAR, widthCell))
	// Draw bottom
	drawText(screen, leftSide+1, bottom, leftSide+widthCell, bottom, style, strings.Repeat(HORIZONTAL_BAR, widthCell))
}

// Remove all currently displayed ueberzug images
func removeGridImages(uebChan chan ueberzug.CommandInfo) {
	for i := 0; i < curPageInfo.GridInfo.W*curPageInfo.GridInfo.H; i++ {
		imgCmd := ueberzug.CommandInfo{
			Action:     "remove",
			Identifier: "img" + strconv.Itoa(i),
		}
		uebChan <- imgCmd
	}
}

// Handle movement around the grid, supports up, down, left, right, home, end, page up, page down and hjkl (vim keys (does vim have pageup etc?)
func gridHandleMovement(key tcell.Key, r rune, curSel CurSelection) (int, CurSelection) {
	switch {
	case key == tcell.KeyPgUp:
		// Possible cases: normal, first page
		if curSel.Page > 0 {
			curSel.Page--
			curSel.Index -= curPageInfo.GridInfo.W * curPageInfo.GridInfo.H
			return REDRAW_NEW_PAGE, curSel
		} else {
			return DO_NOTHING, curSel
		}
	case key == tcell.KeyPgDn:
		// Possible cases: Normal, on last page, not on last page but offset
		// Last page
		if curSel.Page == curPageInfo.GridInfo.NumPages - 1 {
			return GET_MORE, curSel
		} else {
			curSel.Page++
			curSel.Index += curPageInfo.GridInfo.W * curPageInfo.GridInfo.H
			// If going one page down lands you out of bounds
			if curSel.Index >= curPageInfo.GridInfo.TotalVids {
				curSel.Index = curPageInfo.GridInfo.TotalVids - 1
				var index int = curSel.Index
				var page int = index / (curPageInfo.GridInfo.H * curPageInfo.GridInfo.W)
				var pageIndex int = index - page * curPageInfo.GridInfo.H * curPageInfo.GridInfo.W
				curSel.X = pageIndex % curPageInfo.GridInfo.W
				curSel.Y = pageIndex / curPageInfo.GridInfo.W
			}
			return REDRAW_NEW_PAGE, curSel
		}
	case key == tcell.KeyHome:
		if curSel.X == 0 {
			return DO_NOTHING, curSel
		} else {
			curSel.Index -= curSel.X
			curSel.X = 0
			return REDRAW, curSel
		}
	case key == tcell.KeyEnd:
		if curSel.X == curPageInfo.GridInfo.W - 1 {
			return DO_NOTHING, curSel
		} else {
			curSel.Index += (curPageInfo.GridInfo.W - 1) - curSel.X
			curSel.X = curPageInfo.GridInfo.W - 1
			if curSel.Index > curPageInfo.GridInfo.TotalVids - 1 {
				curSel.X -= curSel.Index + 1 - curPageInfo.GridInfo.TotalVids
				curSel.Index = curPageInfo.GridInfo.TotalVids - 1
			}
			return REDRAW, curSel
		}
	case key == tcell.KeyRight, r == 'l':
		// Possible cases: Normal, On edge of line, at end of page, at end of results
		// End of results
		if curSel.Index == curPageInfo.GridInfo.TotalVids - 1 {
			return GET_MORE, curSel
		// End of line but not end of page
		} else if curSel.X == curPageInfo.GridInfo.W - 1 && curSel.Y < curPageInfo.GridInfo.H - 1 {
			curSel.Y++
			curSel.X = 0
			curSel.Index++
			return REDRAW, curSel
		// End of page (that isn't the last page)
		} else if (curSel.X+1) * (curSel.Y+1) == curPageInfo.GridInfo.H * curPageInfo.GridInfo.W && curSel.Page < curPageInfo.GridInfo.NumPages - 1 {
			curSel.X = 0
			curSel.Y = 0
			curSel.Page++
			curSel.Index++
			return REDRAW_NEW_PAGE, curSel
		// Normal case
		} else if curSel.Index < curPageInfo.GridInfo.TotalVids {
			curSel.X++
			curSel.Index++
			return REDRAW, curSel
		}
	case key == tcell.KeyLeft, r == 'h':
		// Possible cases: Normal, start of line, start of page, start of results
		// Start of results (do nothing)
		if curSel.X == 0 && curSel.Y == 0 && curSel.Page == 0 {
			return DO_NOTHING, curSel
		// Start of page but not start of results (go back a page)
		} else if curSel.X == 0 && curSel.Y == 0 {
			curSel.X = curPageInfo.GridInfo.W - 1
			curSel.Y = curPageInfo.GridInfo.H - 1
			curSel.Page--
			curSel.Index--
			return REDRAW_NEW_PAGE, curSel
		// Start of line but not start of page (go to end of the line above)
		} else if curSel.X == 0 && curSel.Y > 0 {
			curSel.X = curPageInfo.GridInfo.W - 1
			curSel.Y--
			curSel.Index--
			return REDRAW, curSel
		// Normal case
		} else {
			curSel.X--
			curSel.Index--
			return REDRAW, curSel
		}
	case key == tcell.KeyDown, r == 'j':
		// Possible cases: normal, last line of last page, normal offset
		// If last line of not last page
		if curSel.Y == curPageInfo.GridInfo.H - 1 && curSel.Page + 1 < curPageInfo.GridInfo.NumPages {
			curSel.Y = 0
			curSel.Page++
			curSel.Index = curSel.Index + curPageInfo.GridInfo.W
			// Page boundary offset
			if curSel.Index >= curPageInfo.GridInfo.TotalVids {
				curSel.X -= curSel.Index + 1 - curPageInfo.GridInfo.TotalVids
				curSel.Index = curPageInfo.GridInfo.TotalVids - 1
			}
			return REDRAW_NEW_PAGE, curSel
		// Normal
		} else if curSel.Index + (curPageInfo.GridInfo.W - curSel.X) < curPageInfo.GridInfo.TotalVids {
			curSel.Y++
			curSel.Index = curSel.Index + curPageInfo.GridInfo.W
			// Normal offset
			if curSel.Index >= curPageInfo.GridInfo.TotalVids {
				curSel.X -= curSel.Index + 1 - curPageInfo.GridInfo.TotalVids
				curSel.Index = curPageInfo.GridInfo.TotalVids - 1
			}
			return REDRAW, curSel
		} else {
			return GET_MORE, curSel
		}
	case key == tcell.KeyUp, r == 'k':
		// Possible cases: normal, on first line, on first line of first page
		// First line of first page (do nothing)
		if curSel.Y == 0 && curSel.Page == 0 {
			return DO_NOTHING, curSel
		// First line not of first page (go up a page)
		} else if curSel.Y == 0 {
			curSel.Y = curPageInfo.GridInfo.H - 1
			curSel.Index = curSel.Index - curPageInfo.GridInfo.W
			curSel.Page--
			return REDRAW_NEW_PAGE, curSel
		// Normal case
		} else {
			curSel.Y--
			curSel.Index = curSel.Index - curPageInfo.GridInfo.W
			return REDRAW, curSel
		}
	}
	return DO_NOTHING, curSel
}
