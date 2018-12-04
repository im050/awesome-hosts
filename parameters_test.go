package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParams(t *testing.T) {
	p := NewParameters()
	data := map[string]interface{}{
		"world": "world",
		"bool" : true,
	}
	p.From(data)
	a := p.GetInt("testInt", 5)
	s := p.GetString("hello", "hello")
	assert.Equal(t, a, 5)
	assert.Equal(t, s, "hello")
	assert.Equal(t, p.GetString("world", "memory"), "world")
	assert.Equal(t, true, p.GetBool("bool", false))
	assert.Equal(t, false, p.GetBool("bool2", false))
	assert.Equal(t, true, p.GetBool("bool3", true))
}
