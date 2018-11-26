package main

import (
	"fmt"
	"os"
	"runtime"
)

var (
	hosts []Host
)

func main() {

	file, err := os.Open(getHostsFile())
	go func() {
		fmt.Println("加载当前hosts")
		hosts = GetCurrentHosts(file)
	}()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()
	ServerStart()
}

func getHostsFile() string {
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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
