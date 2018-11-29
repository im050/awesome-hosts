package main

import (
	"encoding/json"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)

func handleMessages(w *astilectron.Window, mi bootstrap.MessageIn) (payload interface{}, err error) {
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
		var s interface{}
		if err = json.Unmarshal(mi.Payload, &s); err != nil {
			payload = err.Error()
			panic(err)
			return
		}
		payload = m.SystemHosts
	case "groups":
		payload = m.GetGroups()
	}

	return
}
