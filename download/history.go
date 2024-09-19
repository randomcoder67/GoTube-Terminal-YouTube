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

func GetHistory() youtube.VideoHolder {
	config.LogEvent("Getting history")
	// Get JSON text from the HTML
	var fullHTML string = network.GetHTML("https://www.youtube.com/feed/history", true)
	config.FileDump("HistoryRaw.html", fullHTML, false)
	var jsonText string = network.ExtractJSON(fullHTML, false)
	config.FileDump("HistoryRaw.json", jsonText, false)
	// Format into correct structure
	var jsonA HistJSON
	if err := json.Unmarshal([]byte(jsonText), &jsonA); err != nil {
		panic(err)
	}
	text, _ := json.MarshalIndent(jsonA, "", "  ")
	config.FileDump("HistoryProcessed.json", string(text), false)

	contents := jsonA.Contents.TwoColumnBrowseResultsRenderer.Tabs[0]
	contentsB := contents.TabRenderer.Content.SectionListRenderer.Contents
	contentsA := contentsB[0].ItemSectionRenderer.Contents
	if 1 > 2 {
		os.Exit(1)
	}

	contentsA = append(contentsA, contentsB[1].ItemSectionRenderer.Contents...)
	videos := []youtube.Video{}

	var doneChan chan int = make(chan int)
	var err error
	_ = err
	var number int = 0
	for _, x := range contentsA {

		videoJSON := x.VideoRenderer
		if videoJSON.Title.Runs != nil {

			// Views
			var views string
			var vidType string
			if videoJSON.ShortViewCountText.Runs == nil {
				views = strings.Split(videoJSON.ShortViewCountText.SimpleText, " ")[0]
				vidType = "Video"
			} else {
				views = videoJSON.ShortViewCountText.Runs[0].Text
				vidType = "Livestream"
			}

			// Published Time (history doesn't contain release date for some reason)
			var releaseDate string = "Unknown"

			// Length
			var length string = "Livestream"
			if videoJSON.LengthText.SimpleText != "" {
				length = videoJSON.LengthText.SimpleText
			}

			number++
			_ = views
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
				ThumbnailLink: videoJSON.Thumbnail.Thumbnails[3].URL,
				ThumbnailFile: youtube.HOME_DIR + ThumbnailDir + strconv.Itoa(number) + ".png",
				DirectLink:    "",
				StartTime:     videoJSON.NavigationEndpoint.WatchEndpoint.StartTimeSeconds,
				Type:          youtube.VIDEO,
			}
			videos = append(videos, video)
			//fmt.Println(video.ThumbnailLink)
			if config.ActiveConfig.Thumbnails {
				go network.DownloadThumbnail(video.ThumbnailLink, video.ThumbnailFile, false, doneChan, false)
			}
		}
	}

	videoHolder := youtube.VideoHolder{
		Videos:            videos,
		PageType:          youtube.HISTORY,
		ContinuationToken: "",
		//ContinuationToken: contentsB[len(contentsB)-1].ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token,
	}
	//Print(contentsB[len(contentsB)-1].ContinuationItemRenderer.ContinuationEndpoint.ContinuationCommand.Token)

	//fmt.Println("DONE Data")
	if config.ActiveConfig.Thumbnails {
		for i := 0; i < number; i++ {
			//fmt.Println("Doing thumbnails")
			_ = <-doneChan
		}
	}
	return videoHolder
}
