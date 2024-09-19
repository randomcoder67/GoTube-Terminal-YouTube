package display

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"gotube/download"
	"gotube/config"
	"gotube/mpv"
	"gotube/youtube"
	"os"
	"os/exec"
	"strconv"
	//"bytes"
)

// This file contains meta functions related to the various on-keypress actions that can occur in the grid. Split into different categories as not all screens will have playlists for example

// Returns the currently selected video/playlist given the current selection
func getCurSelVid(content MainContent) youtube.Video {
	return content.GetVidHolder().Videos[content.getCurSel().Index]
}

// General functions not specific to either videos or playlists, and always avalible
func handleGeneralFunctions(key tcell.Key, r rune, mod tcell.ModMask, content MainContent) (int, []string) {
	// Exit
	if key == tcell.KeyEscape || key == tcell.KeyCtrlW || r == 'q' || r == 'Q' {
		return youtube.EXIT, []string{""}
	// Exit with message (useful for scripts)
	} else if key == tcell.KeyCtrlC {
		return youtube.EXIT, []string{"exit"}
	// Copy linx (share)
	} else if r == 's' && mod == tcell.ModAlt && len(content.GetVidHolder().Videos) > 0 {
		copyLink(content.getScreen(), getCurSelVid(content).Id, getCurSelVid(content).StartTime, getCurSelVid(content).Type, true)
	} else if r == 's' && mod == tcell.ModNone && len(content.GetVidHolder().Videos) > 0 {
		copyLink(content.getScreen(), getCurSelVid(content).Id, getCurSelVid(content).StartTime, getCurSelVid(content).Type, false)
	} else if mod == tcell.ModAlt && r == 'f' && len(content.GetVidHolder().Videos) > 0 {
		openInBrowser(getCurSelVid(content).Id, getCurSelVid(content).Type)
	// Go to channel
	} else if r == 'c' && len(content.GetVidHolder().Videos) > 0 {
		Print("Go to channel")
	// Download
	} else if r == 'd' && len(content.GetVidHolder().Videos) > 0 {
		Print("Download")
	// Switch focus to search box
	} else if key == tcell.KeyTab {
		ret, data := FocusSearchBox(content, false, false)
		if ret != youtube.NONE {
			return ret, data
		}
	} else if r == '/' {
		currentSearchTerm = "/"
		cursorLoc = 1
		ret, data := FocusSearchBox(content, false, false)
		if ret != youtube.NONE {
			return ret, data
		}
	}
	return youtube.NONE, nil
}

// Functions that only apply to videos
func handleVideoFunctions(key tcell.Key, r rune, mod tcell.ModMask, content MainContent) (int, []string) {
	//Print("Handling video functions")
	// Launch video in foreground
	if key == tcell.KeyEnter && mod == tcell.ModAlt && getCurSelVid(content).Type == youtube.VIDEO {
		var timestamp string = playVideoForeground(content, false)
		return youtube.VIDEO_PAGE, []string{getCurSelVid(content).Id, timestamp}
	// Launch video in the background
	} else if key == tcell.KeyEnter && mod == tcell.ModNone && getCurSelVid(content).Type == youtube.VIDEO {
		playVideoBackground(content, false)
	// Launch video with quality options in foreground
	} else if r == '#' && mod == tcell.ModAlt && getCurSelVid(content).Type == youtube.VIDEO {
		var timestamp string = playVideoForeground(content, true)
		return youtube.VIDEO_PAGE, []string{getCurSelVid(content).Id, timestamp}
	// Launch video with quality options in background
	} else if r == '#' && mod == tcell.ModNone && getCurSelVid(content).Type == youtube.VIDEO {
		playVideoBackground(content, true)
	// Add to Watch Later
	} else if r == 'w' && getCurSelVid(content).Type == youtube.VIDEO {
		addToPlaylist(content.getScreen(), getCurSelVid(content).Id, "WL", "Watch later")
	// Add to playlist
	} else if r == 'a' && getCurSelVid(content).Type == youtube.VIDEO {
		addToPlaylistOptions(content)
	// Remove from playlist
	} else if r == 'r' && content.GetVidHolder().PageType == youtube.MY_PLAYLIST {
		removeFromPlaylist(content)
	}

	return youtube.NONE, nil
}

