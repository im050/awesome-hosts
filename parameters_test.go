package main

import (
	"awesome-hosts/parameters"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParams(t *testing.T) {
	p := parameters.New()
	payload := "{\"world\": \"world\", \"bool\":true, \"number\": [1,2,3]}"
	var data map[string]interface{}
	data = make(map[string]interface{})
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		t.Error("err")
	}
	//data := map[string]interface{}{
	//	"world": "world",
	//	"bool" : true,
	//}
	p.From(data)
	number := p.Get("number", nil).([]interface{})
	fmt.Println(number)

	v1,_ := p.GetInt("testInt", 5)
	v2,_ := p.GetString("hello", "hello")
	v3, e1 := p.GetString("world", "memory")
	v4, _ := p.GetBool("bool", false)
	v5, _ := p.GetBool("bool2", false)
	v6, _ := p.GetBool("bool3", true)
	_, e2 := p.GetString("come")
	assert.Equal(t, 5, v1)
	assert.Equal(t, "hello", v2)
	assert.Equal(t, "world", v3)
	assert.Equal(t, true, e1)
	assert.Equal(t, true, v4)
	assert.Equal(t, false, v5)
	assert.Equal(t, true, v6)
	assert.Equal(t, false, e2)
}
