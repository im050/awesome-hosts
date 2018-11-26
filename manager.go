package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func GetCurrentHosts(file *os.File) []Host {
	br := bufio.NewReader(file)
	//each file line by line
	for {
		line, _, err := br.ReadLine()
		lineString := strings.TrimSpace(string(line))
		if err == io.EOF {
			break
		}
		//if empty, continue
		if len(lineString) == 0 {
			continue
		}
		//if notice, continue
		if strings.Index(lineString, "#") == 0 {
			continue
		}
		hostSplit := strings.Split(lineString, " ")
		//if domain nonexistent, continue
		if len(hostSplit) < 2 {
			continue
		}
		hosts = append(hosts, Host{
			Domain: hostSplit[0],
			IP: hostSplit[1],
		})
	}
	return hosts
}
