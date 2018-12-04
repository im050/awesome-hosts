package main

import (
	"awesome-hosts/parameters"
	"encoding/json"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"reflect"
)

type Handler struct {
	Parameters *parameters.Parameters
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func (handler *Handler) toFirstCharUpper(str string) string {
	r := []rune(str)
	if r[0] >= 97 && r[0] <= 122 {
		r[0] -= 32
	}
	return string(r)
}

func (handler *Handler) handleMessages(w *astilectron.Window, messageIn bootstrap.MessageIn) (payload interface{}, handleErr error) {
	//explain data
	var data map[string]interface{}
	data = make(map[string]interface{})
	if err := json.Unmarshal(messageIn.Payload, &data); err != nil {
		payload = nil
		return
	}
	//set parameters
	handler.Parameters.From(data)
	//call func
	reflectVal := reflect.ValueOf(handler)
	method := reflectVal.MethodByName(h.toFirstCharUpper(messageIn.Name) + "Handler")
	if method.IsValid() {
		retVal := method.Call(nil)
		return retVal[0].Interface().(Response), nil
	}
	return Response{Code: 0, Message: "Not Found"}, nil
}