package download

import (
	"encoding/json"
	"fmt"
	"gotube/config"
	"gotube/download/network"
	"gotube/youtube"
	"math/rand"
	"strings"
	"os"
	"strconv"
)

var _ = os.ReadFile

// This file contains various functions for interacting with the YouTube API

const API_KEY string = "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
const PLAYLIST_ADD_URL = "https://www.youtube.com/youtubei/v1/browse/edit_playlist?key=" + API_KEY
const BROWSE_URL = "https://www.youtube.com/youtubei/v1/browse?key=" + API_KEY
const GET_ADD_TO_PLAYLIST_URL string = "https://www.youtube.com/youtubei/v1/playlist/get_add_to_playlist?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
const ADD_LIKE_URL string = "https://www.youtube.com/youtubei/v1/like/like?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
const REMOVE_LIKE_URL string = "https://www.youtube.com/youtubei/v1/like/removelike?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
const CREATE_PLAYLIST_URL string = "https://www.youtube.com/youtubei/v1/playlist/create?prettyPrint=false&key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
//const CREATE_PLAYLIST_LOG_URL string = "https://www.youtube.com/youtubei/v1/att/log?prettyPrint=false&key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
const DELETE_PLAYLIST_URL string = "https://www.youtube.com/youtubei/v1/playlist/delete?prettyPrint=false&key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"

type markWatchedKeys struct {
	cpn    string
	docid  string
	ns     string
	el     string
	uga    string
	ver    string
	st     string
	cl     string
	ei     string
	plid   string
	length string
	of     string
	vm     string
	cmt    string
	et     string
}

const cpnOptions string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

func InitThumbnailDir() {
	ThumbnailDir = THUMBNAIL_DIR_START + strconv.Itoa(config.ActiveConfig.PID) + "/"
}

// Helper function for the process of marking a video as watched
func getCPN() string {
	var toReturn string = ""
	for i := 0; i < 16; i++ {
		toReturn += string(cpnOptions[rand.Intn(63)])
	}
	return toReturn
}

// Extract info from the url contained in the page data
func extractInfo(url string, time string) markWatchedKeys {
	split := strings.Split(url, "&")
	var toReturn markWatchedKeys = markWatchedKeys{
		cpn:    getCPN(),
		docid:  strings.Split(split[1], "=")[1],
		ns:     strings.Split(split[4], "=")[1],
		el:     strings.Split(split[6], "=")[1],
		uga:    strings.Split(split[len(split)-2], "=")[1],
		ver:    "2",
		st:     "0",
		cl:     strings.Split(split[0], "=")[1],
		ei:     strings.Split(split[2], "=")[1],
		plid:   strings.Split(split[5], "=")[1],
		length: strings.Split(split[7], "=")[1],
		of:     strings.Split(split[8], "=")[1],
		vm:     strings.Split(split[len(split)-1], "=")[1],
		cmt:    time + ".0",
		et:     time + ".0",
	}

	return toReturn
}

// Mark a video as watched
func MarkWatched(videoId string, videoStatsPlaybackURL string, videoStatsWatchtimeURL string, time string) {
	config.LogEvent("Adding video: " + videoId + " to history with time: " + time)
	time = strings.Split(time, ".")[0]

	playbackData := extractInfo(videoStatsPlaybackURL, time)
	watchtimeData := extractInfo(videoStatsWatchtimeURL, time)
	//Print(watchtimeData.cmt)

	var playbackURL string = fmt.Sprintf("https://s.youtube.com/api/stats/playback?cl=%s&docid=%s&ei=%s&ns=%s&plid=%s&el=%s&len=%s&of=%s&uga=%s&vm=%s&ver=%s&cpn=%s&cmt=%s", playbackData.cl, playbackData.docid, playbackData.ei, playbackData.ns, playbackData.plid, playbackData.el, playbackData.length, playbackData.of, playbackData.uga, playbackData.vm, playbackData.ver, playbackData.cpn, playbackData.cmt)

	var watchtimeURL string = fmt.Sprintf("https://s.youtube.com/api/stats/watchtime?cl=%s&docid=%s&ei=%s&ns=%s&plid=%s&el=%s&len=%s&of=%s&uga=%s&vm=%s&ver=%s&cpn=%s&cmt=%s&st=%s&et=%s", watchtimeData.cl, watchtimeData.docid, watchtimeData.ei, watchtimeData.ns, watchtimeData.plid, watchtimeData.el, watchtimeData.length, watchtimeData.of, watchtimeData.uga, watchtimeData.vm, watchtimeData.ver, watchtimeData.cpn, watchtimeData.cmt, watchtimeData.st, watchtimeData.et)

	config.FileDump("PlaybackURLFinal.txt", playbackURL, false)
	config.FileDump("WatchtimeURLFinal.txt", playbackURL, false)

	network.GetHTML(playbackURL, true)
	network.GetHTML(watchtimeURL, true)
}

