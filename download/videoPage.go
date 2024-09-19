package download

import (
	"encoding/json"
	"fmt"
	"gotube/config"
	"gotube/download/network"
	"gotube/youtube"
	"strconv"
	//"os"
)

func GetVideoPage(videoID string, playbackTrackingFilename string, skipThumbnails bool) (youtube.VideoPage, youtube.VideoHolder) {
	config.LogEvent("Getting video page " + videoID)
	// Get JSON text from the HTML
	var fullHTML string = network.GetHTML("https://www.youtube.com/watch?v=" + videoID, true)
	config.FileDump("VideoPageRaw.html", fullHTML, false)
	ytInitialData, initialPlayerResponse := network.ExtractJSONVideoPage(fullHTML)
	config.FileDump("VideoPageRawYTInitialData.json", ytInitialData, false)
	config.FileDump("VideoPageRawInitialPlayerResponse.json", initialPlayerResponse, false)
	
	//os.WriteFile(videoID + "data.json", []byte(ytInitialData), 0666)
	//os.WriteFile(videoID + "response.json", []byte(initialPlayerResponse), 0666)
	
	// Format into correct structure
	var initialData VidPageInitialData
	var playerResponse VidPagePlayerResp

	if err := json.Unmarshal([]byte(ytInitialData), &initialData); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(initialPlayerResponse), &playerResponse); err != nil {
		panic(err)
	}

	textInitialData, _ := json.MarshalIndent(initialData, "", "  ")
	textPlayerResponse, _ := json.MarshalIndent(playerResponse, "", "  ")
	config.FileDump("VideoPageProcessedYTInitialData.json", string(textInitialData), false)
	config.FileDump("VideoPageProcessedInitialPlayerResponse.json", string(textPlayerResponse), false)

	// First get the main video info
	primaryVideoInfo := initialData.Contents.TwoColumnWatchNextResults.Results.Results.Contents[0].VideoPrimaryInfoRenderer

	var subStatus string = "Unsubbed"
	var subParam string
	var unSubParam string
	if len(playerResponse.Annotations) > 0 {
		subInfo := playerResponse.Annotations[0].PlayerAnnotationsExpandedRenderer.FeaturedChannel.SubscribeButton.SubscribeButtonRenderer
		unSubInfo := subInfo.ServiceEndpoints[1].SignalServiceEndpoint.Actions[0].OpenPopupAction.Popup.ConfirmDialogRenderer.ConfirmButton.ButtonRenderer.ServiceEndpoint
		if subInfo.Subscribed {
			subStatus = "Subbed"
		}
		subParam = subInfo.ServiceEndpoints[0].SubscribeEndpoint.Params
		unSubParam = unSubInfo.UnsubscribeEndpoint.Params
	} else {

		if len(playerResponse.Endscreen.EndscreenRenderer.Elements) == 0 || len(playerResponse.Endscreen.EndscreenRenderer.Elements[0].EndscreenElementRenderer.HovercardButton.SubscribeButtonRenderer.ServiceEndpoints) == 0 {
			subInfo := playerResponse.PlayerConfig.WebPlayerConfig.WebPlayerActionsPorting
			subParam = subInfo.SubscribeCommand.SubscribeEndpoint.Params
			unSubParam = subInfo.UnsubscribeCommand.UnsubscribeEndpoint.Params
			// THIS CAN'T JUST BE 2, NEED TO ITERATE THROUGH AND CHECK ALL OF THEM
			if initialData.FrameworkUpdates.EntityBatchUpdate.Mutations[2].Payload.SubscriptionStateEntity.Subscribed {
				subStatus = "Subbed"
			}

		} else {
			subInfo := playerResponse.Endscreen.EndscreenRenderer.Elements[0].EndscreenElementRenderer.HovercardButton.SubscribeButtonRenderer
			if subInfo.Subscribed {
				subStatus = "Subbed"
			}
			subParam = subInfo.ServiceEndpoints[0].SubscribeEndpoint.Params
			unSubParam = subInfo.ServiceEndpoints[1].SignalServiceEndpoint.Actions[0].OpenPopupAction.Popup.ConfirmDialogRenderer.ConfirmButton.ButtonRenderer.ServiceEndpoint.UnsubscribeEndpoint.Params
		}

	}

	if subParam == "" {
		panic("Empty sub param")
	}
	if unSubParam == "" {
		panic("Empty unsub param")
	}

	addLikeInfo := primaryVideoInfo.VideoActions.MenuRenderer.TopLevelButtons[0].SegmentedLikeDislikeButtonViewModel.LikeButtonViewModel.LikeButtonViewModel.ToggleButtonViewModel.ToggleButtonViewModel.DefaultButtonViewModel.ButtonViewModel
	if addLikeInfo.IconName != "LIKE" || addLikeInfo.OnTap.SerialCommand.Commands[1].InnertubeCommand.ModalEndpoint.Modal.ModalWithTitleAndButtonRenderer.Button.ButtonRenderer.NavigationEndpoint.SignInEndpoint.NextEndpoint.CommandMetadata.WebCommandMetadata.ApiURL != "/youtubei/v1/like/like" {
		//panic("Add like misplaced")
	}

	removeLikeInfo := primaryVideoInfo.VideoActions.MenuRenderer.TopLevelButtons[0].SegmentedLikeDislikeButtonViewModel.LikeButtonViewModel.LikeButtonViewModel.ToggleButtonViewModel.ToggleButtonViewModel.ToggledButtonViewModel.ButtonViewModel
	if removeLikeInfo.IconName != "LIKE" || removeLikeInfo.OnTap.SerialCommand.Commands[1].InnertubeCommand.CommandMetadata.WebCommandMetadata.ApiURL != "/youtubei/v1/like/removelike" {
		//panic("Remove like misplaced")
	}

	addDislikeInfo := primaryVideoInfo.VideoActions.MenuRenderer.TopLevelButtons[0].SegmentedLikeDislikeButtonViewModel.DislikeButtonViewModel.DislikeButtonViewModel.ToggleButtonViewModel.ToggleButtonViewModel.DefaultButtonViewModel.ButtonViewModel
	if addDislikeInfo.IconName != "DISLIKE" || addDislikeInfo.OnTap.SerialCommand.Commands[1].InnertubeCommand.ModalEndpoint.Modal.ModalWithTitleAndButtonRenderer.Button.ButtonRenderer.NavigationEndpoint.SignInEndpoint.NextEndpoint.CommandMetadata.WebCommandMetadata.ApiURL != "/youtubei/v1/like/dislike" {
		//panic("Add dislike misplaced")
	}

	removeDislikeInfo := primaryVideoInfo.VideoActions.MenuRenderer.TopLevelButtons[0].SegmentedLikeDislikeButtonViewModel.DislikeButtonViewModel.DislikeButtonViewModel.ToggleButtonViewModel.ToggleButtonViewModel.ToggledButtonViewModel.ButtonViewModel
	if removeDislikeInfo.IconName != "DISLIKE" || removeDislikeInfo.OnTap.SerialCommand.Commands[1].InnertubeCommand.CommandMetadata.WebCommandMetadata.ApiURL != "/youtubei/v1/like/removelike" {
		//panic("Remove dislike misplaced")
	}

	mainVideo := youtube.VideoPage{
		Title:                  playerResponse.VideoDetails.Title,
		Views:                  primaryVideoInfo.ViewCount.VideoViewCountRenderer.ViewCount.SimpleText,
		ViewsShort:             primaryVideoInfo.ViewCount.VideoViewCountRenderer.ShortViewCount.SimpleText,
		VidType:                "",
		ReleaseDate:            primaryVideoInfo.DateText.SimpleText,
		ReleaseDateShort:       primaryVideoInfo.RelativeDateText.SimpleText,
		Length:                 playerResponse.VideoDetails.LengthSeconds,
		Likes:                  primaryVideoInfo.VideoActions.MenuRenderer.TopLevelButtons[0].SegmentedLikeDislikeButtonViewModel.LikeButtonViewModel.LikeButtonViewModel.ToggleButtonViewModel.ToggleButtonViewModel.DefaultButtonViewModel.ButtonViewModel.Title,
		Id:                     playerResponse.VideoDetails.VideoID,
		Channel:                playerResponse.VideoDetails.Author,
		ChannelID:              playerResponse.VideoDetails.ChannelID,
		ChannelThumbnailLink:   initialData.Contents.TwoColumnWatchNextResults.Results.Results.Contents[1].VideoSecondaryInfoRenderer.Owner.VideoOwnerRenderer.Thumbnail.Thumbnails[2].URL,
		ChannelThumbnailFile:   youtube.HOME_DIR + ThumbnailDir + "mainChannel.png",
		ThumbnailLink:          playerResponse.VideoDetails.Thumbnail.Thumbnails[len(playerResponse.VideoDetails.Thumbnail.Thumbnails)-1].URL,
		ThumbnailFile:          youtube.HOME_DIR + ThumbnailDir + "main.png",
		DirectLink:             "",
		Description:            playerResponse.VideoDetails.ShortDescription,
		SubStatus:              subStatus,
		SubParam:               subParam,
		UnSubParam:             unSubParam,
		LikeStatus:             primaryVideoInfo.VideoActions.MenuRenderer.TopLevelButtons[0].SegmentedLikeDislikeButtonViewModel.LikeButtonViewModel.LikeButtonViewModel.LikeStatusEntity.LikeStatus,
		AddLikeParam:           addLikeInfo.OnTap.SerialCommand.Commands[1].InnertubeCommand.ModalEndpoint.Modal.ModalWithTitleAndButtonRenderer.Button.ButtonRenderer.NavigationEndpoint.SignInEndpoint.NextEndpoint.LikeEndpoint.LikeParams,
		RemoveLikeParam:        removeLikeInfo.OnTap.SerialCommand.Commands[1].InnertubeCommand.LikeEndpoint.RemoveLikeParams,
		AddDislikeParam:        addDislikeInfo.OnTap.SerialCommand.Commands[1].InnertubeCommand.ModalEndpoint.Modal.ModalWithTitleAndButtonRenderer.Button.ButtonRenderer.NavigationEndpoint.SignInEndpoint.NextEndpoint.LikeEndpoint.DislikeParams,
		RemoveDislikeParam:     removeDislikeInfo.OnTap.SerialCommand.Commands[1].InnertubeCommand.LikeEndpoint.RemoveLikeParams,
		VideoStatsPlaybackURL:  playerResponse.PlaybackTracking.VideoStatsPlaybackURL.BaseURL,
		VideoStatsWatchtimeURL: playerResponse.PlaybackTracking.VideoStatsWatchtimeURL.BaseURL,
	}

	var doneChan chan int = make(chan int)

	// Download main video thumbnail and channel thumbnail
	if !skipThumbnails && config.ActiveConfig.Thumbnails {
		go network.DownloadThumbnail(mainVideo.ThumbnailLink, mainVideo.ThumbnailFile, false, doneChan, true)
		go network.DownloadThumbnail(mainVideo.ChannelThumbnailLink, mainVideo.ChannelThumbnailFile, false, doneChan, true)
		_ = <-doneChan
		_ = <-doneChan
	}

	// Then get the suggestions
	videos := []youtube.Video{}

	contents := initialData.Contents.TwoColumnWatchNextResults.SecondaryResults.SecondaryResults.Results[1].ItemSectionRenderer.Contents

	var number int = 0

	for _, entry := range contents {
		if entry.CompactVideoRenderer.VideoID != "" {

			views := entry.CompactVideoRenderer.ShortViewCountText.SimpleText
			vidType := "Video"
			if len(entry.CompactVideoRenderer.ShortViewCountText.Runs) != 0 {
				views = entry.CompactVideoRenderer.ShortViewCountText.Runs[0].Text
				vidType = "Livestream"
			}

			length := "Unknown"
			if entry.CompactVideoRenderer.LengthText.SimpleText != "" {
				length = entry.CompactVideoRenderer.LengthText.SimpleText
			}

			video := youtube.Video{
				Title:         entry.CompactVideoRenderer.Title.SimpleText,
				Views:         views,
				VidType:       vidType,
				ReleaseDate:   entry.CompactVideoRenderer.PublishedTimeText.SimpleText,
				Length:        length,
				Id:            entry.CompactVideoRenderer.VideoID,
				Channel:       entry.CompactVideoRenderer.ShortBylineText.Runs[0].Text,
				ChannelID:     entry.CompactVideoRenderer.ShortBylineText.Runs[0].NavigationEndpoint.CommandMetadata.WebCommandMetadata.URL,
				ThumbnailLink: entry.CompactVideoRenderer.Thumbnail.Thumbnails[1].URL,
				ThumbnailFile: youtube.HOME_DIR + ThumbnailDir + strconv.Itoa(number) + ".png",
				DirectLink:    "",
				StartTime:     entry.CompactVideoRenderer.NavigationEndpoint.WatchEndpoint.StartTimeSeconds,
				Type:          youtube.VIDEO,
			}
			number++
			videos = append(videos, video)

			if !skipThumbnails && config.ActiveConfig.Thumbnails {
				go network.DownloadThumbnail(video.ThumbnailLink, video.ThumbnailFile, false, doneChan, true)
			}
		}
	}

	videoHolder := youtube.VideoHolder{
		Videos:            videos,
		PageType:          youtube.VIDEO_PAGE,
		ContinuationToken: "",
	}

	//fmt.Println("HERE")
	//os.WriteFile("VideoPageProcessedYTInitialData.json", []byte(ytInitialData), 0666)
	//os.WriteFile("VideoPageProcessedInitialPlayerResponse.json", []byte(initialPlayerResponse), 0666)

	// Chapters
	var chaptersString string = ""
	if initialData.PlayerOverlays.PlayerOverlayRenderer.DecoratedPlayerBarRenderer.DecoratedPlayerBarRenderer.PlayerBar.MultiMarkersPlayerBarRenderer.MarkersMap != nil {
		chapters := initialData.PlayerOverlays.PlayerOverlayRenderer.DecoratedPlayerBarRenderer.DecoratedPlayerBarRenderer.PlayerBar.MultiMarkersPlayerBarRenderer.MarkersMap[0].Value.Chapters

		for _, chapter := range chapters {
			chaptersString = chaptersString + fmt.Sprintf("%sDELIM%s\n", chapter.ChapterRenderer.Title.SimpleText, strconv.Itoa(chapter.ChapterRenderer.TimeRangeStartMillis/1000))
		}
	}

	if skipThumbnails {
		//Print("saving to file")
		//os.WriteFile("THISraw.json", []byte(initialPlayerResponse), 0666)
		//thing, _ := json.MarshalIndent(initialData, "", "  ")
		//os.WriteFile("THISdone.json", []byte(string(thing)), 0666)
		//Print("saved to file")
		//var dirName string = "/tmp/" + strconv.Itoa(os.Getpid())
		//os.WriteFile(dirName + "/playbackTracking" + videoID + ".txt", []byte(mainVideo.VideoStatsPlaybackURL + "\n" + mainVideo.VideoStatsWatchtimeURL), 0666)
		//os.WriteFile(dirName + "/chapters" + videoID + ".txt", []byte(chaptersString), 0666)
		fmt.Println(mainVideo.VideoStatsPlaybackURL)
		fmt.Println(mainVideo.VideoStatsWatchtimeURL)
		fmt.Println(chaptersString)
	}

	if !skipThumbnails && config.ActiveConfig.Thumbnails {
		for i := 0; i < number; i++ {
			_ = <-doneChan
		}
	}

	return mainVideo, videoHolder
}
