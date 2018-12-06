package manager

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
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
	InstalledTimestamp   int64
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
	//third, sync your changed with an interval
	defer m.syncStart()
	//second, init host groups
	defer m.initGroups()
	//first, init or load config into m.Config
	defer m.loadConfig()
	//init current system hosts
	m.initSystemHosts()
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
	file, err := os.Open(GetHostsFile())
	m.SystemHosts = m.GetHosts(file)
	defer file.Close()
	if err != nil {
		panic(err)
	}
}

func (m *Manager) backupSystemHosts() {
	if err := m.WriteHosts(m.GetGroupFilePath(m.DefaultGroupName), m.SystemHosts); err != nil {
		panic(err)
	}
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
			InstalledTimestamp:   0,
		}
	}
	if m.Config.InstalledTimestamp <= 0 {
		//update install timestamp
		m.Config.InstalledTimestamp = GetNowTimestamp()
		//first install, backup your system hosts
		m.backupSystemHosts()
		//and save the config
		m.persistConfig()
	}
	//create index for group config
	m.initGroupConfigIndex()
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
	if m.CheckIP(hostSplit[0]) != nil {
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
	lines := strings.Split(content, "\n") //because of textarea, so the line separator is "\n"
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
		c := &m.Config.Groups[i]
		m.GroupConfigIndex[c.Name] = c
	}
}

func (m *Manager) GetGroups() []Group {
	var groups []Group
	for _, item := range m.Config.Groups {
		groupFileName := m.GetGroupFilePath(item.Name)
		file, err := os.OpenFile(groupFileName, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			continue
		}
		hosts := m.GetHosts(file)
		groups = append(groups, Group{Name: item.Name, Enabled: item.Enabled, CreatedAt: item.CreatedTimestamp, Hosts: hosts})
		if err := file.Close(); err != nil {
			panic(err)
		}
	}
	//sort by createdTimestamp
	sort.SliceStable(groups, func(i, j int) bool {
		return groups[i].CreatedAt < groups[j].CreatedAt
	})
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
			//nothing changed if updated time less then sync time
			if m.Config.LastUpdatedTimestamp-m.Config.LastSyncTimestamp <= 0 {
				continue
			}
			m.Config.LastSyncTimestamp = GetNowTimestamp()
			m.persistConfig()
			tmpHosts := ""
			//build hosts
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
			//check system hosts file content is the same as tmpHosts
			if !m.needSync(tmpHosts) {
				continue
			}
			err := m.WriteContent(m.GetHostDir()+"/"+m.TempFileName, tmpHosts)
			ErrorAndExitWithLog(err)
			//sync hosts to host file
			m.SyncSystemHosts()
			//latest hosts list
			m.SystemHosts = m.explainHostsString(tmpHosts)
			//refresh client system hosts list
			m.SendMessage("updateSystemHosts", m.SystemHosts)
		}
	}()
}

func (m *Manager) needSync(tmpHosts string) bool {
	systemHostFile, err := os.Open(GetHostsFile())
	defer systemHostFile.Close()
	ErrorAndExitWithLog(err)
	systemHostMD5 := md5.New()
	_, ioErr := io.Copy(systemHostMD5, systemHostFile)
	ErrorAndExitWithLog(ioErr)
	systemHostMD5String := hex.EncodeToString(systemHostMD5.Sum(nil))
	tempHostMD5 := md5.New()
	tempHostMD5.Write([]byte(tmpHosts))
	tempHostMD5String := hex.EncodeToString(tempHostMD5.Sum(nil))
	if systemHostMD5String == tempHostMD5String {
		return false
	}
	return true
}

//refresh config
//when group has changed, remember call this func to updated config file and var `m.Config`.
func (m *Manager) refreshGroupsConfig(groupName string) {
	timestamp := GetNowTimestamp()
	config := m.FindGroupConfig(groupName)
	if config == nil {
		return
	}
	m.Config.LastUpdatedTimestamp = timestamp
	config.LastUpdatedTimestamp = timestamp
}

//find host group data
func (m *Manager) FindGroupConfig(groupName string) *GroupConfig {
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
	m.GroupConfigIndex[groupName] = &m.Config.Groups[len(m.Config.Groups)-1]
}

func (m *Manager) SyncSystemHosts() bool {
	if runtime.GOOS == "windows" {
		return m.SyncSystemHostsWin()
	} else {
		return m.SyncSystemHostsUnix()
	}
}

