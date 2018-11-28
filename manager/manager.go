package manager

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// `json:"xxx"` 转换json时相对应的字段
type Host struct {
	Domain     string `json:"domain"`
	IP         string `json:"ip"`
	Enabled    bool   `json:"enabled"`
	LineNumber int    `json:"lineNumber"`
}

type Manager struct {
	hostsDir    string
	SystemHosts []Host
}

func New(hostDir string) *Manager {
	m := new(Manager)
	m.hostsDir = hostDir
	return m
}

func (h *Manager) GetHosts(file *os.File) (hosts []Host) {
	br := bufio.NewReader(file)
	//each file line by line
	lineIndex := -1
	for {
		line, _, err := br.ReadLine()
		lineString := strings.TrimSpace(string(line))
		lineIndex++
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
		hosts = append(hosts, Host{
			Domain:     hostSplit[1],
			IP:         hostSplit[0],
			Enabled:    enabled,
			LineNumber: lineIndex,
		})
	}
	return hosts
}

func (h *Manager) WriteContent(name string, content string) {
	data := []byte(content)
	if ioutil.WriteFile(h.hostsDir+"/"+name, data, 0644) == nil {
		fmt.Println("写入文件成功:", content)
	}
}

func (h *Manager) WriteHosts(name string, hosts []Host) {
	hostsContent := ""
	eol := GetLineSeparator()
	for i, _ := range hosts {
		host := hosts[i]
		if !host.Enabled {
			hostsContent += "#"
		}
		hostsContent += host.IP + " " + host.Domain + eol
	}
	h.WriteContent(name, hostsContent)
}
