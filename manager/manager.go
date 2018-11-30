package manager

import (
	"bufio"
	"encoding/json"
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

type HostGroupData struct {
	Name                 string
	Enabled              bool
	LastUpdatedTimestamp uint32
}

type Config struct {
	Groups               []HostGroupData
	LastUpdatedTimestamp uint32 //last timestamp of hosts data was updated
	LastSyncTimestamp    uint32 //last timestamp of refresh hosts data to system hosts
}

type Manager struct {
	hostsDir         string
	ConfigFileName   string
	DefaultGroupName string
	SystemHosts      Hosts
	Groups           Groups
	Config           Config
}

func New(hostsDir string) *Manager {
	m := new(Manager)
	m.hostsDir = hostsDir
	m.DefaultGroupName = "Default Hosts"
	m.ConfigFileName = "data.config"
	return m
}

func (h *Manager) Init() *Manager {
	//third, init host groups
	defer h.initGroups()
	//second, init or load config into h.Config
	defer h.loadConfig()
	//first, backup system hosts as a new group
	defer h.initSystemHosts()
	exists, err := PathExists(h.hostsDir)
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
	defer file.Close()
	hosts := h.GetHosts(file)
	h.SystemHosts = hosts
	exists, err := PathExists(h.hostsDir + "/" + GetHostFileName(h.DefaultGroupName))
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	if exists {
		return
	}
	h.WriteHosts(GetHostFileName(h.DefaultGroupName), h.SystemHosts)
}

func (h *Manager) initGroups() {
	groups := h.GetGroups()
	h.Groups = groups
}

func (h *Manager) GetHostDir() string {
	return h.hostsDir
}

func (h *Manager) loadConfig() {
	exists, _ := PathExists(h.hostsDir + "/" + h.ConfigFileName)
	if exists {
		h.loadConfigFromFile()
	} else {
		h.Config = Config{
			Groups:               []HostGroupData{{Name: h.DefaultGroupName, Enabled: true, LastUpdatedTimestamp: 0}},
			LastUpdatedTimestamp: 0,
			LastSyncTimestamp:    0,
		}
		h.persistConfig()
	}
}

func (h *Manager) loadConfigFromFile() {
	fileByte, err := ioutil.ReadFile(h.hostsDir + "/" + h.ConfigFileName)
	ErrorAndExitWithLog(err)
	ErrorAndExitWithLog(json.Unmarshal(fileByte, &h.Config))
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
	groupEnableMap := map[string]bool{}
	for _, c := range h.Config.Groups {
		groupEnableMap[c.Name] = c.Enabled
	}
	fmt.Println(groupEnableMap)
	files, _ := ioutil.ReadDir(h.hostsDir)
	for _, f := range files {
		groupInfo := strings.Split(f.Name(), ".")
		if groupInfo[len(groupInfo)-1] != "host" {
			continue
		}
		groupName := groupInfo[0]
		if len(groupInfo) >= 3 {
			groupName = strings.Join(groupInfo[0:], ".")
		}
		transferGroupName(&groupName, true)
		var enabled bool
		value, exists := groupEnableMap[groupName]
		if !exists {
			enabled = false
			//if group doesn't exists in config file, then add it.
			h.addGroupToConfig(groupName, enabled, 0)
		} else {
			enabled = value
		}
		//read host file
		file, err := os.Open(h.hostsDir + "/" + f.Name())
		hosts := h.GetHosts(file)
		ErrorAndExitWithLog(file.Close())
		ErrorAndExitWithLog(err)
		//append to groups
		groups = append(groups, Group{Name: groupName, Enabled: enabled, Active: false, Hosts: hosts})
	}
	//refresh config file content
	h.persistConfig()
	return groups
}

func (h *Manager) persistConfig() {
	jsonText, err := json.Marshal(h.Config)
	ErrorAndExitWithLog(err)
	err = ioutil.WriteFile(h.hostsDir+"/"+h.ConfigFileName, jsonText, 0666)
	ErrorAndExitWithLog(err)
}

func (h *Manager) addGroupToConfig(groupName string, enabled bool, lastUpdatedTimestamp uint32) {
	h.Config.Groups = append(h.Config.Groups, HostGroupData{
		Name:                 groupName,
		Enabled:              enabled,
		LastUpdatedTimestamp: lastUpdatedTimestamp,
	})
}
