package network

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"gotube/config"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

// This file contains every functions which actually makes network requests

//const YOUTUBE_API_URL = "https://youtubei.googleapis.com/youtubei/v1/player?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
const YOUTUBE_API_URL = "https://www.youtube.com/youtubei/v1/player?key=AIzaSyA8eiZmM1FaDVjRy-df2KTyQ_vz_yYM39w&prettyPrint=false"
const ORIGIN_URL string = "https://www.youtube.com"

const RETRY_COUNT = 3

type PostJSON struct {
	VideoID        string `json:"videoId"`
	ContentCheckOK bool   `json:"contentCheckOk"`
	RacyCheckOK    bool   `json:"racyCheckOk"`
	Params         string `json:"params"`
	Context        struct {
		Client struct {
			ClientName        string `json:"clientName"`
			ClientVersion     string `json:"clientVersion"`
			AndroidSDKVersion int    `json:"androidSdkVersion"`
			UserAgent         string `json:"userAgent"`
			HL                string `json:"hl"`
			TimeZone          string `json:"timeZone"`
			UTCOffsetMinutes  int    `json:"utcOffsetMinutes"`
		} `json:"client"`
	} `json:"context"`
	PlaybackContext struct {
		ContentPlaybackContext struct {
			HTML5Preference string `json:"html5Preference"`
		} `json:"contentPlaybackContext"`
	} `json:"playbackContext"`
}

// Function to download a thumbnail, will be ran in a goroutine. Sends an int on chan finished when done
func DownloadThumbnail(url string, filename string, resize bool, finished chan int, edit bool) {
	var resp *http.Response
	var err error
	for i := 0; i < RETRY_COUNT; i++ {
		resp, err = http.Get(url)
		if err == nil {
			break
		} else if err != nil && i == 2 {
			panic(err)
		}
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filename, content, 0666)
	if err != nil {
		panic(err)
	}
	finished <- 1
}

// Basic function for GET request - used for downloading pages
func GetHTML(url1 string, cookies bool) string {
	jar := GetCookies()
	client := http.Client{
		Jar: jar,
	}
	
	if config.ActiveConfig.BrowserCookies == "chromium" {
		jar := GetCookiesChromium()
		client = http.Client{
			Jar: jar,
		}
	}

	req, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.18 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("accept-langauge", "en-us,en;q=0.5")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-YouTube-Client-Name", "3")
	req.Header.Set("X-YouTube-Client-Version", "19.09.37")
	req.Header.Set("origin", "https://www.youtube.com")
	
	
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	responseHTML, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(responseHTML)
}

// Basic post request function, takes request JSON as input, headers set within
func PostRequest(structJSON *PostJSON) string {
	client := http.Client{}

	if config.ActiveConfig.BrowserCookies == "firefox" {
		jar := GetCookies()
		client = http.Client{
			Jar: jar,
		}
	} else if config.ActiveConfig.BrowserCookies == "chromium" {
		jar := GetCookiesChromium()
		client = http.Client{
			Jar: jar,
		}
	}

	properJSON, err := json.Marshal(structJSON)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", YOUTUBE_API_URL, bytes.NewBuffer(properJSON))
	req.Header.Set("user-agent", "com.google.android.youtube/19.09.37 (Linux; U; Android 11) gzip")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("accept-langauge", "en-us,en;q=0.5")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-YouTube-Client-Name", "3")
	req.Header.Set("X-YouTube-Client-Version", "19.09.37")
	req.Header.Set("origin", "https://www.youtube.com")
	//req.Header.Set("accept-encoding", "gzip, deflate")

	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	responseHTML, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(responseHTML)
}

// Post request for the YouTube website API, requires some more stuff including the SAPIDID hash so this is a seperate function
func PostRequestAPI(jsonString string, url string, refererURL string) (int, string) {
	// Get cookies and form client
	client := http.Client{}
	var jar *cookiejar.Jar
	
	if config.ActiveConfig.BrowserCookies == "firefox" {
		jar = GetCookies()
		client = http.Client{
			Jar: jar,
		}
	} else if config.ActiveConfig.BrowserCookies == "chromium" {
		jar = GetCookiesChromium()
		client = http.Client{
			Jar: jar,
		}
	}

	// Get sapisid and hash
	var sapisid string = getSapis(jar)
	curTime := time.Now()
	var toHashString string = fmt.Sprintf("%d %s %s", curTime.Unix(), sapisid, ORIGIN_URL)
	var hashPart string = fmt.Sprintf("%x", sha1.Sum([]byte(toHashString)))
	var sapisidHash string = fmt.Sprintf("SAPISIDHASH %d_%s", curTime.Unix(), hashPart)

	config.LogEvent("SAPISID hash is: " + sapisidHash)
	// Format JSON
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonMap)
	properJSON, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}

	// Create request and add headers
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(properJSON))
	if err != nil {
		panic(err)
	}
	
	req.Header.Set("authorization", sapisidHash)
	req.Header.Set("Host", "www.youtube.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	req.Header.Set("Referer", refererURL)
	req.Header.Set("Origin", "https://www.youtube.com")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("TE", "trailers")

	// Perform request and return response HTML
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	responseHTML, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//err = os.WriteFile("out.json", responseHTML, 0666)
	//err = os.WriteFile("sapisid.txt", []byte(sapisidHash), 0666)
	return resp.StatusCode, string(responseHTML)
}