func getPLAddRemove(playlistId string) PLAddRemove {
	var fullHTML string = network.GetHTML(PLAYLIST_URL + playlistId, true)
	config.FileDump("PLAddRemoveRaw.html", fullHTML, false)
	var jsonText string = network.ExtractJSON(fullHTML, true)
	config.FileDump("PLAddRemoveRaw.json", jsonText, false)

	//os.WriteFile("pl.json", []byte(jsonText), 0666)

	var jsonA PLAddRemove
	if err := json.Unmarshal([]byte(jsonText), &jsonA); err != nil {
		panic(err)
	}

	text, _ := json.MarshalIndent(jsonA, "", "  ")
	config.FileDump("PLAddRemoveProcessed.json", string(text), false)

	return jsonA
}

// Add a playlist to library
func AddToLibrary(playlistId string) bool {
	config.LogEvent("Adding playlist: " + playlistId + " to library")
	jsonA := getPLAddRemove(playlistId)

	jsonString := `{
		"context": {
			"client": {
				"clientName":"WEB",
				"clientVersion":"2.20240624.06.00"
			}
		},
		"target": {
			"playlistId":"PLAYLIST_ID_PLACEHOLDER"
		},
		"params":"PARAMS_PLACEHOLDER"
	}`

	jsonString = strings.ReplaceAll(strings.ReplaceAll(jsonString, "PLAYLIST_ID_PLACEHOLDER", playlistId), "PARAMS_PLACEHOLDER", jsonA.Header.PlaylistHeaderRenderer.SaveButton.ToggleButtonRenderer.DefaultServiceEndpoint.LikeEndpoint.LikeParams)

	status, response := network.PostRequestAPI(jsonString, ADD_LIKE_URL, "https://www.youtube.com/playlist?list=" + playlistId)

	config.FileDump("AddToLibraryResponse.json", response, false)

	if status == 200 {
		return true
	}
	return false
}

// Remove a playlist from library
func RemoveFromLibrary(playlistId string) bool {
	config.LogEvent("Removing playlist: " + playlistId + " from library")
	jsonA := getPLAddRemove(playlistId)

	jsonString := `{
		"context": {
			"client": {
				"clientName":"WEB",
				"clientVersion":"2.20240624.06.00"
			}
		},
		"target": {
			"playlistId":"PLAYLIST_ID_PLACEHOLDER"
		},
		"params":"PARAMS_PLACEHOLDER"
	}`

	jsonString = strings.ReplaceAll(strings.ReplaceAll(jsonString, "PLAYLIST_ID_PLACEHOLDER", playlistId), "PARAMS_PLACEHOLDER", jsonA.Header.PlaylistHeaderRenderer.SaveButton.ToggleButtonRenderer.ToggledServiceEndpoint.LikeEndpoint.DislikeParams)

	status, response := network.PostRequestAPI(jsonString, REMOVE_LIKE_URL, "https://www.youtube.com/playlist?list=" + playlistId)

	config.FileDump("RemoveFromLibraryResponse.json", response, false)

	if status == 200 {
		return true
	}
	return false
}

