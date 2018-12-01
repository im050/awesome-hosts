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

func handleMessages(w *astilectron.Window, mi bootstrap.MessageIn) (payload interface{}, err error) {
	var data map[string]string
	data = make(map[string]string, 0)
	if err = json.Unmarshal(mi.Payload, &data); err != nil {
		payload = err.Error()
		fmt.Println(err)
		return
	}
	switch mi.Name {
	case "event.name":
		// Unmarshal payload
		var s string
		if err = json.Unmarshal(mi.Payload, &s); err != nil {
			payload = err.Error()
			return
		}
		fmt.Println(s, mi.Payload, 111)
		payload = s + " world"
	case "list":
		payload = m.SystemHosts
	case "groups":
		payload = m.Groups
	case "intranet":
		payload = manager.GetIntranetIp()
	case "addHost":
		fmt.Println(data)
		m.AddHost(data["groupName"], manager.Host{IP: data["ip"], Domain: data["domain"], Enabled: true})
		payload = ElectronResponse(1, "success", nil)
		default:
		payload = "not found"
	}

	return
}

func ElectronResponse(code int, message string, payload interface{}) Response {
	return Response{Code: code, Message: message, Payload: payload}
}