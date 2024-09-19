package display

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"gotube/youtube"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Below the various sidebar permutations are defined

var MAIN_SIDEBAR = [][]string{
	{"  Home", "(/home)"},
	{"  Subs", "(/subs)"},
	{"  Watch Later", "(/wl)"},
	{"  History", "(/his)"},
	{"  Liked Videos", "(/lik)"},
	{"  Playlists", "(/p)"},
	{"?  Help", "(/help)"},
	{"  Exit", "(/q)"},
}

var BROWSE_VIDEO = [][]string{
	{"  Play", "(Enter)"},
	{"  Quality", "(#)"},
	{"  Add WL", "(w)"},
	{"+  Add PL", "(a)"},
	{"  Download", "(d)"},
	{"  Channel", "(c)"},
	{"  Share", "(s)"},
	{"", ""},
}

var BROWSE_PLAYLIST = [][]string{
	{"  Open", "(Enter)"},
	{"+  Save PL", "(a)"},
	{"  Download", "(d)"},
	{"  Channel", "(c)"},
	{"  Share", "(s)"},
	{"", ""},
	{"", ""},
	{"", ""},
}

var LIBRARY_MY_PLAYLIST = [][]string{
	{"  Open", "(Enter)"},
	{"  Delete PL", "(r)"},
	{"  New PL", "(n)"},
	{"  Download", "(d)"},
	{"  Channel", "(c)"},
	{"  Share", "(s)"},
	{"", ""},
	{"", ""},
}

var LIBRARY_OTHER_PLAYLIST = [][]string{
	{"  Open", "(Enter)"},
	{"  Remove PL", "(r)"},
	{"  New PL", "(n)"},
	{"  Download", "(d)"},
	{"  Channel", "(c)"},
	{"  Share", "(s)"},
	{"", ""},
	{"", ""},
}

var IN_MY_PLAYLIST = [][]string{
	{"  Play", "(Enter)"},
	{"  Quality", "(#)"},
	{"  Add WL", "(w)"},
	{"+  Add PL", "(a)"},
	{"  Remov PL", "(r)"},
	{"  Download", "(d)"},
	{"  Channel", "(c)"},
	{"  Share", "(s)"},
}

var IN_OTHER_PLAYLIST = [][]string{
	{"  Play", "(Enter)"},
	{"  Quality", "(#)"},
	{"  Add WL", "(w)"},
	{"+  Add PL", "(a)"},
	{"  Download", "(d)"},
	{"  Channel", "(c)"},
	{"  Share", "(s)"},
	{"", ""},
}

var VIDEO_PAGE = [][]string{
	{"  Like", "(l)"},
	{"  Dislike", "(d)"},
	{"  Add WL", "(w)"},
	{"+  Add PL", "(a)"},
	{"  Download", "(d)"},
	{"  Channel", "(c)"},
	{"  Add Sub", "(s)"},
	{"  Unsub", "(r)"},
}

var BLANK = [][]string{
	{"", ""},
	{"", ""},
	{"", ""},
	{"", ""},
	{"", ""},
	{"", ""},
	{"", ""},
	{"", ""},
}

var RECENT_PLAYLIST_SIDEBAR [][]string
var RECENT_PLAYLIST_IDS map[int]string
var RECENT_PLAYLIST_NAMES map[int]string

var emptySidebar = [][][]string{BLANK, BLANK}
var browseVideoSidebar = [][][]string{BROWSE_VIDEO, BLANK}
var browsePlaylistSidebar = [][][]string{BROWSE_PLAYLIST, BLANK}
var myPlaylistSidebar = [][][]string{IN_MY_PLAYLIST, BLANK}
var otherPlaylistSidebar = [][][]string{IN_OTHER_PLAYLIST, BLANK}
var videoPageSidebar = [][][]string{VIDEO_PAGE, BLANK}
var libraryMyPlaylistSidebar = [][][]string{LIBRARY_MY_PLAYLIST, BLANK}
var libraryOtherPlaylistSidebar = [][][]string{LIBRARY_OTHER_PLAYLIST, BLANK}