func CreatePlaylist(playlistName string, visibility int) ([]string, bool) {
	config.LogEvent("Creating Playlist: " + playlistName)
	jsonString := `{
		"context": {
			"client": {
				"clientName":"WEB",
				"clientVersion":"2.20240624.06.00"
			}
		},
		"title":"TITLE_PLACEHOLDER",
		"privacyStatus": "PRIVACY_PLACEHOLDER",
		"videoIDs": []
	}`
	
	var visibilityString string = youtube.DecodeVisibility(visibility)
	if visibilityString == "unknown" {
		panic("Unknown privacy, most likely a coding error")
	}
	
	jsonString = strings.ReplaceAll(strings.ReplaceAll(jsonString, "TITLE_PLACEHOLDER", playlistName), "PRIVACY_PLACEHOLDER", visibilityString)
	
	status, response := network.PostRequestAPI(jsonString, CREATE_PLAYLIST_URL, "https://www.youtube.com/")

	config.FileDump("CreatePlaylistResponseRaw.json", response, false)
	
	var responseJSON CreatePlaylistResponseJSON;
	if err := json.Unmarshal([]byte(response), &responseJSON); err != nil {
		return []string{}, false
	}
	if responseJSON.Error.Code != 0 {
		return []string{}, false
	}
	
	text, _ := json.MarshalIndent(responseJSON, "", "  ")
	config.FileDump("CreatePlaylistResponseProcessed.json", string(text), false)
	
	var createdPlaylistId string = responseJSON.PlaylistId
	var createdPlaylistName string = responseJSON.Actions[0].AddToGuideSectionAction.Items[0].GuideEntryRenderer.FormattedTitle.SimpleText

	//youtube.Print("Code: " + strconv.Itoa(responseJSON.Error.Code))


	if status == 200 {
		return []string{createdPlaylistId, createdPlaylistName}, true
	}
	return []string{}, false
}

func DeletePlaylist(playlistID string) bool {
	config.LogEvent("Deleting Playlist: " + playlistID)
	jsonString := `{
		"context": {
			"client": {
				"clientName":"WEB",
				"clientVersion":"2.20240624.06.00"
			}
		},
		"playlistId":"PLAYLIST_ID_PLACEHOLDER"
	}`

	jsonString = strings.ReplaceAll(jsonString, "PLAYLIST_ID_PLACEHOLDER", playlistID)
	status, response := network.PostRequestAPI(jsonString, DELETE_PLAYLIST_URL, "https://www.youtube.com/playlist?list=" + playlistID)

	config.FileDump("DeletePlaylistResponse.json", response, false)

	if status == 200 {
		return true
	}
	return false
}


func AddToPlaylist(videoID string, playlistID string) bool {
	config.LogEvent("Adding video: " + videoID + " to playlist: " + playlistID)
	jsonString := `{
		"context": {
			"client": {
				"clientName":"WEB",
				"clientVersion":"2.20240624.06.00"
			}
		},
		"actions": [
		{
			"addedVideoId":"VIDEO_ID_PLACEHOLDER",
			"action":"ACTION_ADD_VIDEO"
		}
		],
		"playlistId":"PLAYLIST_ID_PLACEHOLDER"
	}`

	jsonString = strings.ReplaceAll(strings.ReplaceAll(jsonString, "VIDEO_ID_PLACEHOLDER", videoID), "PLAYLIST_ID_PLACEHOLDER", playlistID)
	status, response := network.PostRequestAPI(jsonString, PLAYLIST_ADD_URL, "https://www.youtube.com/watch?v=" + videoID)

	config.FileDump("AddToPlaylistResponse.json", response, false)

	if status == 200 {
		return true
	}
	return false
}

