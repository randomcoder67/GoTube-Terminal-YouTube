package download

import (
	"encoding/json"
	"gotube/config"
	"gotube/download/network"
	"gotube/youtube"
	"strconv"
	"strings"
	"os"
	"fmt"
)

var _ = os.WriteFile

func GetLibrary() youtube.VideoHolder {
	config.LogEvent("Getting library")
	// Get JSON text from the HTML
	var jsonText string
	var count int = 0
	for {
		var fullHTML string = network.GetHTML("https://www.youtube.com/feed/library", true)
		config.FileDump("LibraryRaw.html", fullHTML, false)
		jsonText = network.ExtractJSON(fullHTML, false)
		if strings.Contains(jsonText, "\"runs\":[{\"text\":\"Playlists\"}]") {
			break
		}
		count++
		config.LogWarning(fmt.Sprintf("Retrying GetLibrary (count: %d)", count))
	}
	config.FileDump("LibraryRaw.json", jsonText, false)
	// Format into correct structure
	var jsonA LibraryJSON
	if err := json.Unmarshal([]byte(jsonText), &jsonA); err != nil {
		panic(err)
	}
	
	//os.WriteFile("output.json", []byte(jsonText), 0666)

	text, _ := json.MarshalIndent(jsonA, "", "  ")
	//os.WriteFile("processed.json", text, 0666)
	config.FileDump("LibraryProcessed.json", string(text), false)
	
	contents := jsonA.Contents.TwoColumnBrowseResultsRenderer.Tabs[0]
	contentsB := contents.TabRenderer.Content.SectionListRenderer.Contents
	contentsA := contentsB[1].ItemSectionRenderer.Contents[0].ShelfRenderer.Content.HorizontalListRenderer.Items
	playlists := []youtube.Video{}
	
	

	var doneChan chan int = make(chan int)
	var err error
	_ = err
	var number int = 0
	var numberOfThumbnails int = 0
	
	
	watchLater := contentsB[2].ItemSectionRenderer.Contents[0].ShelfRenderer
		
	if watchLater.Title.Runs[0].Text == "Watch Later" {
	
		numVideos, err := strconv.Atoi(watchLater.TitleAnnotation.SimpleText)
		if err != nil {
			panic(err)
		}
		
		var thumbnailFile string
		if numVideos > 0 {
			thumbnailFile = youtube.HOME_DIR + ThumbnailDir + "wl.png"
			numberOfThumbnails++
		} else {
			thumbnailFile = youtube.HOME_DIR + youtube.DATA_FOLDER + "thumbnails/emptyPlaylist.jpg"
		}
		
		playlist := youtube.Video{
			Title:         "Watch Later",
			LastUpdated:   "Unknown",
			NumVideos:     numVideos,
			Channel:       "Unknown",
			Visibility:    "Private",
			Id:            "WL",
			ThumbnailLink: watchLater.Content.HorizontalListRenderer.Items[0].GridVideoRenderer.Thumbnail.Thumbnails[0].URL,
			ThumbnailFile: thumbnailFile,
			Type:          youtube.OTHER_PLAYLIST,
		}
		
		if numVideos > 0 && config.ActiveConfig.Thumbnails {
			go network.DownloadThumbnail(playlist.ThumbnailLink, playlist.ThumbnailFile, false, doneChan, false)
		}
		playlists = append(playlists, playlist)
	}
	
	
	for _, x := range contentsA {

		playlistJSON := x.GridPlaylistRenderer
		if playlistJSON.Title.SimpleText != "" {

			// Last Updated
			var lastUpdated string = "Unknown"
			if playlistJSON.PublishedTimeText.SimpleText != "" {
				lastUpdated = playlistJSON.PublishedTimeText.SimpleText
				if strings.Contains(lastUpdated, "yesterday") {
					lastUpdated = "Yesterday"
				} else if strings.Contains(lastUpdated, "today") {
					lastUpdated = "Today"
				} else if strings.Contains(lastUpdated, "days ago") {

				} else if strings.Contains(lastUpdated, "months ago") {

				} else if strings.Contains(lastUpdated, "years ago") {

				}
			}

			// Num Videos
			var numVideos int = 0
			if playlistJSON.VideoCountText.Runs != nil {
				var videosString string = playlistJSON.VideoCountText.Runs[0].Text
				if videosString == "No videos" {
					numVideos = 0
				} else if videosString == "1 video" {
					numVideos = 1
				} else {
					//Print(playlistJSON.Title.SimpleText + ": " + playlistJSON.VideoCountText.Runs[0].Text)
					// Playlists with more than 999 videos will have a comma in the number (e.g. "1,120")
					var videosString string = strings.ReplaceAll(playlistJSON.VideoCountText.Runs[0].Text, ",", "")
					numVideos, err = strconv.Atoi(videosString)
				}
			}

			var visibility string = "Unknown"

			var author string = "Unknown"
			if playlistJSON.ShortBylineText.Runs[0].NavigationEndpoint.ClickTrackingParams != "" {
				author = playlistJSON.ShortBylineText.Runs[0].Text
				visibility = "Public"
			} else {
				visibility = playlistJSON.ShortBylineText.Runs[0].Text
			}
			
			number++

			var typeA int = youtube.OTHER_PLAYLIST
			if author == "Unknown" {
				typeA = youtube.MY_PLAYLIST
			}
			
			var thumbnailFile string
			if numVideos > 0 {
				thumbnailFile = youtube.HOME_DIR + ThumbnailDir + strconv.Itoa(number) + ".png"
				numberOfThumbnails++
			} else {
				thumbnailFile = youtube.HOME_DIR + youtube.DATA_FOLDER + "thumbnails/emptyPlaylist.jpg"
			}

			// Put it all together
			playlist := youtube.Video{
				Title:         playlistJSON.Title.SimpleText,
				LastUpdated:   lastUpdated,
				NumVideos:     numVideos,
				Channel:       author,
				Visibility:    visibility,
				Id:            playlistJSON.PlaylistID,
				ThumbnailLink: playlistJSON.Thumbnail.Thumbnails[0].URL,
				ThumbnailFile: thumbnailFile,
				Type:          typeA,
			}
			playlists = append(playlists, playlist)
			if numVideos > 0 && config.ActiveConfig.Thumbnails {
				go network.DownloadThumbnail(playlist.ThumbnailLink, playlist.ThumbnailFile, false, doneChan, false)
			}
		}
	}
	if config.ActiveConfig.Thumbnails {
		for i := 0; i < numberOfThumbnails; i++ {
			_ = <-doneChan
		}
	}

	holder := youtube.VideoHolder{
		Videos:            playlists,
		PageType:          youtube.LIBRARY,
		ContinuationToken: "",
	}

	return holder
}