// Functions that only apply to playlists
func handlePlaylistFunctions(key tcell.Key, r rune, mod tcell.ModMask, content MainContent) (int, []string) {
	// Open playlist
	if key == tcell.KeyEnter && mod == tcell.ModNone && (getCurSelVid(content).Type == youtube.OTHER_PLAYLIST || getCurSelVid(content).Type == youtube.MY_PLAYLIST) {
		return youtube.GET_PLAYLIST, []string{getCurSelVid(content).Id, getCurSelVid(content).Title}
	// Open playlist in new window
	} else if key == tcell.KeyEnter && mod == tcell.ModAlt && (getCurSelVid(content).Type == youtube.OTHER_PLAYLIST || getCurSelVid(content).Type == youtube.MY_PLAYLIST) {
		openPlaylistInNewWindow(content)
	// Save to library
	} else if r == 'a' && getCurSelVid(content).Type == youtube.OTHER_PLAYLIST {
		addToLibrary(content.getScreen(), getCurSelVid(content).Id, getCurSelVid(content).Title)
	// Remove from library
	} else if r == 'r' && content.GetVidHolder().PageType == youtube.LIBRARY && getCurSelVid(content).Type == youtube.OTHER_PLAYLIST {
		removeFromLibrary(content)
	} else if r == 'r' && content.GetVidHolder().PageType == youtube.LIBRARY && getCurSelVid(content).Type == youtube.MY_PLAYLIST {
		deletePlaylist(content)
	} else if r == 'n' && content.GetVidHolder().PageType == youtube.LIBRARY {
		createPlaylist(content)
	}
	return youtube.NONE, nil
}

// Below are the individual functions which handle keypress requests

func openPlaylistInNewWindow(content MainContent) {
	cmd := exec.Command("nohup", config.ActiveConfig.Term, "-e", youtube.HOME_DIR + "/.local/bin/gotube", "-p", getCurSelVid(content).Id, getCurSelVid(content).Title)
	
	// "window.title=\"GoTube\"", "-e",
	
	cmd.Start()
	
	/*
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	
	
	os.WriteFile("command.cmd", []byte(cmd.String()), 0666)
	os.WriteFile("command.out", outb.Bytes(), 0666)
	os.WriteFile("command.err", errb.Bytes(), 0666)
	*/
}

func playVideoBackground(content MainContent, qualitySelection bool) {
	timestamp := getTimestampFilename()
	StartLoading(content.getScreen())
	playVideo(content, qualitySelection, timestamp)
	EndLoading()
}

func playVideoForeground(content MainContent, qualitySelection bool) string {
	timestamp := getTimestampFilename()
	StartLoading(content.getScreen()) // End Loading called in main, it's not missing
	content.setCurSel(playVideo(content, qualitySelection, timestamp))
	return timestamp
}

func addToLibrary(screen tcell.Screen, playlistId string, playlistName string) {
	StartLoading(screen)
	ok := download.AddToLibrary(playlistId)
	EndLoading()
	if ok {
		drawStatusBar(screen, []string{"Added playlist " + playlistName + " to library"})
	} else {
		drawStatusBar(screen, []string{"Error, could not add to library"})
	}
	screen.Sync()
}

func removeFromLibrary(content MainContent) {
	screen := content.getScreen()
	StartLoading(screen)
	ok := download.RemoveFromLibrary(getCurSelVid(content).Id)
	EndLoading()
	if ok {
		var curIndex int = content.getCurSel().Index
		editedVideos := append(content.GetVidHolder().Videos[:curIndex], content.GetVidHolder().Videos[curIndex+1:]...)
		content.SetVideosList(editedVideos)
		content.calcSizing()
		content.recalibrate()
		content.redraw(REDRAW_IMAGES, HIDE_CURSOR)
	} else {
		drawStatusBar(screen, []string{"Error, could not remove from library"})
	}
	screen.Sync()
}