func AddToPlaylistMany(videoIDs []string, playlistID string) bool {
	actionsString := ""
	for _, id := range videoIDs {
		actionsString += "{ \"addedVideoId\": \"" + id + "\", \"action\": \"ACTION_ADD_VIDEO\" },"
	}
	actionsString = actionsString[:len(actionsString)-1]
	
	jsonString := `{
		"context": {
			"client": {
				"clientName":"WEB",
				"clientVersion":"2.20240624.06.00"
			}
		},
		"actions": [
			HERE
		],
		"playlistId":"PLAYLISTID"
	}`
	
	jsonString = strings.ReplaceAll(strings.ReplaceAll(jsonString, "HERE", actionsString), "PLAYLISTID", playlistID)
	status, response := network.PostRequestAPI(jsonString, PLAYLIST_ADD_URL, "https://www.youtube.com/watch?v=" + videoIDs[0])

	config.FileDump("AddToPlaylistResponse.json", response, true)

	if status == 200 {
		return true
	}
	return false
}

func RemoveFromPlaylist(videoID string, playlistID string, removeID string, removeParams string) bool {
	config.LogEvent("Removing video: " + videoID + " from playlist: " + playlistID)
	jsonString := `{
		"context": {
			"client": {
				"clientName":"WEB",
				"clientVersion":"2.20240624.06.00"
			}
		},
		"actions": [
		{
			"action":"ACTION_REMOVE_VIDEO",
			"setVideoId":"REMOVE_ID_PLACEHOLDER"
		}
		],
		"params": "REMOVE_PARAMS_PLACEHOLDER",
		"playlistId":"PLAYLIST_ID_PLACEHOLDER"
	}`

	jsonString = strings.ReplaceAll(jsonString, "PLAYLIST_ID_PLACEHOLDER", playlistID)
	jsonString = strings.ReplaceAll(strings.ReplaceAll(jsonString, "REMOVE_ID_PLACEHOLDER", removeID), "REMOVE_PARAMS_PLACEHOLDER", removeParams)
	status, response := network.PostRequestAPI(jsonString, PLAYLIST_ADD_URL, "https://www.youtube.com/playlist?list=" + playlistID)

	config.FileDump("RemoveFromPlaylistResponse.json", response, false)

	if status == 200 {
		return true
	}
	return false
}

func GetAddToPlaylist(videoID string) (map[string]string, []string) {
	config.LogEvent("Getting AddToPlaylist information for video: " + videoID)
	jsonString := `{
		"context": {
			"client": {
				"clientName": "WEB",
				"clientVersion": "2.20240624.06.00"
			}
		},
		"videoIds": [
			"VIDEO_ID_PLACEHOLDER"
		],
		"excludeWatchLater": false
	}`

	jsonString = strings.ReplaceAll(jsonString, "VIDEO_ID_PLACEHOLDER", videoID)

	var status int = 0
	var count int = 0
	var returnedJSONString string
	for {
		status, returnedJSONString = network.PostRequestAPI(jsonString, GET_ADD_TO_PLAYLIST_URL, "https://www.youtube.com/watch?v=" + videoID)
		if status == 200 {
			break
		}
		count++
		config.LogWarning(fmt.Sprintf("Retrying GetAddToPlaylist (count: %d)", count))
	}
	config.FileDump("GetAddToPlaylistRaw.json", returnedJSONString, false)

	var infoJSON PlaylistAddDataJSON
	if err := json.Unmarshal([]byte(returnedJSONString), &infoJSON); err != nil {
		panic(err)
	}

	text, _ := json.MarshalIndent(infoJSON, "", "  ")
	config.FileDump("GetAddToPlaylistProcessed.json", string(text), false)

	_ = status

	playlistsMap := make(map[string]string)
	playlistsSlice := []string{}
	contents := infoJSON.Contents[0].AddToPlaylistRenderer.Playlists
	for _, entry := range contents {
		info := entry.PlaylistAddToOptionRenderer
		playlistName := info.Title.SimpleText

		// Check if the playlist already contains the target video (only works for single selected video)
		if info.ContainsSelectedVideos  == "ALL" {
			playlistName = "* " + playlistName
		}
		
		playlistsSlice = append(playlistsSlice, playlistName)
		playlistsMap[playlistName] = info.PlaylistID
	}
	
	// Watch later should always be first on the list, then all the other playlist in order of most recently added to. Sometimes Watch later will be the second in the list, this code fixes that
	if len(playlistsSlice) > 1 && playlistsSlice[1] == "Watch later" {
		temp := playlistsSlice[0]
		playlistsSlice[0] = playlistsSlice[1]
		playlistsSlice[1] = temp
	}
	
	return playlistsMap, playlistsSlice
}

