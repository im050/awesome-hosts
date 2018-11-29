package manager

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type Host struct {
	Domain     string `json:"domain"`
	IP         string `json:"ip"`
	Enabled    bool   `json:"enabled"`
	LineNumber int    `json:"lineNumber"`
}

type Hosts []Host //line=>number
type Groups []Group

type Group struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Active  bool   `json:"active"`
	Hosts   Hosts  `json:"hosts"`
}

type Manager struct {
	hostsDir    string
	SystemHosts Hosts
	Groups      Groups
}

func New(hostsDir string) *Manager {
	m := new(Manager)
	m.hostsDir = hostsDir

	return m
}

func (h *Manager) Init() *Manager {
	exists, err := PathExists(h.hostsDir)

	defer h.initSystemHosts()

	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	if exists {
		return h
	}

	//create hosts dir
	err = os.Mkdir(h.hostsDir, 0777)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	return h
}

func (h *Manager) initSystemHosts() {
	file, _ := os.Open(GetHostsFile())
	defer ErrorAndExitWithLog(file.Close())
	hosts := h.GetHosts(file)
	h.SystemHosts = hosts
	exists1, err := PathExists(h.hostsDir + "/Default_Group.enable")
	exists2, err := PathExists(h.hostsDir + "/Default_Group.disable")
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	if exists1 || exists2 {
		return
	}
	h.WriteHosts("Default_Group.enable", h.SystemHosts)
}

func (h *Manager) initGroups() {
	groups := h.GetGroups()
	h.Groups = groups
}

func (h *Manager) GetHostDir() string {
	return h.hostsDir
}

func (h *Manager) GetHosts(file *os.File) Hosts {
	br := bufio.NewReader(file)
	//each file line by line
	lineIndex := -1
	var hosts []Host
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
	err := ioutil.WriteFile(h.hostsDir+"/"+name, data, 0666)
	if err != nil {
		_, _ = fmt.Println("写入文件失败", err)
		os.Exit(0)
	}
}

func (h *Manager) WriteHosts(name string, hosts Hosts) {
	hostsContent := ""
	eol := GetLineSeparator()
	for _, host := range hosts {
		if !host.Enabled {
			hostsContent += "#"
		}
		hostsContent += host.IP + " " + host.Domain + eol
	}
	h.WriteContent(name, hostsContent)
}

func (h *Manager) GetGroups() []Group {
	var groups []Group
	fmt.Println(h.hostsDir)
	files, err := ioutil.ReadDir(h.hostsDir)
	ErrorAndExitWithLog(err)
	for _, f := range files {
		groupInfo := strings.Split(f.Name(), ".")
		var enabled = false
		if groupInfo[len(groupInfo)-1] == "enable" {
			enabled = true
		}
		groupName := groupInfo[0]
		if len(groupInfo) >= 3 {
			groupName = strings.Join(groupInfo[0:], ".")
		}
		transferGroupName(&groupName, true)
		//read host file
		file, err := os.Open(h.hostsDir + "/" + f.Name())
		hosts := h.GetHosts(file)
		ErrorAndExitWithLog(file.Close())
		ErrorAndExitWithLog(err)
		//append to groups
		groups = append(groups, Group{Name: groupName, Enabled: enabled, Active: false, Hosts: hosts})
	}
	return groups
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