func createPlaylist(content MainContent) {
	screen := content.getScreen()
	// Save current search term to restore later, then blank
	var savedSearchTerm string = currentSearchTerm
	var savedCursorLoc int = cursorLoc
	currentSearchTerm = ""
	cursorLoc = 0
	
	// Get playlist name from user
	drawStatusBar(screen, []string{"Enter new playlist name"})
	screen.Sync()
	ret, data := FocusSearchBox(content, false, true)
	var playlistName string = data[0]
	// Cancel if the name is blank, or Escape pressed (restoring previous search box contents)
	if playlistName == "" || ret == youtube.EXIT {
		drawStatusBar(screen, []string{"Cancelling"})
		currentSearchTerm = savedSearchTerm
		cursorLoc = savedCursorLoc
		screen.Sync()
		return
	}
	
	// Get playlist visibility from user
	drawStatusBar(screen, []string{"Select playlist visibility"})
	screen.Sync()
	
	visibilityOptions := []string{"Private", "Unlisted", "Public"}
	var visibilityChoice string = selectionTUI(content, visibilityOptions, false)
	var visibilityEncoded int = youtube.EncodeVisibility(visibilityChoice)
	
	// Create playlist
	StartLoading(screen)
	data, ok := download.CreatePlaylist(playlistName, visibilityEncoded)
	EndLoading()
	
	// If ok, add new empty playlist with correct data, then refresh library view
	if ok {
		newPlaylist := youtube.Video {
			Title: data[1],
			Id: data[0],
			ThumbnailFile: youtube.HOME_DIR + youtube.DATA_FOLDER + "thumbnails/emptyPlaylist.jpg",
			Channel: "Unknown",
			LastUpdated: "Never",
			NumVideos: 0,
			Visibility: visibilityChoice,
			Type: youtube.MY_PLAYLIST,
		}
		
		editedVideos := append([]youtube.Video{newPlaylist}, content.GetVidHolder().Videos...)
		temp := editedVideos[0]
		editedVideos[0] = editedVideos[1]
		editedVideos[1] = temp
		content.SetVideosList(editedVideos)
		content.calcSizing()
		content.recalibrate()
		content.redraw(REDRAW_IMAGES, HIDE_CURSOR)
		drawStatusBar(screen, []string{"Created playlist " + playlistName})
	} else {
		drawStatusBar(screen, []string{"Error, failed to create playlist"})
	}
	
	// Restore previous search box contents
	currentSearchTerm = savedSearchTerm
	cursorLoc = savedCursorLoc
	screen.Sync()
}

func deletePlaylist(content MainContent) {
	screen := content.getScreen()
	// Save current search term to restore later, then blank
	var savedSearchTerm string = currentSearchTerm
	var savedCursorLoc int = cursorLoc
	currentSearchTerm = ""
	cursorLoc = 0
	
	// Prompt user to confirm, requires typing "Delete Playlist" to be absolutely sure, as deleting a playlist is final
	// (Future improvement, save playlists somewhere so it could be undone)
	drawStatusBar(screen, []string{"Deleting " + getCurSelVid(content).Title + " permanently, type \"Delete Playlist\" to confirm"})
	ret, data := FocusSearchBox(content, false, true)
	
	// Cancel if the text is not "Delete Playlist", or Escape pressed (restoring previous search box contents)
	if ret == youtube.EXIT || data[0] != "Delete Playlist" {
		drawStatusBar(screen, []string{"Cancelling"})
		currentSearchTerm = savedSearchTerm
		cursorLoc = savedCursorLoc
		screen.Sync()
		return
	}
	
	StartLoading(screen)
	ok := download.DeletePlaylist(getCurSelVid(content).Id)
	EndLoading()

	if ok {
		var curIndex int = content.getCurSel().Index
		editedVideos := append(content.GetVidHolder().Videos[:curIndex], content.GetVidHolder().Videos[curIndex+1:]...)
		content.SetVideosList(editedVideos)
		content.calcSizing()
		content.recalibrate()
		content.redraw(REDRAW_IMAGES, HIDE_CURSOR)
	} else {
		drawStatusBar(screen, []string{"Error, could not delete playlist"})
	}
	
	// Restore previous search box contents
	currentSearchTerm = savedSearchTerm
	cursorLoc = savedCursorLoc
	screen.Sync()
}

