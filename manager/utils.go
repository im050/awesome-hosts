package manager

import (
	"log"
	"net"
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

func GetUserHome() string {
	home := ""
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	} else {
		home = os.Getenv("HOME")
	}

	return home
}

func transferGroupName(name *string, isDisplay bool) {
	if isDisplay {
		*name = strings.Replace(*name, "_", " ", -1)
	} else {
		*name = strings.Replace(*name, " ", "_", -1)
	}
}

func ErrorAndExitWithLog(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}

func GetHostFileName(name string) string {
	transferGroupName(&name, false)
	return name + ".host"
}

func GetIntranetIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Println(err)
		return ""
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
}
