package manager

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

type Host struct {
	Domain  string `json:"domain"`
	IP      string `json:"ip"`
	Enabled bool   `json:"enabled"`
}

type Hosts []Host //line=>number
type Groups []Group

type Group struct {
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	Active    bool   `json:"active"`
	Hosts     Hosts  `json:"hosts"`
	CreatedAt int64  `json:"createAt"`
}

type GroupConfig struct {
	Name                 string
	Enabled              bool
	LastUpdatedTimestamp int64
	CreatedTimestamp     int64
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
	TempFileName     string
	Window           *astilectron.Window
	SudoPassword     string
	GroupConfigIndex map[string]*GroupConfig
}

func New(hostsDir string) *Manager {
	m := new(Manager)
	m.hostsDir = hostsDir
	m.DefaultGroupName = "Default Hosts"
	m.ConfigFileName = "data.config"
	m.TempFileName = "hosts.temp"
	m.SudoPassword = ""
	return m
}

func (m *Manager) Init() *Manager {
	//fourth, every x ms sync system hosts
	defer m.syncStart()
	//third, init host groups
	defer m.initGroups()
	//second, init or load config into m.Config
	defer m.loadConfig()
	//first, backup system hosts as a new group
	defer m.initSystemHosts()
	exists, err := PathExists(m.hostsDir)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	if exists {
		return m
	}
	//create hosts dir
	err = os.Mkdir(m.hostsDir, 0777)
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	return m
}

func (m *Manager) initSystemHosts() {
	file, _ := os.Open(GetHostsFile())
	defer file.Close()
	hosts := m.GetHosts(file)
	m.SystemHosts = hosts
	exists, err := PathExists(m.GetGroupFilePath(m.DefaultGroupName))
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	if exists {
		return
	}
	m.WriteHosts(m.GetGroupFilePath(m.DefaultGroupName), m.SystemHosts)
}

func (m *Manager) GetGroupFilePath(groupName string) string {
	return m.GetHostDir() + "/" + groupName + ".host"
}

func (m *Manager) initGroups() {
	groups := m.GetGroups()
	m.Groups = groups
}

func (m *Manager) GetHostDir() string {
	return m.hostsDir
}

func (m *Manager) loadConfig() {
	exists, _ := PathExists(m.hostsDir + "/" + m.ConfigFileName)
	if exists {
		m.loadConfigFromFile()
	} else {
		m.Config = Config{
			Groups:               []GroupConfig{{Name: m.DefaultGroupName, Enabled: true, LastUpdatedTimestamp: 0, CreatedTimestamp: GetNowTimestamp()}},
			LastUpdatedTimestamp: 0,
			LastSyncTimestamp:    0,
		}
		m.persistConfig()
	}
}

func (m *Manager) loadConfigFromFile() {
	fileByte, err := ioutil.ReadFile(m.hostsDir + "/" + m.ConfigFileName)
	ErrorAndExitWithLog(err)
	ErrorAndExitWithLog(json.Unmarshal(fileByte, &m.Config))
}

func (m *Manager) GetHosts(file *os.File) Hosts {
	br := bufio.NewReader(file)
	//each file line by line
	var hosts []Host
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		host, ok := m.explainHostsLine(string(line))
		if !ok {
			continue
		}
		hosts = append(hosts, host)
	}
	return hosts
}

func (m *Manager) explainHostsLine(line string) (Host, bool) {
	lineString := strings.TrimSpace(line)
	//if empty, continue
	if len(lineString) == 0 {
		return Host{}, false
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
		return Host{}, false
	}
	if !enabled {
		hostSplit[0] = strings.TrimSpace(strings.TrimLeft(hostSplit[0], "#"))
	}
	IPv4Pattern := `((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`
	IPv6Pattern := `(([\da-fA-F]{1,4}):){8}`
	//if not ip, continue
	if !regexp.MustCompile(IPv4Pattern).MatchString(hostSplit[0]) && !regexp.MustCompile(IPv6Pattern).MatchString(hostSplit[0]) {
		return Host{}, false
	}
	return Host{
		Domain:  hostSplit[1],
		IP:      hostSplit[0],
		Enabled: enabled,
	}, true
}

func (m *Manager) explainHostsString(content string) Hosts {
	var hosts []Host
	lines := strings.Split(content, GetLineSeparator())
	for _, line := range lines {
		host, ok := m.explainHostsLine(line)
		if !ok {
			continue
		}
		hosts = append(hosts, host)
	}
	return hosts
}