func addToPlaylistOptions(content MainContent) {
	var videoId string = getCurSelVid(content).Id

	playlistOptionsMap, playlistOptionsList := download.GetAddToPlaylist(videoId)
	chosen := selectionTUI(content, playlistOptionsList, false)
	if chosen == "" {
		return
	}
	addToPlaylist(content.getScreen(), videoId, playlistOptionsMap[chosen], chosen)
}

func addToPlaylist(screen tcell.Screen, videoId, playlistId, playlistName string) {
	StartLoading(screen)
	ok := download.AddToPlaylist(videoId, playlistId)
	EndLoading()
	if ok {
		drawStatusBar(screen, []string{"Added to " + playlistName})
	} else {
		drawStatusBar(screen, []string{"Error, could not add to playlist"})
	}
	screen.Sync()
}

func removeFromPlaylist(content MainContent) {
	screen := content.getScreen()
	StartLoading(screen)
	ok := download.RemoveFromPlaylist(getCurSelVid(content).Id, content.GetVidHolder().PlaylistID, getCurSelVid(content).PlaylistRemoveId, getCurSelVid(content).PlaylistRemoveParams)
	EndLoading()

	if ok {
		var curIndex int = content.getCurSel().Index
		editedVideos := append(content.GetVidHolder().Videos[:curIndex], content.GetVidHolder().Videos[curIndex+1:]...)
		content.SetVideosList(editedVideos)
		content.calcSizing()
		content.recalibrate()
		content.redraw(REDRAW_IMAGES, HIDE_CURSOR)
	} else {
		drawStatusBar(screen, []string{"Error, could not remove from playlist"})
	}
	screen.Sync()
}

func copyLink(screen tcell.Screen, id string, startTime int, itemType int, timestamp bool) {
	if itemType == youtube.VIDEO {
		if timestamp {
			copyToClipboard("https://www.youtube.com/watch?v=" + id + "&t=" + strconv.Itoa(startTime))
		} else {
			copyToClipboard("https://www.youtube.com/watch?v=" + id)
		}
	} else if itemType == youtube.MY_PLAYLIST || itemType == youtube.OTHER_PLAYLIST {
		copyToClipboard("https://www.youtube.com/playlist?list=" + id)
	}
	drawStatusBar(screen, []string{"Copied link to clipboard"})
	screen.Sync()
}

func getExtension(screen tcell.Screen, videosHolder youtube.VideoHolder) youtube.VideoHolder {
	StartLoading(screen)
	videosHolder = download.GetPlaylistContinuation(videosHolder, videosHolder.ContinuationToken)
	EndLoading()
	return videosHolder
}

func openInBrowser(id string, contentType int) {
	var link string
	if contentType == youtube.VIDEO {
		link = "https://www.youtube.com/watch?v=" + id
	} else if contentType == youtube.MY_PLAYLIST || contentType == youtube.OTHER_PLAYLIST {
		link = "https://www.youtube.com/playlist?list=" + id
	}
	cmd := exec.Command("nohup", config.ActiveConfig.BrowserEnv, link)
	cmd.Start()
}

//func DetachVideo(title string, channel string, startTime string, startNum string, folderName string, quality string)

func playVideo(content MainContent, qualitySelection bool, timestamp string) CurSelection {
	//var qualityOptions map[string]youtube.Format = download.GetDirectLinks(getCurSelVid(content).Id)
	qualityOptions := []string{"2160p", "1440p", "1080p", "720p", "360p"}
	mpv.WritePlaylistFile(content.GetVidHolder())

	var desiredQuality string = "720p"
	var curSel CurSelection
	if qualitySelection {
		desiredQuality = selectionTUI(content, qualityOptions, false)
		if desiredQuality == "" {
			return curSel
		}
	}
	video := getCurSelVid(content)

	var windowWidth, windowHeight, windowPosX, windowPosY int = getWindowSizeAndPosition()
	var geometryArgument string = fmt.Sprintf("%dx%d+%d+%d", windowWidth, windowHeight, windowPosX, windowPosY)

	go mpv.DetachVideo(video.Title, video.Channel, strconv.Itoa(video.StartTime), strconv.Itoa(content.getCurSel().Index), "/tmp/gotube_" + strconv.Itoa(os.Getpid()), desiredQuality, geometryArgument)
	return curSel
}
