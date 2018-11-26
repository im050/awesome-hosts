package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
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
		enabled := true
		if strings.Index(lineString, "#") == 0 {
			enabled = false
		}
		lineString = regexp.MustCompile(`\t+`).ReplaceAllLiteralString(lineString, " ")
		reg := regexp.MustCompile(`\s+`)
		hostSplit := reg.Split(lineString, -1)
		//if domain nonexistent, continue
		if len(hostSplit) < 2 {
			continue
		}
		if !enabled {
			hostSplit[0] = strings.TrimSpace(strings.TrimLeft(hostSplit[0], "#"))
		}
		IPv4Pattern := `((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`
		IPv6Pattern := `(([\da-fA-F]{1,4}):){8}`
		//if not ip, continue
		if !regexp.MustCompile(IPv4Pattern).MatchString(hostSplit[0]) && !regexp.MustCompile(IPv6Pattern).MatchString(hostSplit[0]) {
			continue
		}
		fmt.Println(hostSplit)
		hosts = append(hosts, Host{
			Domain:  hostSplit[1],
			IP:      hostSplit[0],
			Enabled: enabled,
		})
	}
	return hosts
}
