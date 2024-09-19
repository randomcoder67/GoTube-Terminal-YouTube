package main

import (
	"fmt"
	"strings"
	"gotube/config"
	"gotube/display"
	"gotube/download"
	"gotube/download/network"
	"gotube/mpv"
	"gotube/ueberzug"
	"gotube/youtube"
	"os"
	"strconv"
)

var homeDir string

const PLAYLIST_URL string = "https://www.youtube.com/playlist?list="
const WATCH_LATER_URL string = "https://www.youtube.com/playlist?list=WL"

func printHelp() {
	fmt.Println("Help info")
}

func MainREPL(request int, data []string) int {
	var currentVideos youtube.VideoHolder // Current videos to display
	//display.Print(strconv.Itoa(initialState))
	var screen display.Screen = display.GetNewScreen(display.InitScreen())
	uebChan := ueberzug.InitUeberzug()
	defer display.DisplayShutdown(screen)

	var curSel display.CurSelection
	// The REPL is in charge of getting requests from the UI, and then giving data back to it (or launching a video) - So the display module only returns to main when it needs some network request to be made
	for {
		if request == youtube.GET_LIBRARY { // If the data is a bunch of playlists
			display.StartLoading(screen)
			currentVideos = download.GetLibrary()
			display.EndLoading()
			request, data, curSel = display.TUIWithVideos(screen, currentVideos, display.CurSelection{}, uebChan)
		} else if request == youtube.VIDEO_PAGE {
			var mainVideo youtube.VideoPage
			mainVideo, currentVideos = download.GetVideoPage(data[0], data[1], false)
			display.EndLoading()
			request, data, curSel = display.VideoPageTUI(screen, currentVideos, mainVideo, display.CurSelection{}, uebChan)
		} else if request == youtube.EXIT {
			display.DisplayShutdown(screen)
			if len(data) > 0 && len(data[0]) > 0 {
				fmt.Println(data[0])
			}
			return 0
		} else { // If the request is to fetch some new videos
			display.StartLoading(screen)
			switch request {
			case youtube.PERFORM_SEARCH:
				display.SetCurrentSearchTerm(data[0])
				currentVideos = download.GetSearch(data[0])
			case youtube.GET_SUBS:
				currentVideos = download.GetSubscriptions()
			case youtube.GET_HISTORY:
				currentVideos = download.GetHistory()
			case youtube.GET_PLAYLIST:
				currentVideos = download.GetPlaylist(data[0], data[1])
			case youtube.GET_WL:
				currentVideos = download.GetPlaylist("WL", "Watch later")
			case youtube.GET_LIKED:
				currentVideos = download.GetPlaylist("LL", "Liked Videos")
			//case GET_LIKED:
				//contents = download.GetLiked()
			case youtube.GET_HOME:
				currentVideos = download.GetRecommendations()
			}
			display.EndLoading()
			request, data, curSel = display.TUIWithVideos(screen, currentVideos, display.CurSelection{}, uebChan)
		}
	}

	_ = curSel
	return 0
}

func checkFolders() {
	youtube.HOME_DIR, _ = os.UserHomeDir()
	_, err := os.Stat(youtube.HOME_DIR + youtube.CACHE_FOLDER)
	if err != nil {
		err = os.Mkdir(youtube.HOME_DIR + youtube.CACHE_FOLDER, 0755)
		if err != nil {
			panic(err)
		}
		err = os.Mkdir(youtube.HOME_DIR + youtube.CACHE_FOLDER + "/log", 0755)
		if err != nil {
			panic(err)
		}
		err = os.Mkdir(youtube.HOME_DIR + youtube.CACHE_FOLDER + "/thumbnails", 0755)
		if err != nil {
			panic(err)
		}
	}

	_, err = os.Stat(youtube.HOME_DIR + youtube.CACHE_FOLDER + youtube.FRECENCY_PLAYLISTS_FILE)
	if err != nil {
		err = os.WriteFile(youtube.HOME_DIR + youtube.CACHE_FOLDER + youtube.FRECENCY_PLAYLISTS_FILE, []byte{}, 0666)
		if err != nil {
			panic(err)
		}
	}
}

