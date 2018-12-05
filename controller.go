package main

import (
	"awesome-hosts/manager"
	"strings"
)

func (handler *Handler) ListHandler() Response {
	return Response{1, "success", m.SystemHosts}
}

func (handler *Handler) GroupsHandler() Response {
	return Response{1, "success", m.Groups}
}

func (handler *Handler) IntranetHandler() Response {
	return Response{1, "success", manager.GetIntranetIp()}
}

func (handler *Handler) AddHostHandler() Response {
	ip, _ := handler.Parameters.GetString("ip", "")
	ip = strings.Trim(ip, " ")
	domain, _ := handler.Parameters.GetString("domain", "")
	domain = strings.Trim(domain, " ")
	if ip == "" || domain == "" {
		return Response{Code: 0, Message: "IP or Domain cannot be empty"}
	}
	if err := m.CheckIP(ip); err != nil {
		return Response{Code: 0, Message: err.Error()}
	}
	if err := m.CheckDomain(domain); err != nil {
		return Response{Code: 0, Message: err.Error()}
	}
	groupName, exists := handler.Parameters.GetString("groupName", "")
	if groupName == "" || !exists {
		return Response{Code: 0, Message: "Group name cannot be empty"}
	}
	m.AddHost(groupName, manager.Host{IP: ip, Domain: domain, Enabled: true})
	return Response{Code: 1, Message: "success"}
}

func (handler *Handler) UpdateHostHandler() Response {
	ip, _ := handler.Parameters.GetString("ip", "")
	domain, _ := handler.Parameters.GetString("domain", "")
	ip = strings.Trim(ip, " ")
	domain = strings.Trim(domain, " ")
	enabled, _ := handler.Parameters.GetBool("enabled", false)
	groupName, _ := handler.Parameters.GetString("groupName", "")
	index, indexExists := handler.Parameters.GetInt("index")
	group := m.FindGroup(groupName)
	if group == nil || index > len(group.Hosts) - 1 {
		return Response{Code:0, Message: "An error occurred while an operation"}
	}
	host := group.Hosts[index]
	if err := m.CheckIP(ip); err != nil {
		return Response{Code: 0, Message: err.Error(), Payload: host}
	}
	if err := m.CheckDomain(domain); err != nil {
		return Response{Code: 0, Message: err.Error(), Payload: host}
	}
	if !indexExists {
		return Response{Code: 0, Message: "An error occurred while an operation", Payload: host}
	}
	if host.IP == ip && host.Domain == domain && host.Enabled == enabled {
		return Response{Code: 1, Message: "success, nothing changed."}
	}
	if m.UpdateHost(groupName, index, ip, domain, enabled) {
		return Response{Code: 1, Message: "success"}
	} else {
		return Response{Code: 0, Message: "An error occurred while an operation", Payload: host}
	}
}

func (handler *Handler) EnableGroupHandler() Response {
	groupName, _ := handler.Parameters.GetString("groupName", "")
	enabled, _ := handler.Parameters.GetBool("enabled", false)
	if m.EnableGroup(groupName, enabled) {
		return Response{Code: 1, Message: "success"}
	} else {
		return Response{Code: 0, Message: "An error occurred while an operation"}
	}
}

func (handler *Handler) SyncSystemHostsUnixHandler() Response {
	m.SudoPassword, _ = handler.Parameters.GetString("password", "")
	if m.SudoPassword == "" {
		return Response{Code: 0, Message: "Password cannot be empty"}
	}
	if m.SyncSystemHostsUnix() {
		return Response{Code: 1, Message: "success"}
	} else {
		return Response{Code: 0, Message: "An error occurred while an operation"}
	}
}

func (handler *Handler) AddGroupHandler() Response {
	groupName, _ := handler.Parameters.GetString("name", "")
	enabled, _ := handler.Parameters.GetBool("enabled", false)
	hosts, _ := handler.Parameters.GetString("hosts", "")
	if groupName == "" {
		return Response{Code: 0, Message: "Group name cannot be empty"}
	}
	if err := m.CheckGroupName(groupName); err != nil {
		return Response{Code: 0, Message: err.Error()}
	}
	if m.FindGroup(groupName) != nil {
		return Response{Code: 0, Message: "Group already exists"}
	}
	if m.AddGroup(groupName, enabled, hosts) {
		return Response{Code: 1, Message: "success", Payload: m.Groups}
	} else {
		return Response{Code: 1, Message: "An error occurred while an operation"}
	}
}

func (handler *Handler) ChangeGroupHandler() Response {
	newName, _ := handler.Parameters.GetString("newName", "")
	oldName, _ := handler.Parameters.GetString("oldName", "")
	if newName == oldName {
		return Response{Code: 0, Message: "Nothing changed"}
	}
	if newName == "" || oldName == "" {
		return Response{Code: 0, Message: "Lost some parameters"}
	}
	if err := m.CheckGroupName(newName); err != nil {
		return Response{Code: 0, Message: err.Error()}
	}
	if m.FindGroupConfig(newName) != nil {
		return Response{Code: 0, Message: "Group already exists"}
	}

	m.ChangeGroupName(oldName, newName)
	return Response{Code: 1, Message: "success"}
}

func (handler *Handler) DeleteGroupHandler() Response {
	groupName, _ := handler.Parameters.GetString("groupName", "")
	if groupName == "" {
		return Response{Code: 0, Message: "Group name cannot be empty"}
	}
	if m.FindGroupConfig(groupName) == nil {
		return Response{Code: 0, Message: "Group already exists"}
	}
	m.DeleteGroup(groupName)
	return Response{Code: 1, Message: "success", Payload: m.Groups}
}

func (handler *Handler) DeleteHostHandler() Response {
	index, _ := handler.Parameters.GetInt("index", -1)
	groupName, _ := handler.Parameters.GetString("groupName", "")
	if index == -1 {
		return Response{Code: 0, Message: "Lost some parameters"}
	}
	m.DeleteHost(groupName, index)
	return Response{Code: 1, Message: "success"}
}