func GetDirectLinks(videoID string) map[string]youtube.Format {
	config.LogEvent("Getting direct links for video: " + videoID)
	structJSON := &network.PostJSON{
		VideoID:        videoID,
		Params:         "CgIIAQ==",
		ContentCheckOK: true,
		RacyCheckOK:    true,
	}

	// https://github.com/trizen/pipe-viewer/commit/729f44744851ee8b11afd136806d6beae5571f65 (09/03/2024)

	structJSON.Context.Client.ClientName = "ANDROID"
	structJSON.Context.Client.ClientVersion = "19.09.37"
	structJSON.Context.Client.AndroidSDKVersion = 30
	structJSON.Context.Client.UserAgent = "com.google.android.youtube/19.09.37 (Linux; U; Android 11) gzip"
	structJSON.Context.Client.HL = "en"
	structJSON.Context.Client.TimeZone = "UTC"
	structJSON.Context.Client.UTCOffsetMinutes = 0
	structJSON.PlaybackContext.ContentPlaybackContext.HTML5Preference = "HTML5_PREF_WANTS"

	var jsonString string = network.PostRequest(structJSON)
	config.FileDump("GetDirectLinksRaw.json", jsonString, false)
	//os.WriteFile("raw.json", []byte(jsonString), 0666)

	var jsonFormats DLResponse
	if err := json.Unmarshal([]byte(jsonString), &jsonFormats); err != nil {
		panic(err)
	}

	text, _ := json.MarshalIndent(jsonFormats, "", " ")
	//os.WriteFile("processed.json", text, 0666)
	config.FileDump("GetDirectLinksProcessed.json", string(text), false)

	qualityOptions := make(map[string]youtube.Format)

	for _, entry := range jsonFormats.StreamingData.Formats {
		if entry.Itag == 22 {
			qualityOptions["720p"] = youtube.Format{VideoURL: entry.URL, AudioURL: "combined"}
		} else if entry.Itag == 18 {
			qualityOptions["360p"] = youtube.Format{VideoURL: entry.URL, AudioURL: "combined"}
		}
	}

	var audioLinkM4A string
	var audioLinkOpus string
	for _, entry := range jsonFormats.StreamingData.AdaptiveFormats {
		if entry.Itag == 140 {
			audioLinkM4A = entry.URL
		} else if entry.Itag == 251 {
			audioLinkOpus = entry.URL
		}
	}

	for _, entry := range jsonFormats.StreamingData.AdaptiveFormats {
		if entry.Itag == 313 || entry.Itag == 315 {
			qualityOptions["2160p"] = youtube.Format{VideoURL: entry.URL, AudioURL: audioLinkOpus}
		} else if entry.Itag == 271 || entry.Itag == 308 {
			qualityOptions["1440p"] = youtube.Format{VideoURL: entry.URL, AudioURL: audioLinkOpus}
		} else if entry.Itag == 248 || entry.Itag == 303 {
			qualityOptions["1080p"] = youtube.Format{VideoURL: entry.URL, AudioURL: audioLinkOpus}
		} else if entry.Itag == 397 {
			qualityOptions["480p"] = youtube.Format{VideoURL: entry.URL, AudioURL: audioLinkM4A}
		}
	}

	//Print("Length: " + strconv.Itoa(len(qualityOptions)))
	// If it's a livestream, just return that
	if jsonFormats.StreamingData.HLSManifestURL != "" {
		qualityOptions = make(map[string]youtube.Format)
		qualityOptions["1080p"] = youtube.Format{VideoURL: jsonFormats.StreamingData.HLSManifestURL, AudioURL: "combined"}
		//Print("Livestream: " + jsonFormats.StreamingData.HLSManifestURL)
	}

	return qualityOptions
}
