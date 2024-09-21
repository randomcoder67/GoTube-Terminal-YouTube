package download

import (
	"encoding/json"
	"gotube/config"
	"gotube/download/network"
	"gotube/youtube"
	"strconv"
	"strings"
)

const PLAYLIST_URL string = "https://www.youtube.com/playlist?list="
const THUMBNAIL_DIR_START string = "/.cache/gotube/thumbnails/"

var ThumbnailDir string

func GetPlaylist(playlistId string, playlistName string) youtube.VideoHolder {
	config.LogEvent("Getting playlist " + playlistName)
	// Add to frecency file
	if playlistId != "WL" && playlistId != "LL" {
		config.LogEvent("Adding playlist to frecency file")
		AddToFile(playlistId, playlistName, youtube.HOME_DIR + youtube.CACHE_FOLDER + youtube.FRECENCY_PLAYLISTS_FILE)
	}
	// Get JSON text from the HTML
	var fullHTML string = network.GetHTML(PLAYLIST_URL + playlistId, true)
	config.FileDump("PlaylistRaw.html", fullHTML, false)
	var jsonText string = network.ExtractJSON(fullHTML, true)
	config.FileDump("PlaylistRaw.json", jsonText, false)
	// Format into correct structure
	var jsonA WLJSON
	if err := json.Unmarshal([]byte(jsonText), &jsonA); err != nil {
		panic(err)
	}

	text, _ := json.MarshalIndent(jsonA, "", "  ")
	config.FileDump("PlaylistProcessed.json", string(text), false)
	if jsonA.Contents.TwoColumnBrowseResultsRenderer.Tabs == nil {
		return youtube.VideoHolder{}
	}
	contents := jsonA.Contents.TwoColumnBrowseResultsRenderer.Tabs[0]
	contentsB := contents.TabRenderer.Content.SectionListRenderer.Contents
	contentsA := contentsB[0].ItemSectionRenderer.Contents[0].PlaylistVideoListRenderer.Contents
	videos := []youtube.Video{}

	var doneChan chan int = make(chan int)
	var err error
	_ = err
	var number int = 0
	for _, x := range contentsA {

		videoJSON := x.PlaylistVideoRenderer
		if videoJSON.Title.Runs != nil {

			// Views
			var views string = "Unknown"
			var vidType string = "Unknown"

			if videoJSON.VideoInfo.Runs != nil {
				if videoJSON.VideoInfo.Runs[1].Text == " watching" {
					views = videoJSON.VideoInfo.Runs[0].Text
					vidType = "Livestream"
				} else {
					views = videoJSON.VideoInfo.Runs[0].Text
					vidType = "Video"
				}
			}

			// Published Time
			var releaseDate string = "Livestream"
			if len(videoJSON.VideoInfo.Runs) > 2 {
				releaseDate = videoJSON.VideoInfo.Runs[2].Text
			}

			// Length
			var length string = "Livestream"
			if videoJSON.LengthText.SimpleText != "" {
				length = videoJSON.LengthText.SimpleText
			}

			// Remove params
			var playlistRemoveId string = ""
			var playlistRemoveParams string = ""
			for _, entry := range videoJSON.Menu.MenuRenderer.Items {
				if len(entry.MenuServiceItemRenderer.ServiceEndpoint.PlaylistEditEndpoint.ClientActions) > 0 {
					if entry.MenuServiceItemRenderer.ServiceEndpoint.PlaylistEditEndpoint.ClientActions[0].PlaylistRemoveVideosAction.SetVideoIds[0] != "" {
						playlistRemoveId = entry.MenuServiceItemRenderer.ServiceEndpoint.PlaylistEditEndpoint.ClientActions[0].PlaylistRemoveVideosAction.SetVideoIds[0]
						playlistRemoveParams = entry.MenuServiceItemRenderer.ServiceEndpoint.PlaylistEditEndpoint.Params
					}
				}
			}

			number++
			_ = views
			// Put it all together
			video := youtube.Video{
				Title:                videoJSON.Title.Runs[0].Text,
				Views:                views,
				VidType:              vidType,
				ReleaseDate:          releaseDate,
				Length:               length,
				Id:                   videoJSON.VideoID,
				Channel:              videoJSON.ShortBylineText.Runs[0].Text,
				ChannelID:            videoJSON.ShortBylineText.Runs[0].NavigationEndpoint.CommandMetadata.WebCommandMetadata.URL,
				ThumbnailLink:        videoJSON.Thumbnail.Thumbnails[3].URL,
				ThumbnailFile:        youtube.HOME_DIR + ThumbnailDir + strconv.Itoa(number) + ".png",
				DirectLink:           "",
				StartTime:            videoJSON.NavigationEndpoint.WatchEndpoint.StartTimeSeconds,
				PlaylistRemoveId:     playlistRemoveId,
				PlaylistRemoveParams: playlistRemoveParams,
				Type:                 youtube.VIDEO,
			}
			videos = append(videos, video)
			if config.ActiveConfig.Thumbnails {
				go network.DownloadThumbnail(video.ThumbnailLink, video.ThumbnailFile, false, doneChan, false)
			}
		}
	}
	if config.ActiveConfig.Thumbnails {
		for i := 0; i < number; i++ {
			_ = <-doneChan
		}
	}

	var pageType int
	if len(videos) == 0 {
		pageType = youtube.OTHER_PLAYLIST
	} else if videos[0].PlaylistRemoveId == "" {
		pageType = youtube.OTHER_PLAYLIST
	} else {
		pageType = youtube.MY_PLAYLIST
	}
	
	var continuationToken string = ""
	if len(contentsA) > 0 {
		continuationToken = contentsA[len(contentsA)-1].ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token
	}

	videoHolder := youtube.VideoHolder{
		Videos:            videos,
		PageType:          pageType,
		PlaylistID:        playlistId,
		ContinuationToken: continuationToken,
	}

	return videoHolder
}

