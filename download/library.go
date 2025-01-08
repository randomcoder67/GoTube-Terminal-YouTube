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
var _ = strconv.Itoa

func GetLibrary(includeHidden bool) youtube.VideoHolder {
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
	playlists := []youtube.Video{}
	
	contents := jsonA.Contents.TwoColumnBrowseResultsRenderer.Tabs[0]
	contentsB := contents.TabRenderer.Content.SectionListRenderer.Contents
	contentsA := contentsB[1].ItemSectionRenderer.Contents[0].ShelfRenderer.Content.HorizontalListRenderer.Items
	
	

	var doneChan chan int = make(chan int)
	var err error
	_ = err
	var number int = 0
	var numberOfThumbnails int = 0
	
	

	watchLater := contentsB[2].ItemSectionRenderer.Contents[0].ShelfRenderer
	
	// Don't add Watch Later if it's disabled in config
	if includeHidden || !config.ActiveConfig.HideWatchLater {
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
	}

	
	for _, x := range contentsA {

		playlistJSON := x.LockupViewModel
		if playlistJSON.Metadata.LockupMetadataViewModel.Title.Content != "" {

			// Title

			var title string = playlistJSON.Metadata.LockupMetadataViewModel.Title.Content
			var inHidden bool = false
			for _, item := range config.ActiveConfig.PlaylistsToHide {
				if item == title {
					inHidden = true
				}
			}

			if !includeHidden && inHidden {
				continue
			}
			
			// Last Updated

			var lastUpdated string = "Unknown"
			lastUpdatedPos := playlistJSON.Metadata.LockupMetadataViewModel.Metadata.ContentMetadataViewModel.MetadataRows
			for _, y := range lastUpdatedPos {
				if strings.Contains(y.MetadataParts[0].Text.Content, "Updated") {
					lastUpdated = strings.ReplaceAll(y.MetadataParts[0].Text.Content, "Updated ", "")
				}
			}

			// Num Videos
			
			var numVideos int = 0
			var vids string = playlistJSON.ContentImage.CollectionThumbnailViewModel.PrimaryThumbnail.ThumbnailViewModel.Overlays[0].ThumbnailOverlayBadgeViewModel.ThumbnailBadges[0].ThumbnailBadgeViewModel.Text
			vids = strings.ReplaceAll(vids, ",", "")
			vids = strings.ReplaceAll(vids, "videos", "")
			vids = strings.ReplaceAll(vids, "video", "")
			vids = strings.ReplaceAll(vids, " ", "")
			i, err := strconv.Atoi(vids)
			if err == nil {
				numVideos = i
			}
			
			// Channel Name

			var author string

			for _, y := range playlistJSON.Metadata.LockupMetadataViewModel.Metadata.ContentMetadataViewModel.MetadataRows {
				if y.MetadataParts[0].Text.Content != "Playlist" && y.MetadataParts[0].Text.Content != "View full playlist" && !strings.Contains(y.MetadataParts[0].Text.Content, "Updated") {
					author = y.MetadataParts[0].Text.Content
					//visibility = "Public"
					//typeA = youtube.OTHER_PLAYLIST
				}
			}
			
			// Visibility and Type
			
			var visibility string = "Unknown"
			var typeA int = youtube.OTHER_PLAYLIST

			for _, y := range playlistJSON.Metadata.LockupMetadataViewModel.Metadata.ContentMetadataViewModel.MetadataRows {
				if len(y.MetadataParts) == 2 && y.MetadataParts[1].Text.Content == "Playlist" {
					if len(y.MetadataParts[1].Text.CommandRuns) == 0 {
						visibility = y.MetadataParts[0].Text.Content
						typeA = youtube.MY_PLAYLIST
					}
				}
			}
			if typeA == youtube.OTHER_PLAYLIST {
				visibility = "Public"
			}
			
			// Playlist ID
			
			var playlistID string
			for _, y := range playlistJSON.Metadata.LockupMetadataViewModel.Metadata.ContentMetadataViewModel.MetadataRows {
				if y.MetadataParts[0].Text.Content == "View full playlist" {
					playlistID = y.MetadataParts[0].Text.CommandRuns[0].OnTap.InnertubeCommand.CommandMetadata.WebCommandMetadata.URL
					playlistID = strings.ReplaceAll(playlistID, "/playlist?list=", "")
				}
			}
			
			// Don't add Liked Videos if it's disabled in config
			if !includeHidden && config.ActiveConfig.HideLikedVideos && playlistID == "LL" {
				continue
			}

			// Thumbnail Link

			var thumbnailLink string = playlistJSON.ContentImage.CollectionThumbnailViewModel.PrimaryThumbnail.ThumbnailViewModel.Image.Sources[0].URL

			// Thumbnail File Name

			var thumbnailFile string
			if numVideos > 0 {
				thumbnailFile = youtube.HOME_DIR + ThumbnailDir + strconv.Itoa(number) + ".png"
				numberOfThumbnails++
			} else {
				thumbnailFile = youtube.HOME_DIR + youtube.DATA_FOLDER + "thumbnails/emptyPlaylist.jpg"
			}


			number++

			// Put it all together
			playlist := youtube.Video{
				Title:         title,
				LastUpdated:   lastUpdated,
				NumVideos:     numVideos,
				Channel:       author,
				Visibility:    visibility,
				Id:            playlistID,
				ThumbnailLink: thumbnailLink,
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
