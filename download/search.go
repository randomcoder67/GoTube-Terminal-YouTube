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
		playlistJSON := x.LockupViewModel
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
		} else if playlistJSON.Metadata.LockupMetadataViewModel.Title.Content != "" {
			// Title
			var title string = playlistJSON.Metadata.LockupMetadataViewModel.Title.Content
			
			// Num Videos
			var numVideos int
			for _, y := range playlistJSON.ContentImage.CollectionThumbnailViewModel.PrimaryThumbnail.ThumbnailViewModel.Overlays {
				z := y.ThumbnailOverlayBadgeViewModel.ThumbnailBadges
				if len(z) == 1 {
					zz := z[0].ThumbnailBadgeViewModel.Icon.Sources
					if len(zz) == 1 {
						if zz[0].ClientResource.ImageName == "PLAYLISTS" {
							num, err := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(z[0].ThumbnailBadgeViewModel.Text, " videos", ""), " video", ""))
							if err != nil {
								panic(err)
							}
							numVideos = num
						}
					}
				}
			}

			// Channel

			var author string = "Unknown"

			y := playlistJSON.Metadata.LockupMetadataViewModel.Metadata.ContentMetadataViewModel.MetadataRows
			if len(y) > 0 {
				z := y[0].MetadataParts
				if len(z) > 1 {
					author = z[0].Text.Content
				}
			}
			
			// Visibility

			var visibility string = "Public"
			
			// Playlist ID

			var playlistID string
			for _, y := range playlistJSON.Metadata.LockupMetadataViewModel.Metadata.ContentMetadataViewModel.MetadataRows {
				if len(y.MetadataParts) > 0 {
					if y.MetadataParts[0].Text.Content == "View full playlist" {
						var url string = y.MetadataParts[0].Text.CommandRuns[0].OnTap.InnertubeCommand.CommandMetadata.WebCommandMetadata.URL
						playlistID = strings.ReplaceAll(url, "/playlist?list=", "")
					}
				}
			}

			// Thumbanil Link and Filename

			var thumbnailLink string
			for _, y := range playlistJSON.ContentImage.CollectionThumbnailViewModel.PrimaryThumbnail.ThumbnailViewModel.Image.Sources {
				thumbnailLink = y.URL
			}

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
				LastUpdated:   "Unknown",
				NumVideos:     numVideos,
				Channel:       author,
				Visibility:    visibility,
				Id:            playlistID,
				ThumbnailLink: thumbnailLink,
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
