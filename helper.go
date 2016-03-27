package web

import (
	"os"
	"strings"
	"time"
)

func FileExists(fileName string) bool {

	info, err := os.Stat(fileName)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
func DirExists(dir string) bool {

	info, err := os.Stat(dir)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func WebTime(t time.Time) string {
	ftime := t.Format(time.RFC1123)
	if strings.HasSuffix(ftime, "UTC") {
		ftime = ftime[0:len(ftime)-3] + "GMT"
	}
	return ftime
}
