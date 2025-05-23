# GoTube Terminal YouTube Client

![Static Badge](https://img.shields.io/badge/Linux-grey?logo=linux)
![Static Badge](https://img.shields.io/badge/Golang-007D9C)
![Static Badge](https://img.shields.io/badge/Usage-Terminal_YouTube_Client-blue)
![GitHub License](https://img.shields.io/github/license/randomcoder67/GoTube-Terminal-YouTube)
![GitHub Release](https://img.shields.io/github/v/tag/randomcoder67/GoTube-Terminal-YouTube)

A YouTube client in the terminal with thumbnails. Allows displaying images in any terminal using ueberzug, and provides access to YouTube account features with no API keys needed.

### Screenshots

* Basic video search

![Video Search](.screenshots/search.png)

* Video page with description and recommendations

![Video Page](.screenshots/videoPage.png)

* Video watch history (synced with account)

![Watch History](.screenshots/history.png)

### Features

#### Video Finding

* Video Search
* Access to user history, subscriptions, watch later and saved playlists
* Recommendations (home page and video page)
* No API keys needed
* Video thumbnails, real images with Ueberzug

#### Video Playing

* Play videos with [mpv](https://github.com/mpv-player/mpv)
* Format options (current options: 480p, 720p, 1080p, 1440p, 2160p)
* Includes YouTube chapters in mpv

#### Video Management

* Add and remove videos from your playlists (including Watch Later)
* Saves watch time to history
* Create and delete playlists

## Version

* Currently in "Beta" state, features are present and working but may be buggy
* UI subject to change (i.e. some parts still look ugly and I may think of a better idea)

### Update (January 2025)

I remade my GitHub account to change the email address, so the original repo was deleted. I just started the releases from 0 again.  
Also I have changed the "scope" of the project a bit. I now plan on using yt-dlp for video playback. Originally I had coded into GoTube the ability to get the direct URL from a YouTube video, but this broke, and maintaining this isn't something I can guarantee. yt-dlp is a bit slower, but it will be more reliable.  
This also means the exact watch time won't be saved back to YouTube, only that the video was watched. I would like to fix this in the future.

## Installation

### Dependancies

* Firefox
* mpv
* ueberzug
* golang
* yt-dlp

### Install

`git clone https://github.com/randomcoder67/GoTube-Terminal-YouTube.git`  
`cd GoTube-Terminal-YouTube`  
`make`  
`make install`  

Log into YouTube in Firefox. Ensure there is a `~/.mozilla/firefox/[something].default-release/` directory. This should be the default save location for the cookies file when Firefox is installed using a package manager.

## Usage

### Command Line

`gotube -h/--help/help` - Display help  
`gotube` - Launch onto blank page, commands can be entered to find content  
`gotube -s/--subscriptions` - Launch and show subscriptions  
`gotube -hs/--history` - Launch and show history  
`gotube -wl/--watchlater` - Launch and show watch later  
`gotube -p/-l/--playlists/--library` [playlist_id] [playlist_name] - Launch and show library (include id and name to go straight to a playlist  
`gotube -hm/--recommendations/--home` - Launch and show home page  

### TUI

Once launched, the search box can be focused with `/` or `Tab`, and the grid can be focused with `Tab`. Navigation around the grid can be done with arrow keys or vim keys. PageUp/Down, Home and End also supported. `Ctrl-C`, `Esc` or `q` to quit. Other keyboard commands are show in the sidebar.  
Commands can be entered in the search box, beginning with a `/`. These are also shown in the sidebar.
### Wayland

* Ueberzug doesn't work on Wayland, to get the thumbnails working on Wayland, you can use the Kitty terminal and set Ueberzug to use the Kitty protocol

## Future Features

### 1.0 Release

Planned before 1.0 release:

* Channel pages
* Chrome cookies support
* More/better handled format options
* Audio only playback
* Like/dislike videos
* Sub/unsub from channels
* Wayland support - Partially done, (for now use a terminal which support kitty image protocol and use [Ueberzugpp](https://github.com/jstkdng/ueberzugpp), with config set to "kitty")
* Config file - Done
* Options menu
* Many bug fixes
* Better error handling

### Future Features

Future features, but not necessary for 1.0 release.  
Organised into varying degrees of likelyhood:

#### Planned

* Cookie support for other browsers (Brave, Vivaldi, Opera etc)
* View/interact with comments
* Change order of videos in playlist
* CLI mode - use functionality in other scripts
* Video downloading TUI (frontend for yt-dlp) - Cancelled, moving this to a seperate project
* Video queue
* Text-only mode (no thumbnails, designed for ssh or console operation, can manage videos without playing them or displaying thumbnails)

#### Likely

* [VLC](https://www.videolan.org/) support
* Search suggestions
* Alternate display mode (more detail vs current compact view)

#### Possible

* Support for other sites (Twitch, Kick, Rumble etc)
* Management of your channel
* Support for other image display protocols (Kitty for example)
* GUI version
* Windows version
* MacOS version
