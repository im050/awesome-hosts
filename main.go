package main

import (
	"os"
)

var (
	hosts []Host
	)

func main() {
	file, err := os.Open("C:\\Windows\\System32\\drivers\\etc\\hosts")
	hosts = GetCurrentHosts(file)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	ServerStart()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}