func (m *Manager) SyncSystemHostsWin() bool {
	file, err := os.Open(m.GetHostDir() + "/" + m.TempFileName)
	if err != nil {
		return false
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}
	err = ioutil.WriteFile(GetHostsFile(), content, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
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
			m.SendMessage("needPassword", "syncSystemHostsUnix")
		}
		return false
	}
	return true
}

func (m *Manager) SendMessage(name string, payload interface{}) bool {
	if err := bootstrap.SendMessage(m.Window, name, payload); err != nil {
		return false
	}
	return true
}

func (m *Manager) ChangeGroupName(oldName string, newName string) {
	//update group config
	groupConfig := m.FindGroupConfig(oldName)
	groupConfig.Name = newName
	//update group
	group := m.FindGroup(oldName)
	group.Name = newName
	//update index
	delete(m.GroupConfigIndex, oldName)
	m.GroupConfigIndex[newName] = groupConfig
	//update file name
	ErrorAndExitWithLog(os.Rename(m.GetGroupFilePath(oldName), m.GetGroupFilePath(newName)))
	m.persistConfig()
}

func (m *Manager) DeleteGroup(groupName string) {
	m.deleteGroupWithGroupName(groupName).deleteGroupConfigWithGroupName(groupName).initGroupConfigIndex()
	if err := os.Remove(m.GetGroupFilePath(groupName)); err != nil {
		panic(err)
	}
	m.Config.LastUpdatedTimestamp = GetNowTimestamp()
}

func (m *Manager) deleteGroupConfigWithGroupName(groupName string) *Manager {
	index := -1
	for i, _ := range m.Config.Groups {
		group := m.Config.Groups[i]
		if groupName == group.Name {
			index = i
			break
		}
	}
	if index == -1 {
		return m
	}
	m.Config.Groups = append(m.Config.Groups[:index], m.Config.Groups[index+1:]...)
	return m
}

func (m *Manager) deleteGroupWithGroupName(groupName string) *Manager {
	index := -1
	for i, _ := range m.Groups {
		group := m.Groups[i]
		if groupName == group.Name {
			index = i
			break
		}
	}
	if index == -1 {
		return m
	}
	m.Groups = append(m.Groups[:index], m.Groups[index+1:]...)
	return m
}

func (m *Manager) DeleteHost(groupName string, index int) {
	group := m.FindGroup(groupName)
	if group == nil {
		return
	}
	group.Hosts = append(group.Hosts[:index], group.Hosts[index+1:]...)
	m.Config.LastUpdatedTimestamp = GetNowTimestamp()
}

func (m *Manager) DeleteHostsByGroup(group *Group, indexes []int) {
	indexes = RemoveRepeatNumber(indexes)
	sort.SliceStable(indexes, func(i, j int) bool {
		return indexes[i] < indexes[j]
	})
	for i, index := range indexes {
		index -= i
		group.Hosts = append(group.Hosts[:index], group.Hosts[index+1:]...)
	}
	m.Config.LastUpdatedTimestamp = GetNowTimestamp()
}

func (m *Manager) CheckIP(ip string) error {
	IPv4Pattern := `((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`
	IPv6Pattern := `([a-f0-9]{1,4}(:[a-f0-9]{1,4}){7}|[a-f0-9]{1,4}(:[a-f0-9]{1,4}){0,7}::[a-f0-9]{0,4}(:[a-f0-9]{1,4}){0,7})`
	if !regexp.MustCompile(IPv4Pattern).MatchString(ip) && !regexp.MustCompile(IPv6Pattern).MatchString(ip) && ip != "::1" {
		return fmt.Errorf("IP [%s] is illegal. ", ip)
	}
	return nil
}

func (m *Manager) CheckDomain(domain string) error {
	pattern := `^[^\.]([A-Za-z\.\-\_0-9]+)[^.]$`
	if !regexp.MustCompile(pattern).MatchString(domain) {
		return fmt.Errorf("Domain [%s] is illegal. ", domain)
	}
	return nil

}

func (m *Manager) CheckGroupName(name string) error {
	pattern := `[\\\\/:*?\"<>|]`
	if regexp.MustCompile(pattern).MatchString(name) {
		return fmt.Errorf("Group name [%s] is illegal. ", name)
	}
	return nil
}