func (m *Manager) WriteContent(filename string, content string) error {
	data := []byte(content)
	err := ioutil.WriteFile(filename, data, 0666)
	return err
}

func (m *Manager) WriteHosts(name string, hosts Hosts) error {
	hostsContent := ""
	eol := GetLineSeparator()
	for _, host := range hosts {
		if !host.Enabled {
			hostsContent += "#"
		}
		hostsContent += host.IP + " " + host.Domain + eol
	}
	return m.WriteContent(name, hostsContent)
}

func (m *Manager) initGroupConfigIndex() {
	m.GroupConfigIndex = make(map[string]*GroupConfig)
	for i, _ := range m.Config.Groups {
		c := &m.Config.Groups[i];
		m.GroupConfigIndex[c.Name] = c
	}
}

func (m *Manager) GetGroups() []Group {
	var groups []Group
	m.initGroupConfigIndex()
	files, _ := ioutil.ReadDir(m.hostsDir)
	timestamp := GetNowTimestamp()
	for _, f := range files {
		groupInfo := strings.Split(f.Name(), ".")
		if groupInfo[len(groupInfo)-1] != "host" {
			continue
		}
		groupName := groupInfo[0]
		if len(groupInfo) >= 3 {
			groupName = strings.Join(groupInfo[0:], ".")
		}
		var enabled bool
		var createdTimestamp int64
		gc := m.FindGroupConfig(groupName)
		if gc == nil {
			enabled = false
			createdTimestamp = timestamp
			//if group doesn't exists in config file, then add it.
			m.addGroupToConfig(groupName, enabled, 0, createdTimestamp)
		} else {
			enabled = gc.Enabled
			createdTimestamp = gc.CreatedTimestamp
		}
		//read host file
		file, err := os.Open(m.hostsDir + "/" + f.Name())
		hosts := m.GetHosts(file)
		ErrorAndExitWithLog(file.Close())
		ErrorAndExitWithLog(err)
		//append to groups
		groups = append(groups, Group{Name: groupName, Enabled: enabled, Active: false, Hosts: hosts, CreatedAt: createdTimestamp})
	}
	fmt.Println("Before sort", groups)
	sort.SliceStable(groups, func(i, j int) bool {
		return groups[i].CreatedAt < groups[j].CreatedAt
	})
	fmt.Println("After sort", groups)

	return groups
}

//Add host to Group
func (m *Manager) AddHost(groupName string, host Host) bool {
	group := m.FindGroup(groupName)
	//when found group
	if group == nil {
		return false
	}
	group.Hosts = append(group.Hosts, host)
	//refresh groups config order by name of group
	m.refreshGroupsConfig(group.Name)
	return true
}

func (m *Manager) UpdateHost(groupName string, index int, ip string, domain string, enabled bool) bool {
	group := m.FindGroup(groupName)
	if (index + 1) > len(group.Hosts) {
		return false
	}
	host := &group.Hosts[index]
	host.Enabled = enabled
	host.IP = ip
	host.Domain = domain
	//refresh groups config order by name of group
	m.refreshGroupsConfig(group.Name)
	return true
}

//find group with group name and return a pointer
func (m *Manager) FindGroup(groupName string) *Group {
	index := -1
	for i, _ := range m.Groups {
		group := &m.Groups[i]
		if group.Name != groupName {
			continue
		}
		index = i
	}
	if index == -1 {
		return nil
	}
	return &m.Groups[index]
}

func (m *Manager) EnableGroup(groupName string, enabled bool) bool {
	config := m.FindGroupConfig(groupName)
	if config == nil {
		return false
	}
	group := m.FindGroup(groupName)
	if group == nil {
		return false
	}
	group.Enabled = enabled
	config.Enabled = enabled
	m.Config.LastUpdatedTimestamp = GetNowTimestamp()
	return true
}