func GetPlaylistContinuation(videosHolder youtube.VideoHolder, continuationToken string) youtube.VideoHolder {
	config.LogEvent("Getting playlist continuation for playlist: " + videosHolder.PlaylistID)
	videos := videosHolder.Videos

	jsonString := `{
	  "context": {
		"client": {
		  "clientName": "WEB",
		  "clientVersion": "2.20231214.06.00"
		},
		"user": {
		  "lockedSafetyMode": "false"
		},
		"request": {
		  "useSsl": "true"
		}
	  },
	  "continuation": "CONTINUE"
	}`

	jsonString = strings.ReplaceAll(jsonString, "CONTINUE", continuationToken)
	status, returnedJSONString := network.PostRequestAPI(jsonString, BROWSE_URL, "https://www.youtube.com/playlist?list=" + videosHolder.PlaylistID)

	config.FileDump("PlaylistContinuationRaw.json", returnedJSONString, false)

	_ = videos
	_ = status
	// Format into correct structure
	var jsonA ContinuationJSON
	if err := json.Unmarshal([]byte(returnedJSONString), &jsonA); err != nil {
		panic(err)
	}

	text, _ := json.MarshalIndent(jsonA, "", "  ")
	config.FileDump("PlaylistContinuationProcessed.json", string(text), false)

	contents := jsonA.OnResponseReceivedActions[0].AppendContinuationItemsAction.ContinuationItems

	videosHolder.ContinuationToken = ""

	var doneChan chan int = make(chan int)
	var err error
	_ = err
	var oldNumber int = len(videos)
	var number int = len(videos)
	for _, x := range contents {

		videoJSON := x.PlaylistVideoRenderer
		continuationJSON := x.ContinuationItemRenderer

		if videoJSON.Title.Runs != nil {

			// Views
			var views string
			var vidType string
			if len(videoJSON.VideoInfo.Runs) == 0 {
				views = "0"
				vidType = "Livestream"
			} else if videoJSON.VideoInfo.Runs[1].Text == " watching" {
				views = videoJSON.VideoInfo.Runs[0].Text
				vidType = "Livestream"
			} else {
				views = videoJSON.VideoInfo.Runs[0].Text
				vidType = "Video"
			}

			// Published Time
			var releaseDate string = "Livestream"
			if len(videoJSON.VideoInfo.Runs) > 2 {
				releaseDate = videoJSON.VideoInfo.Runs[2].Text
			}

			// Length
			var length string = "Livestream"
			if videoJSON.LengthText.SimpleText != "" {
				length = videoJSON.LengthText.SimpleText
			}

			// Remove params
			var playlistRemoveId string = ""
			var playlistRemoveParams string = ""
			for _, entry := range videoJSON.Menu.MenuRenderer.Items {
				if len(entry.MenuServiceItemRenderer.ServiceEndpoint.PlaylistEditEndpoint.ClientActions) > 0 {
					if len(entry.MenuServiceItemRenderer.ServiceEndpoint.PlaylistEditEndpoint.ClientActions[0].PlaylistRemoveVideosAction.SetVideoIds) > 0 {
						playlistRemoveId = entry.MenuServiceItemRenderer.ServiceEndpoint.PlaylistEditEndpoint.ClientActions[0].PlaylistRemoveVideosAction.SetVideoIds[0]
						playlistRemoveParams = entry.MenuServiceItemRenderer.ServiceEndpoint.PlaylistEditEndpoint.Params
					}
				}
			}

			if playlistRemoveId == "" {
				//Print("ERROR, no reomve ID")
			}

			number++
			_ = views
			// Put it all together
			video := youtube.Video{
				Title:                videoJSON.Title.Runs[0].Text,
				Views:                views,
				VidType:              vidType,
				ReleaseDate:          releaseDate,
				Length:               length,
				Id:                   videoJSON.VideoID,
				Channel:              videoJSON.ShortBylineText.Runs[0].Text,
				ChannelID:            videoJSON.ShortBylineText.Runs[0].NavigationEndpoint.CommandMetadata.WebCommandMetadata.URL,
				ThumbnailLink:        videoJSON.Thumbnail.Thumbnails[3].URL,
				ThumbnailFile:        youtube.HOME_DIR + ThumbnailDir + strconv.Itoa(number) + ".png",
				DirectLink:           "",
				StartTime:            videoJSON.NavigationEndpoint.WatchEndpoint.StartTimeSeconds,
				PlaylistRemoveId:     playlistRemoveId,
				PlaylistRemoveParams: playlistRemoveParams,
				Type:                 youtube.VIDEO,
			}
			videos = append(videos, video)
			if config.ActiveConfig.Thumbnails {
				go network.DownloadThumbnail(video.ThumbnailLink, video.ThumbnailFile, false, doneChan, false)
			}
		} else if continuationJSON.ContinuationEndpoint.ContinuationCommand.Token != "" {
			videosHolder.ContinuationToken = continuationJSON.ContinuationEndpoint.ContinuationCommand.Token
		}
	}
	//fmt.Println("DONE Data")
	if config.ActiveConfig.Thumbnails {
		for i := 0; i < number-oldNumber; i++ {
			//fmt.Println("Doing thumbnails")
			_ = <-doneChan
		}
	}

	videosHolder.Videos = videos
	return videosHolder
}