// Draw the empty boxes for the sidebar
func drawSidebarBoxes(screen tcell.Screen, sidebarWidth, topOfBox, size int) {
	var height = 8 + 10*(size-1)
	// Draw top and bottom
	drawText(screen, termWidth-sidebarWidth-1, topOfBox, termWidth, topOfBox, styles["white"], TOP_LEFT_CORNER+strings.Repeat(HORIZONTAL_BAR, sidebarWidth-2)+TOP_RIGHT_CORNER)
	drawText(screen, termWidth-sidebarWidth-1, topOfBox+height+1, termWidth, topOfBox+height+1, styles["white"], BOTTOM_LEFT_CORNER+strings.Repeat(HORIZONTAL_BAR, sidebarWidth-2)+BOTTOM_RIGHT_CORNER)
	// Draw sides
	drawText(screen, termWidth-sidebarWidth-1, topOfBox+1, termWidth-sidebarWidth-1, topOfBox+height, styles["white"], strings.Repeat(VERTICAL_BAR, height))
	drawText(screen, termWidth-2, topOfBox+1, termWidth-2, topOfBox+height, styles["white"], strings.Repeat(VERTICAL_BAR, height))
}

// Draw a selection menu
func drawSelectionMenu(screen tcell.Screen, options []string, selection int) {
	drawSidebarBoxes(screen, curPageInfo.SpareX-2, 3, curPageInfo.GridInfo.H)

	var insideLen int = 8 + 10 * (curPageInfo.GridInfo.H-1)
	// If multiple pages needed
	if selection >= insideLen {
		// Find page (0 indexed)
		var page int = selection / insideLen
		//Print("Page:" + strconv.Itoa(page))
		//Print("InsideLen:" + strconv.Itoa(insideLen))
		//Print("Option:" + options[0])
		options = options[insideLen*page:]
		selection -= insideLen * page
	}
	

	var leftEdge int = termWidth - curPageInfo.SpareX + 2
	var rightEdge int = termWidth - 2
	var spaces string = strings.Repeat(" ", rightEdge-leftEdge)

	// Then render the stuff, highlighting the selection
	for i := 0; i < insideLen; i++ {
		// Blank out cells
		drawText(screen, leftEdge, 4+i, rightEdge, 4+i, styles["white"], spaces)
		
		// Render entry if necessary
		if i < len(options) {
			if i == selection {
				drawText(screen, leftEdge+1, 4+i, rightEdge-1, 4+i, styles["green"], options[i])
			} else {
				drawText(screen, leftEdge+1, 4+i, rightEdge-1, 4+i, styles["white"], options[i])
			}
		}
	}
	screen.Show()
}

// TUI/event-loop for the selection menu
func selectionTUI(content MainContent, options []string, sortOptions bool) string {
	if sortOptions {
		sort.Strings(options)
	}
	var selection int = 0
	var insideLen int = 8 + 10 * (curPageInfo.GridInfo.H-1)
	
	drawSelectionMenu(content.getScreen(), options, selection)
	for {
		ev := content.getScreen().PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			content.handleResize(REDRAW_IMAGES, HIDE_CURSOR)
			drawSelectionMenu(content.getScreen(), options, selection)
			content.getScreen().Show()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyUp {
				if selection > 0 {
					selection--
				} else {
					selection = len(options) - 1
				}
			} else if ev.Key() == tcell.KeyDown {
				if selection < len(options)-1 {
					selection++
				} else {
					selection = 0
				}
			} else if ev.Key() == tcell.KeyPgUp {
				if selection - insideLen >= 0 {
					selection = selection - insideLen
				} else {
					selection = 0
				}
			} else if ev.Key() == tcell.KeyPgDn {
				if selection + insideLen <= len(options) - 1 {
					selection = selection + insideLen
				} else {
					selection = len(options)-1
				}
			} else if ev.Key() == tcell.KeyHome {
				selection = 0
			} else if ev.Key() == tcell.KeyEnd {
				selection = len(options)-1
				
			} else if ev.Rune() == 'q' || ev.Key() == tcell.KeyEscape {
				content.redraw(DONT_REDRAW_IMAGES, HIDE_CURSOR)
				return ""
			} else if ev.Key() == tcell.KeyEnter {
				content.redraw(DONT_REDRAW_IMAGES, HIDE_CURSOR)
				return options[selection]
			} else if ev.Key() == tcell.KeyBackspace2 {
				content.redraw(DONT_REDRAW_IMAGES, HIDE_CURSOR)
				return ""
			}
			drawSelectionMenu(content.getScreen(), options, selection)
		}
	}
}