func (m *Manager) syncStart() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for _ = range ticker.C {
			if m.Config.LastUpdatedTimestamp-m.Config.LastSyncTimestamp <= 0 {
				//fmt.Println("continue")
				continue
			}
			m.Config.LastSyncTimestamp = GetNowTimestamp()
			m.persistConfig()
			tmpHosts := ""
			for _, config := range m.Config.Groups {
				group := m.FindGroup(config.Name)
				if group == nil {
					continue
				}
				str := m.persistGroup(group)
				if config.Enabled {
					tmpHosts += "#Group Name: " + config.Name + GetLineSeparator()
					tmpHosts += str + GetLineSeparator() + GetLineSeparator()
				}
			}
			err := ioutil.WriteFile(m.hostsDir+"/"+m.TempFileName, []byte(tmpHosts), 0666)
			ErrorAndExitWithLog(err)
			m.SyncSystemHosts()
			fmt.Println("updated")
		}
	}()
}

//refresh config
//when group has changed, remember call this func to updated config file and var `m.Config`.
func (m *Manager) refreshGroupsConfig(groupName string) {
	timestamp := GetNowTimestamp()
	config := m.FindGroupConfig(groupName)
	if config == nil {
		return;
	}
	m.Config.LastUpdatedTimestamp = timestamp
	config.LastUpdatedTimestamp = timestamp
}

//find host group data
func (m *Manager) FindGroupConfig(groupName string) *GroupConfig {
	//for i, _ := range m.Config.Groups {
	//	config := &m.Config.Groups[i]
	//	if config.Name == groupName {
	//		return config
	//	}
	//}
	v, ok := m.GroupConfigIndex[groupName]
	if !ok {
		return nil
	}
	return v
}

func (m *Manager) persistGroup(group *Group) string {
	groupName := group.Name
	filePath := m.GetGroupFilePath(groupName)
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
	//fmt.Println("write here", str)
	err := ioutil.WriteFile(filePath, []byte(str), 0666)
	ErrorAndExitWithLog(err)
	return str
}

func (m *Manager) persistConfig() {
	jsonText, err := json.Marshal(m.Config)
	ErrorAndExitWithLog(err)
	err = ioutil.WriteFile(m.hostsDir+"/"+m.ConfigFileName, jsonText, 0666)
	ErrorAndExitWithLog(err)
}

func (m *Manager) addGroupToConfig(groupName string, enabled bool, lastUpdatedTimestamp int64, createdTimestamp int64) {
	m.Config.Groups = append(m.Config.Groups, GroupConfig{
		Name:                 groupName,
		Enabled:              enabled,
		LastUpdatedTimestamp: lastUpdatedTimestamp,
		CreatedTimestamp:     createdTimestamp,
	})
	m.GroupConfigIndex[groupName] = &m.Config.Groups[len(m.Config.Groups) - 1]
}

func (m *Manager) SyncSystemHosts() bool {
	if runtime.GOOS == "windows" {
		return m.SyncSystemHostsWin()
	} else {
		return m.SyncSystemHostsUnix()
	}
}

func (m *Manager) SyncSystemHostsWin() bool {
	return true
}

func (m *Manager) AddGroup(name string, enabled bool, hosts string) bool {
	timestamp := GetNowTimestamp()
	m.addGroupToConfig(name, enabled, timestamp, timestamp)
	group := Group{Name: name, Enabled: enabled, Hosts: m.explainHostsString(hosts)}
	m.Groups = append(m.Groups, group)
	m.Config.LastUpdatedTimestamp = timestamp
	return true
}

//
//'Permission denied'
//    , 'incorrect password'
//    , 'Password:Sorry, try again.'
func (m *Manager) SyncSystemHostsUnix() bool {
	var (
		output string
		err    error
	)
	if m.SudoPassword != "" {
		commands := []string{"echo '" + m.SudoPassword + "' | sudo -S chmod 777 " + GetHostsFile(),
			"cat " + m.hostsDir + "/" + m.TempFileName + " > " + GetHostsFile(),
			"echo '" + m.SudoPassword + "' | sudo -S chmod 644 " + GetHostsFile()}
		command := strings.Join(commands, " && ")
		output, err = ShellCommand(command)
	} else {
		command := "cat " + m.hostsDir + "/" + m.TempFileName + " > " + GetHostsFile()
		output, err = ShellCommand(command)
	}
	needPassString := [3]string{"Permission denied", "incorrect password", "Password:Sorry, try again."}
	if err != nil {
		isNeedPass := false
		for _, str := range needPassString {
			if strings.Index(output, str) != -1 {
				isNeedPass = true
				break
			}
		}
		if isNeedPass {
			bootstrap.SendMessage(m.Window, "needPassword", "syncSystemHostsUnix")
		}
		return false
	}
	return true
}
