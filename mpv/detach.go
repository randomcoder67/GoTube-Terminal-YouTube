package mpv

import (
	"fmt"
	"gotube/config"
	"gotube/download"
	"gotube/youtube"
	"os"
	"os/exec"
	"strconv"
)

func Print(str string) {
	cmd := exec.Command("notify-send", str)
	cmd.Run()
}

func WritePlaylistFile(videoHolder youtube.VideoHolder) {
	for videoHolder.ContinuationToken != "" {
		videoHolder = download.GetPlaylistContinuation(videoHolder, videoHolder.ContinuationToken)
	}

	// Get relevant data
	var detailsFile string = ""
	var playlistFile string = ""
	for _, video := range videoHolder.Videos {
		detailsFile = detailsFile + fmt.Sprintf("%sDELIM%sDELIM%s\n", video.Title, video.Channel, video.Id)
		playlistFile = playlistFile + fmt.Sprintf("%s - %s\n", video.Title, video.Channel)
	}

	// Write data to files
	var dirName string = "/tmp/gotube_" + strconv.Itoa(os.Getpid())
	config.Mkdir(dirName)
	os.WriteFile(dirName + "/playlist.m3u", []byte(playlistFile), 0666)
	os.WriteFile(dirName + "/details.txt", []byte(detailsFile), 0666)
}

/*
func Playlist() {
	var videos youtube.VideoHolder = download.GetPlaylist("WL", "Watch later")
	for videos.ContinuationToken != "" {
		videos = download.GetPlaylistContinuation(videos, videos.ContinuationToken)
	}

	for _, video := range videos.Videos {
		fmt.Println(video.Title, video.Channel)
	}

	var playlistDummy string = ""
	for i:=1; i<len(videos.Videos); i++ {
		playlistDummy += videos.Videos[i].Title + "DELIM" + videos.Videos[i].Id + "\n"
	}
	os.WriteFile("test.m3u", []byte(playlistDummy), 0666)


	var directLink string = download.GetDirectLinks(videos.Videos[0].Id)["720p"].VideoURL
	mpvCommandArgs := []string{"--title=" + videos.Videos[0].Title + " - " + videos.Videos[0].Channel, directLink, "--playlist=test.m3u", "--script=mpv/addFile.lua"}
	mpvCommand := exec.Command("mpv", mpvCommandArgs...)

	mpvCommand.Start()
}
*/

func DetachVideo(title string, channel string, startTime string, startNum string, folderName string, quality string, geometryArg string) {
	forkCommand := exec.Command("nohup", youtube.HOME_DIR + "/.local/bin/" + "gotube", "--fork", "--play", title, channel, startTime, startNum, folderName, quality, geometryArg)
	//var errb bytes.Buffer
	//forkCommand.Stderr = &errb
	//forkCommand.Run()
	//os.WriteFile("command.err", errb.Bytes(), 0666)
	
	forkCommand.Start()
}
