package web

import "os"

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
