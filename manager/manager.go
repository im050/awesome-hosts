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
	Domain  string `json:"domain"`
	IP      string `json:"ip"`
	Enabled bool   `json:"enabled"`
}

type Hosts []Host //line=>number
type Groups []Group

type Group struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Active  bool   `json:"active"`
	Hosts   Hosts  `json:"hosts"`
}

type GroupConfig struct {
	Name                 string
	Enabled              bool
	LastUpdatedTimestamp int64
}

type Config struct {
	Groups               []GroupConfig
	LastUpdatedTimestamp int64 //last timestamp of hosts data was updated
	LastSyncTimestamp    int64 //last timestamp of refresh hosts data to system hosts
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
			Groups:               []GroupConfig{{Name: h.DefaultGroupName, Enabled: true, LastUpdatedTimestamp: 0}},
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
	var hosts []Host
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
		hosts = append(hosts, Host{
			Domain:  hostSplit[1],
			IP:      hostSplit[0],
			Enabled: enabled,
		})
	}
	return hosts
}

func (h *Manager) WriteContent(name string, content string) {
	data := []byte(content)
	err := ioutil.WriteFile(h.hostsDir+"/"+name, data, 0666)
	if err != nil {
		_, _ = fmt.Println("Write to file failed.", err)
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

//Add host to Group
func (h *Manager) AddHost(groupName string, host Host) bool {
	group := h.FindGroup(groupName)
	//when found group
	if group == nil {
		return false
	}
	group.Hosts = append(group.Hosts, host)
	//save group data to file
	h.persistGroup(group)
	//refresh groups config order by name of group
	h.refreshGroupsConfig(group.Name)
	h.persistConfig()
	return true
}

func (h *Manager) UpdateHost(groupName string, index int, ip string, domain string, enabled bool) bool {
	group := h.FindGroup(groupName)
	if (index + 1) > len(group.Hosts) {
		return false
	}
	host := &group.Hosts[index]
	host.Enabled = enabled
	host.IP = ip
	host.Domain = domain
	//save group data to file
	h.persistGroup(group)
	//refresh groups config order by name of group
	h.refreshGroupsConfig(group.Name)
	h.persistConfig()
	return true
}

//find group with group name and return a pointer
func (h *Manager) FindGroup(groupName string) *Group {
	index := -1
	for i, _ := range h.Groups {
		group := &h.Groups[i]
		if group.Name != groupName {
			continue
		}
		index = i
	}
	if index == -1 {
		return nil
	}
	return &h.Groups[index]
}

func (h *Manager) EnableGroup(groupName string, enabled bool) bool {
	config := h.FindGroupConfig(groupName)
	if config == nil {
		return false
	}
	config.Enabled = enabled
	h.Config.LastUpdatedTimestamp = GetNowTimestamp()
	h.persistConfig()
	return true
}



//refresh config
//when group has changed, remember call this func to updated config file and var `h.Config`.
func (h *Manager) refreshGroupsConfig(groupName string) {
	timestamp := GetNowTimestamp()
	config := h.FindGroupConfig(groupName)
	if config == nil {
		return ;
	}
	h.Config.LastUpdatedTimestamp = timestamp
	config.LastUpdatedTimestamp = timestamp
}

//find host group data
func (h *Manager) FindGroupConfig(groupName string) *GroupConfig {
	for i, _ := range h.Config.Groups {
		config := &h.Config.Groups[i]
		if config.Name == groupName {
			return config
		}
	}
	return nil
}

func (h *Manager) persistGroup(group *Group) {
	groupName := group.Name
	filePath := h.hostsDir + "/" + GetHostFileName(groupName)
	str := ""
	for _, host := range group.Hosts {
		enabled := ""
		if !host.Enabled {
			enabled = "#"
		}
		str += enabled + host.IP + " " + host.Domain + GetLineSeparator()
	}
	//remove "\r\n" at last line
	str = strings.TrimRight(str, GetLineSeparator())
	fmt.Println("write here", str)
	err := ioutil.WriteFile(filePath, []byte(str), 0666)
	ErrorAndExitWithLog(err)
}

func (h *Manager) persistConfig() {
	jsonText, err := json.Marshal(h.Config)
	ErrorAndExitWithLog(err)
	err = ioutil.WriteFile(h.hostsDir+"/"+h.ConfigFileName, jsonText, 0666)
	ErrorAndExitWithLog(err)
}

func (h *Manager) addGroupToConfig(groupName string, enabled bool, lastUpdatedTimestamp int64) {
	h.Config.Groups = append(h.Config.Groups, GroupConfig{
		Name:                 groupName,
		Enabled:              enabled,
		LastUpdatedTimestamp: lastUpdatedTimestamp,
	})
}
