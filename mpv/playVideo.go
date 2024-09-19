package mpv

import (
	"bytes"
	"fmt"
	"strings"
	"gotube/download"
	"gotube/youtube"
	"os/exec"
	//"os"
	//"strconv"
)

/*
func parseMPVVidTime(input string) string {
	splitA := strings.Split(input, "AV:")
	split := strings.Split(splitA[len(splitA)-1], ":")
	var hours string = split[0][1:]
	var minutes string = split[1]
	var seconds string = split[2][:2]
	fmt.Printf("H:%s,M:%s,S:%s\n", hours, minutes, seconds)
	hoursInt, err := strconv.Atoi(hours)
	if err != nil {
		panic(err)
	}
	minutesInt, err := strconv.Atoi(minutes)
	if err != nil {
		panic(err)
	}
	secondsInt, err := strconv.Atoi(seconds)
	if err != nil {
		panic(err)
	}

	return strconv.Itoa(hoursInt*3600 + minutesInt*60 + secondsInt)
}
*/

func GetVideoData(videoId string) {
	//Print("IN GET VIDEO DATA: " + videoId)
	download.GetVideoPage(videoId, "/dev/null", true)
	//Print("DONE GET VIDEO DATA")
}

func GetQualityLinks(videoId string, requestedQuality string) {
	/*
	qualityOptions := download.GetDirectLinks(videoId)
	var videoLink, audioLink string
	
	// Try to get desired quality then 720p then 360p the just return empty URLs (will skip video)
	qualityOrder := []string{"2160p", "1440p", "1080p", "720p", "360p"}
	
	var desiredQuality int = 0
	for i, quality := range qualityOrder {
		desiredQuality = i
		if quality == requestedQuality {
			break
		}
	}
	
	for i:=desiredQuality; i<5; i++ {
		result, ok := qualityOptions[qualityOrder[i]]
		if !ok {
			continue
		} else {
			videoLink = result.VideoURL
			audioLink = result.AudioURL
			break
		}
	}
	*/
	
	var qualityString string = ""
	
	switch requestedQuality {
		case "360p":
			qualityString = "18/bestvideo[protocol*=m3u8][height<=360][vcodec*=avc1]+bestaudio[protocol*=m3u8]/bestvideo[height<=360]+bestaudio"
		case "720p":
			qualityString = "22/bestvideo[protocol*=m3u8][height<=720][vcodec*=avc1]+bestaudio[protocol*=m3u8]/bestvideo[height<=720]+bestaudio"
		case "1080p":
			qualityString = "bestvideo[protocol*=m3u8][height<=1080][vcodec*=avc1]+bestaudio[protocol*=m3u8]/bestvideo[height<=1080]+bestaudio"
		case "1440p":
			qualityString = "bestvideo[height<=1440]+bestaudio[protocol*=m3u8]/bestvideo[height<=1440]+bestaudio"
		case "2160p":
			qualityString = "bestvideo[height<=2160]+bestaudio[protocol*=m3u8]/bestvideo[height<=2160]+bestaudio"
		default:
			qualityString = "22/bestvideo[protocol*=m3u8][height<=720][vcodec*=avc1]+bestaudio[protocol*=m3u8]/bestvideo[height<=720]+bestaudio"
	}
	
	cmd := exec.Command("yt-dlp", "-f", qualityString, "--get-url", "--", videoId)
	out, err := cmd.Output()
	
	//err = os.WriteFile("ytdlp.cmd", []byte(cmd.String()), 0666)
	//err = os.WriteFile("ytdlp.out", []byte(out), 0666)
	//fmt.Println(out, err)
	
	if err != nil {
		//Print("Error: " + err.Error())
		panic(err)
	}
	
	split := strings.Split(string(out), "\n")
	var videoLink string = split[0]
	var audioLink string = "combined"
	
	//Print("Len: " + strconv.Itoa(len(split)))
	if len(split) > 1 && split[1] != "" {
		audioLink = split[1]
	}
	
	//var test string
	//test = fmt.Sprintf("%s\n%s\n", videoLink, audioLink)
	//Print("Test: " + test)
	//err = os.WriteFile("test.out", []byte(test), 0666)
	
	fmt.Printf("%s\n%s\n", videoLink, audioLink)
}

func MarkWatched(videoId string, finalTime string, videoStatsPlaybackURL string, videoStatsWatchtimeURL string) {
	/*
		const MAX_TRIES int = 3
		var i int = 0
		var videoStatsPlaybackURL, videoStatsWatchtimeURL string
		for {
			if _, err := os.Stat(fileName); err == nil {
				i++
				dat, err := os.ReadFile(fileName)
				if err != nil {
					if i > MAX_TRIES {
						Print("PANIC")
						panic(err)
					}
					continue
				}
				split := strings.Split(string(dat), "\n")
				videoStatsPlaybackURL = split[0]
				videoStatsWatchtimeURL = split[1]
				break
			}
			time.Sleep(time.Millisecond * 500)
		}
	*/
	//Print("IN MARK WATCHED")
	//Print(videoId)
	//Print(finalTime)
	//Print(videoStatsPlaybackURL)
	//Print(videoStatsWatchtimeURL)
	download.MarkWatched(videoId, videoStatsPlaybackURL, videoStatsWatchtimeURL, finalTime)
}

func StartPlayback(title string, channel string, startTime string, startNum string, folderName string, quality string, geometryArg string) {
	mpvCommandArgs := []string{"--start=" + startTime, "--playlist-start=" + startNum, "--script=" + youtube.HOME_DIR + "/.local/bin/gotube.lua", "--resume-playback=no", "--geometry=" + geometryArg, "--script-opts-add=gotube-folderName=" + folderName + ",gotube-quality=" + quality, "--playlist=" + folderName + "/playlist.m3u", "--force-window=immediate"}

	var thing string = ""
	for _, x := range mpvCommandArgs {
		thing += x
	}
	//os.WriteFile("mpv.command", []byte(thing), 0666)

	mpvVideoCommand := exec.Command("mpv", mpvCommandArgs...)
	//os.WriteFile("mpv.command", []byte(mpvVideoCommand.String()), 0666)

	var outb, errb bytes.Buffer
	mpvVideoCommand.Stdout = &outb
	mpvVideoCommand.Stderr = &errb

	err := mpvVideoCommand.Run()
	if err != nil {
		Print(err.Error())
		//panic(err)
	}

	//os.WriteFile("mpv.out", []byte(outb.String()), 0666)
	//os.WriteFile("mpv.err", []byte(errb.String()), 0666)
	//if exitErr, ok := err.(*exec.ExitError); ok {
	//	exitCode := exitErr.ExitCode()
	//	Print("Exit Code: " + strconv.Itoa(exitCode))
	//}
}
