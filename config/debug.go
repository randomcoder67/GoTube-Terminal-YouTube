package config

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

const LOG_DIR_NAME string = "/.cache/gotube/log/"

var HOME_DIR string

var logFileD *os.File

func timeNowLong() string {
	return time.Now().Format("2006-01-02_15-04-05")
}

func Mkdir(path string) bool {
	err := os.Mkdir(path, 0755)
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
}

func FileDump(fileName string, content string, force bool) {
	if ActiveConfig.DumpJSON || force {
		var dirName string = HOME_DIR + LOG_DIR_NAME + strconv.Itoa(ActiveConfig.PID)
		if !Mkdir(dirName) {
			panic("Could not make dir: " + dirName)
		}
		var fullFileName = timeNowLong() + "-" + strconv.Itoa(ActiveConfig.PID) + "-" + fileName
		err := os.WriteFile(dirName + "/" + fullFileName, []byte(content), 0666)
		if err != nil {
			panic(err)
		}
	}
}

func LogEvent(input string) {
	if ActiveConfig.Log {
		fmt.Fprintf(logFileD, "%s: %s\n", timeNowLong(), input)
	}
}

func LogWarning(input string) {
	fmt.Fprintf(logFileD, "%s: %s\n", timeNowLong(), input)
}

// When logging an error, the stacktrace is also saved to the log
func LogError(input string) {
	var thing []uint8 = debug.Stack()
	fmt.Fprintf(logFileD, "%s: %s\nStacktrace:\n%s", timeNowLong(), input, thing)
}

func OpenLogFile() {
	HOME_DIR, _ = os.UserHomeDir()

	var fileName string = timeNowLong() + "-" + strconv.Itoa(os.Getpid()) + ".log"
	var err error
	logFileD, err = os.OpenFile(HOME_DIR + LOG_DIR_NAME + fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
}

func CloseLogFile() {
	fmt.Fprintf(logFileD, "Closed Log Safely\n")
	logFileD.Close()
}
