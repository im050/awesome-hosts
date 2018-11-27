package main

import "net"

// `json:"xxx"` 转换json时相对应的字段
type Host struct {
	Domain  string `json:"domain"`
	IP      string `json:"ip"`
	Enabled bool   `json:"enabled"`
}

var (
	hosts []Host
	ln    net.Listener
)
