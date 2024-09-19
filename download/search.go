package download

import (
	"encoding/json"
	"gotube/config"
	"gotube/download/network"
	"gotube/youtube"
	"os"
	"strconv"
	"strings"
)

func GetSearch(searchTerm string) youtube.VideoHolder {
	config.LogEvent("Getting search: " + searchTerm)
	// Get JSON text from the HTML
	var fullHTML string = network.GetHTML("https://www.youtube.com/results?search_query=" + strings.ReplaceAll(searchTerm, " ", "+"), true)
	config.FileDump("SearchRaw.html", fullHTML, false)
	var jsonText string = network.ExtractJSON(fullHTML, false)
	config.FileDump("SearchRaw.json", jsonText, false)
	// Format into correct structure
	var jsonA SearchJSON
	if err := json.Unmarshal([]byte(jsonText), &jsonA); err != nil {
		panic(err)
	}

	text, _ := json.MarshalIndent(jsonA, "", "  ")
	config.FileDump("SearchProcessed.json", string(text), false)

	contents := jsonA.Contents.TwoColumnSearchResultsRenderer.PrimaryContents
	contentsB := contents.SectionListRenderer.Contents
	contentsA := contentsB[0].ItemSectionRenderer.Contents
	if 1 > 2 {
		os.Exit(1)
	}

	contentsA = append(contentsA, contentsB[1].ItemSectionRenderer.Contents...)
	videos := []youtube.Video{}

	var doneChan chan int = make(chan int)
	var err error
	_ = err
	var numberOfThumbnails int = 0
	var number int = 0
	for _, x := range contentsA {

		videoJSON := x.VideoRenderer
		playlistJSON := x.PlaylistRenderer
		if videoJSON.Title.Runs != nil {
			// Views

			var views string = ""
			var vidType string = ""
			if videoJSON.ShortViewCountText.Runs == nil {
				views = strings.Split(videoJSON.ShortViewCountText.SimpleText, " ")[0]
				vidType = "Video"
			} else {
				views = videoJSON.ShortViewCountText.Runs[0].Text
				vidType = "Livestream"
			}

			// Published Time
			var releaseDate string = "Unknown"
			if videoJSON.PublishedTimeText.SimpleText != "" {
				releaseDate = videoJSON.PublishedTimeText.SimpleText
			}

			// Length
			var length string = "Livestream"
			if videoJSON.LengthText.SimpleText != "" {
				length = videoJSON.LengthText.SimpleText
			}
			_ = views
			
			numberOfThumbnails++
			number++
			
			// Put it all together
			video := youtube.Video{
				Title:         videoJSON.Title.Runs[0].Text,
				Views:         views,
				VidType:       vidType,
				ReleaseDate:   releaseDate,
				Length:        length,
				Id:            videoJSON.VideoID,
				Channel:       videoJSON.OwnerText.Runs[0].Text,
				ChannelID:     videoJSON.OwnerText.Runs[0].NavigationEndpoint.CommandMetadata.WebCommandMetadata.URL,
				ThumbnailLink: videoJSON.Thumbnail.Thumbnails[len(videoJSON.Thumbnail.Thumbnails)-1].URL,
				ThumbnailFile: youtube.HOME_DIR + ThumbnailDir + strconv.Itoa(number) + ".png",
				DirectLink:    "",
				StartTime:     videoJSON.NavigationEndpoint.WatchEndpoint.StartTimeSeconds,
				Type:          youtube.VIDEO,
			}
			videos = append(videos, video)
			if config.ActiveConfig.Thumbnails {
				go network.DownloadThumbnail(video.ThumbnailLink, video.ThumbnailFile, false, doneChan, false)
			}
		} else if playlistJSON.Thumbnails != nil {
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
				ThumbnailLink: playlistJSON.Thumbnails[0].Thumbnails[0].URL,
				ThumbnailFile: thumbnailFile,
				Type:          youtube.OTHER_PLAYLIST,
			}
			videos = append(videos, playlist)
			if numVideos > 0 && config.ActiveConfig.Thumbnails {
				go network.DownloadThumbnail(playlist.ThumbnailLink, playlist.ThumbnailFile, false, doneChan, false)
			}
		}
	}

	videoHolder := youtube.VideoHolder{
		Videos:            videos,
		PageType:          youtube.SEARCH,
		ContinuationToken: "",
	}
	if config.ActiveConfig.Thumbnails {
		for i := 0; i < numberOfThumbnails; i++ {
			_ = <-doneChan
		}
	}

	return videoHolder
}
