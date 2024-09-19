package youtube

import (
	"os/exec"
)

const NONE int = 0
const ERROR int = -1
const RESIZE = 100

// Types, these are used to define what is in a struct
const SEARCH int = 1
const HOME int = 2
const SUBS int = 3
const HISTORY int = 4
const LIBRARY int = 5
const MY_PLAYLIST int = 6
const WL int = 7
const LIKED int = 8
const VIDEO_PAGE int = 9
const VIDEO int = 10
const OTHER_PLAYLIST int = 11

// These are used to define what action should be performed.
const PERFORM_SEARCH int = 1
const GET_HOME int = 2
const GET_SUBS int = 3
const GET_HISTORY int = 4
const GET_LIBRARY int = 5
const GET_PLAYLIST int = 6
const GET_WL int = 7
const GET_LIKED int = 8
const PLAY_VIDEO int = 9
const PLAY_VIDEO_BACKGROUND int = 10
const EXIT int = 15
const PREVIOUS int = 16
const NEXT int = 17

const PRIVATE int = 1
const UNLISTED int = 2
const PUBLIC int = 3

const CACHE_FOLDER string = "/.cache/gotube/"
const DATA_FOLDER string = "/.local/share/gotube/"
const CONFIG_FOLDER string = "/.config/gotube/"
const FRECENCY_PLAYLISTS_FILE string = "playlistsFrecency.txt"

var HOME_DIR string

func DecodeVisibility(visibilityInt int) string {
	if visibilityInt == PRIVATE {
		return "PRIVATE"
	} else if visibilityInt == UNLISTED {
		return "UNLISTED"
	} else if visibilityInt == PUBLIC {
		return "PUBLIC"
	}
	return "unknown"
}

func EncodeVisibility(visibilityString string) int {
	if visibilityString == "Private" {
		return PRIVATE
	} else if visibilityString == "Unlisted" {
		return UNLISTED
	} else if visibilityString == "Public" {
		return PUBLIC
	}
	return 0
}

func (vidHolder VideoHolder) GetVidHolder() VideoHolder {
	return vidHolder
}

func Print(input string) {
	exec.Command("notify-send", input).Run()
}

func (vidHolder *VideoHolder) SetVideosList(newSlice []Video) {
	vidHolder.Videos = newSlice
}

type VideoHolder struct {
	Videos            []Video
	PageType          int
	PlaylistID        string
	ContinuationToken string
}

type Format struct {
	VideoURL string
	AudioURL string
}

type Video struct {
	Title         string
	Id            string
	ThumbnailLink string
	ThumbnailFile string
	Channel       string
	ChannelID     string
	// Video Unique Stuff
	Views                string
	VidType              string
	ReleaseDate          string
	Length               string
	DirectLink           string
	DirectLinkAudio      string
	StartTime            int
	PlaylistRemoveId     string
	PlaylistRemoveParams string
	// Playlist Unique Stuff
	LastUpdated string
	NumVideos   int
	Visibility  string
	Type        int
}

type VideoPage struct {
	Title                  string
	Views                  string
	ViewsShort             string
	Likes                  string
	VidType                string
	ReleaseDate            string
	ReleaseDateShort       string
	Length                 string
	LengthShort            string
	Id                     string
	Channel                string
	ChannelID              string
	ChannelThumbnailLink   string
	ChannelThumbnailFile   string
	ThumbnailLink          string
	ThumbnailFile          string
	DirectLink             string
	Description            string
	LikeStatus             string
	SubStatus              string
	SubParam               string
	UnSubParam             string
	AddLikeParam           string
	RemoveLikeParam        string
	AddDislikeParam        string
	RemoveDislikeParam     string
	VideoStatsPlaybackURL  string
	VideoStatsWatchtimeURL string
}
