package main

import (
	"awesome-hosts/manager"
	"encoding/json"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func handleMessages(w *astilectron.Window, mi bootstrap.MessageIn) (payload interface{}, handleErr error) {
	var data map[string]interface{}
	data = make(map[string]interface{}, 0)
	if err := json.Unmarshal(mi.Payload, &data); err != nil {
		payload = ElectronResponse(0, err.Error(), nil)
		fmt.Println(err)
		return
	}
	switch mi.Name {
	case "list":
		payload = ElectronResponse(1, "success", m.SystemHosts)
	case "groups":
		payload = ElectronResponse(1, "success", m.Groups)
	case "intranet":
		payload = ElectronResponse(1, "success", manager.GetIntranetIp())
	case "addHost":
		ip := data["ip"].(string)
		domain := data["domain"].(string)
		if err := m.CheckIP(ip); err != nil {
			payload = ElectronResponse(0, err.Error(), nil)
			return
		}
		if err := m.CheckDomain(domain); err != nil {
			payload = ElectronResponse(0, err.Error(), nil)
			return
		}
		m.AddHost(data["groupName"].(string), manager.Host{IP: ip, Domain: domain, Enabled: true})
		payload = ElectronResponse(1, "success", nil)
	case "updateHost":
		ip := data["ip"].(string)
		domain := data["domain"].(string)
		if err := m.CheckIP(ip); err != nil {
			payload = ElectronResponse(0, err.Error(), nil)
			return
		}
		if err := m.CheckDomain(domain); err != nil {
			payload = ElectronResponse(0, err.Error(), nil)
			return
		}
		if m.UpdateHost(data["groupName"].(string), int(data["index"].(float64)), ip, domain, data["enabled"].(bool)) {
			payload = ElectronResponse(1, "success", nil)
		} else {
			payload = ElectronResponse(0, "An error occurred while an operation", nil)
		}
	case "enableGroup":
		if m.EnableGroup(data["groupName"].(string), data["enabled"].(bool)) {
			payload = ElectronResponse(1, "success", nil)
		} else {
			payload = ElectronResponse(0, "An error occurred while an operation", nil)
		}
	case "syncSystemHostsUnix":
		m.SudoPassword = data["password"].(string)
		if m.SyncSystemHostsUnix() {
			payload = ElectronResponse(1, "success", nil)
		} else {
			payload = ElectronResponse(0, "An error occurred while an operation", nil)
		}
	case "addGroup":
		groupName := data["name"].(string)
		if err := m.CheckGroupName(groupName); err != nil {
			payload = ElectronResponse(0, err.Error(), nil)
			return
		}
		if m.FindGroup(groupName) != nil {
			payload = ElectronResponse(0, "Group already exists", nil)
			return
		}
		if m.AddGroup(groupName, data["enabled"].(bool), data["hosts"].(string)) {
			payload = ElectronResponse(1, "success", m.Groups)
		} else {
			payload = ElectronResponse(0, "An error occurred while an operation", nil)
		}
	case "changeGroup":
		groupName := data["newName"].(string)
		if err := m.CheckGroupName(groupName); err != nil {
			payload = ElectronResponse(0, err.Error(), nil)
			return
		}
		if m.FindGroupConfig(groupName) != nil {
			payload = ElectronResponse(0, "Group already exists", nil)
			return
		}
		oldName := data["oldName"].(string)
		m.ChangeGroupName(oldName, groupName)
		payload = ElectronResponse(1, "success", m.Groups)
	case "deleteGroup":
		if m.FindGroupConfig(data["groupName"].(string)) == nil {
			payload = ElectronResponse(0, "Group not exists", nil)
			return
		}
		m.DeleteGroup(data["groupName"].(string))
		payload = ElectronResponse(1, "success", m.Groups)
	case "deleteHost":
		index := int(data["index"].(float64))
		groupName := data["groupName"].(string)
		m.DeleteHost(groupName, index)
		payload = ElectronResponse(1, "success", nil)
	default:
		payload = ElectronResponse(404, "Not Found", nil)
	}

	return
}

//func GetParams(data map[string]interface{}, name string) interface{} {
//	v, ok := data[name]
//	if !ok {
//		return nil
//	}
//	return v
//}

func ElectronResponse(code int, message string, payload interface{}) Response {
	return Response{Code: code, Message: message, Payload: payload}
}