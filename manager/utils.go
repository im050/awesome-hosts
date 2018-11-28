package manager

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetHostsFile() string {
	switch runtime.GOOS {
	case "darwin":
		return "/etc/hosts"
	case "windows":
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux":
		return "/etc/hosts"
	default:
		return ""
	}
}

func GetLineSeparator() string {
	switch runtime.GOOS {
	case "darwin":
		return "\n"
	case "windows":
		return "\r\n"
	case "linux":
		return "\n"
	default:
		return "\n"
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
