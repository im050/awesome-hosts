package main

import (
	"awesome-hosts/parameters"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParams(t *testing.T) {
	p := parameters.New()
	data := map[string]interface{}{
		"world": "world",
		"bool" : true,
	}
	p.From(data)
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
