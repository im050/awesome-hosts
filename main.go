package main

import (
	"fmt"
	"github.com/zserge/webview"
	"os"
	"runtime"
	"time"
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
	fmt.Println("http://" + ln.Addr().String())
	openWebview()
	//checkError(webview.Open("Host Manager", "http://" + ln.Addr().String(), 850, 630, false))
	wait()
}

func openWebview() {
	wv := webview.New(webview.Settings{
		Title:                  "Host Manager",
		URL:                    "http://" + ln.Addr().String(),
		Width:                  850,
		Height:                 630,
		Resizable:              true,
		Debug:                  false,
		ExternalInvokeCallback: nil,
	})
	wv.Run()
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

func wait() {
	for {
		time.Sleep(time.Second * 86400)
	}
}