func checkThumbnailFolder() {
	_, err := os.Stat(youtube.HOME_DIR + download.ThumbnailDir)
	if err != nil {
		err = os.Mkdir(youtube.HOME_DIR + download.ThumbnailDir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func deleteThumbnailFolder() {
	// Added this check, just in case ThumbnailDir was ever blank, to stop user home dir being deleted
	// With this check, the worse that could happen was ~/.cache/gotube is deleted, which isn't a big deal
	if strings.Contains(download.ThumbnailDir, ".cache/gotube") {
		os.RemoveAll(youtube.HOME_DIR + download.ThumbnailDir)
	}
}

func initFrecency() {
	display.InitRecentPlaylists(download.GetTopN(youtube.HOME_DIR + youtube.CACHE_FOLDER + youtube.FRECENCY_PLAYLISTS_FILE, 8))
}

func fork(args []string) {
	//download.Print("in fork")
	switch args[0] {
	case "--play":
		mpv.StartPlayback(args[1], args[2], args[3], args[4], args[5], args[6], args[7])
	case "--get-quality":
		mpv.GetQualityLinks(args[1], args[2])
	case "--mark-watched":
		//download.Print("1: " + args[1])
		//download.Print("2: " + args[2])
		//download.Print("3: " + args[3])
		//download.Print("4: " + args[4])
		//download.Print("going into mark watched")
		mpv.MarkWatched(args[1], args[2], args[3], args[4])
	case "--get-video-data":
		mpv.GetVideoData(args[1])
	case "--do-http":
		fmt.Println(network.GetHTML(args[1], true))
	case "--cookies":
		network.GetCookies()
	case "--testing":
		download.GetLibrary()
	case "--dump-json":
		config.OpenLogFile()
		config.InitConfig(false, false, false, "firefox")
		download.InitThumbnailDir()
		checkThumbnailFolder()
		defer config.CloseLogFile()
		defer deleteThumbnailFolder()
		html := network.GetHTML("https://www.youtube.com/@chudlogic/playlists/", true)
		json := network.ExtractJSON(html, true);
		os.WriteFile("out.json", []byte(json), 0666)
	}
}

func batchAdd(filename string, playlistID string, start int) {
	config.OpenLogFile()
	config.InitConfig(false, false, false, "firefox")
	download.InitThumbnailDir()
	checkThumbnailFolder()
	defer config.CloseLogFile()
	defer deleteThumbnailFolder()
	dat, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	
	var fileContents string = strings.TrimSuffix(string(dat), "\n")
	videoIDs := strings.Split(fileContents, "\n");

	for i:=0; i<len(videoIDs); {
		if i + 10 > len(videoIDs) {
			fmt.Println(videoIDs[i:])
			var ret bool = download.AddToPlaylistMany(videoIDs[i:], playlistID)
			fmt.Println(ret)
		} else {
			fmt.Println(videoIDs[i:i+10])
			var ret bool = download.AddToPlaylistMany(videoIDs[i:i+10], playlistID)
			fmt.Println(ret)
		}
		i += 10
	}
	/*
	for i:=start; i<len(videoIDs); i++ {
		var ret bool = download.AddToPlaylist(videoIDs[i], playlistID)
		if ret {
			fmt.Printf("Added (%04d): %s\n", i, videoIDs[i])
		} else {
			fmt.Printf("Error (%04d): %s\n", i, videoIDs[i])
			return
		}
	}
	*/
}

func main() {
	checkFolders()
	initFrecency()

	//download.Print("START")
	// Only used by program internally to fork
	if len(os.Args) > 1 && os.Args[1] == "--fork" {
		fork(os.Args[2:])
		return
	} else if len(os.Args) == 5 && os.Args[1] == "--add" {
		start, _ := strconv.Atoi(os.Args[4])
		batchAdd(os.Args[2], os.Args[3], start)
		return
	}

	var initialState int = 0
	var initialData []string = []string{"", ""}
	var logEvents bool = false
	var dumpJSON bool = false
	var thumbnails bool = true
	var browserCookies string = "firefox"

	for i := 1; i < len(os.Args); {
		//fmt.Println(os.Args[i])
		switch os.Args[i] {
		case "-h", "--help", "help":
			printHelp()
			os.Exit(0)
		case "-lik", "--liked-videos":
			i++
			if initialState == 0 {
				initialState = youtube.GET_PLAYLIST
				initialData[0] = "LL"
				initialData[1] = "Liked Videos"
			}
		case "-s", "--subscriptions":
			i++
			if initialState == 0 {
				initialState = youtube.GET_SUBS
			}
		case "-hs", "--history":
			i++
			if initialState == 0 {
				initialState = youtube.GET_HISTORY
				continue
			}
		case "-wl", "--watchlater":
			i++
			if initialState == 0 {
				initialState = youtube.GET_PLAYLIST
				initialData[0] = "WL"
				initialData[1] = "Watch later"
				continue
			}
		case "-p", "-l", "--playlists", "--library":
			i++
			if initialState == 0 {
				initialState = youtube.GET_LIBRARY
				if len(os.Args) > i+1 {
					initialState = youtube.GET_PLAYLIST
					initialData[0] = os.Args[i]
					initialData[1] = os.Args[i+1]
					i++
					i++
				}
				continue
			}
		case "--search":
			i++
			if initialState == 0 {
				initialState = youtube.PERFORM_SEARCH
				if len(os.Args) > i {
					initialData[0] = os.Args[i]
					i++
				}
				continue
			}
		case "-hm", "--recommendations", "--home":
			i++
			if initialState == 0 {
				initialState = youtube.GET_HOME
				continue
			}
		case "--log":
			logEvents = true
			i++
		case "--dump-json":
			dumpJSON = true
			i++
		case "--no-thumbnails":
			thumbnails = false
			i++
		case "--chromium":
			browserCookies = "chromium"
			i++
		default:
			fmt.Println("Error, incorrect arguments")
			printHelp()
			os.Exit(1)
		}
	}
	config.OpenLogFile()
	config.InitConfig(logEvents, dumpJSON, thumbnails, browserCookies)
	download.InitThumbnailDir()
	checkThumbnailFolder()
	defer config.CloseLogFile()
	defer deleteThumbnailFolder()
	//display.TUI()
	// This is all testing stuff, will be gone when program is ready
	if len(os.Args) > 1 && os.Args[1] == "-play" {
		download.GetLibrary()
	} else if len(os.Args) > 1 && os.Args[1] == "-ha" {
		download.GetHistory()
	} else if len(os.Args) > 1 && os.Args[1] == "-d" {
		download.GetDirectLinks(os.Args[2])
	} else if len(os.Args) > 1 && os.Args[1] == "-test" {
		fmt.Println(download.GetAddToPlaylist(os.Args[2]))
	} else {
		/*ls
		if len(os.Args) == 3 {
			download.GetDirectLink(os.Args[2])
		} else if os.Args[1] == "-s" {
			download.GetSubscriptions()
		} else if os.Args[1] == "-w" {
		} else {
			//download.GetHTMLTesting(os.Args[1])
		}
		*/
		MainREPL(initialState, initialData)
	}

}
