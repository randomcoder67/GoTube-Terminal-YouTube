package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ConfigOpts struct {
	Log            bool   // Whether to write all events to the log (errors are always written)
	DumpJSON       bool   // Whether to dump recieved and processed data to files
	SessionType    string // X11 or Wayland, needed for copying
	PID            int
	Term           string // The current $TERM environmental variable (i.e. the terminal you are currently using)
	BrowserEnv     string // The current $BROWSER environmental variable
	BrowserCookies string // The current $BROWSER environmental variable
	Thumbnails     bool // Option to disable thumbnails for bad internet connections
}

var ActiveConfig ConfigOpts

func InitConfig(log bool, dumpJSON bool, thumbnails bool, browserCookies string) {
	ActiveConfig = ConfigOpts{
		Log:            log,
		DumpJSON:       dumpJSON,
		SessionType:    checkSessionType(),
		PID:            os.Getpid(),
		Term:           os.Getenv("TERM"),
		BrowserEnv:     os.Getenv("BROWSER"),
		BrowserCookies: browserCookies,
		Thumbnails:     thumbnails,
	}

	fmt.Fprintf(logFileD, "Config Options: %+v\n", ActiveConfig)
}

func checkSessionType() string {
	username, err := exec.Command("whoami").Output()
	var actualUsername string = strings.ReplaceAll(string(username), "\n", "")
	listSessions, err := exec.Command("loginctl", "list-sessions").Output()
	lines := strings.Split(string(listSessions), "\n")
	var session string = ""
	for _, line := range lines[1:] {
		parts := strings.Fields(line)
		if len(parts) < 3 {
			return "unknown"

		}
		if parts[2] == string(actualUsername) {
			session = parts[0]
			break
		}
	}

	sessionType, err := exec.Command("loginctl", "show-session", session, "-p", "Type").Output()
	if err != nil {
		panic(err)
	}
	var sessionTypeTrimmed string = strings.ReplaceAll((string(sessionType))[5:], "\n", "")
	return sessionTypeTrimmed
}