// Draw sidebar, deciding which to draw based on the type of the page and the type of the currently selected entry
func drawSidebar(screen tcell.Screen, pageType, entryType int) {
	var sidebars [][][]string
	if pageType == youtube.VIDEO_PAGE {
		sidebars = videoPageSidebar
	} else if pageType == youtube.MY_PLAYLIST {
		sidebars = myPlaylistSidebar
	} else if pageType == youtube.OTHER_PLAYLIST {
		sidebars = otherPlaylistSidebar
	} else if pageType == youtube.LIBRARY {
		if entryType == youtube.MY_PLAYLIST {
			sidebars = libraryMyPlaylistSidebar
		} else if entryType == youtube.OTHER_PLAYLIST {
			sidebars = libraryOtherPlaylistSidebar
		}
	} else if pageType == youtube.SEARCH {
		if entryType == youtube.OTHER_PLAYLIST {
			sidebars = browsePlaylistSidebar
		} else if entryType == youtube.VIDEO {
			sidebars = browseVideoSidebar
		}
	} else if pageType == 0 {
		sidebars = emptySidebar
	} else {
		sidebars = browseVideoSidebar
	}

	sidebars = append(sidebars, MAIN_SIDEBAR)

	var sidebarWidth int = curPageInfo.SpareX - 2
	var sidebarNumBoxes int = curPageInfo.GridInfo.H
	for i := 0; i < sidebarNumBoxes; i++ {
		var topOfBox int = 3 + 10*i
		// First draw the box
		drawSidebarBoxes(screen, sidebarWidth, topOfBox, 1)

		// Then the contents
		if i < len(sidebars) {
			for j, line := range sidebars[i] {
				var avaliableSpace int = sidebarWidth - 4
				var spaceForDesc int = avaliableSpace - properLen(line[1]) - 1 // -1 for space
				var numPadding int = 0
				var desc string = line[0]
				if utf8.RuneCountInString(desc) > spaceForDesc {
					var extra int = properLen(desc) - spaceForDesc
					desc = toNSlice(desc, properLen(desc)-extra-2) + ".."
				} else if utf8.RuneCountInString(desc) < spaceForDesc {
					numPadding = spaceForDesc - utf8.RuneCountInString(desc)
				}

				var finalString string = fmt.Sprintf("%s%s %s", desc, strings.Repeat(" ", numPadding), line[1])

				drawText(screen, termWidth-sidebarWidth+1, topOfBox+1+j, termWidth-3, topOfBox+1+j, styles["white"], finalString)
				//DrawText(screen, termWidth-sidebarWidth+1, topOfBox+1+j, termWidth-3, topOfBox+1+j, styles["white"], strconv.Itoa(spaceForDesc) + " " + strconv.Itoa(utf8.RuneCountInString(line[0])))
			}
		}
	}
}

func InitRecentPlaylists(frecencyData [][]string) {
	RECENT_PLAYLIST_SIDEBAR = [][]string{}
	RECENT_PLAYLIST_IDS = make(map[int]string)
	RECENT_PLAYLIST_NAMES = make(map[int]string)
	for i, entry := range frecencyData {
		RECENT_PLAYLIST_SIDEBAR = append(RECENT_PLAYLIST_SIDEBAR, []string{entry[1], "(/" + strconv.Itoa(i+1) + ")"})
		RECENT_PLAYLIST_IDS[i+1] = entry[0]
		RECENT_PLAYLIST_NAMES[i+1] = entry[1]
	}
	emptySidebar[1] = RECENT_PLAYLIST_SIDEBAR
	browseVideoSidebar[1] = RECENT_PLAYLIST_SIDEBAR
	browsePlaylistSidebar[1] = RECENT_PLAYLIST_SIDEBAR
	myPlaylistSidebar[1] = RECENT_PLAYLIST_SIDEBAR
	otherPlaylistSidebar[1] = RECENT_PLAYLIST_SIDEBAR
	videoPageSidebar[1] = RECENT_PLAYLIST_SIDEBAR
	libraryMyPlaylistSidebar[1] = RECENT_PLAYLIST_SIDEBAR
	libraryOtherPlaylistSidebar[1] = RECENT_PLAYLIST_SIDEBAR
}
