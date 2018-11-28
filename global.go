package main

import (
	"flag"
	"github.com/asticode/go-astilectron"
	"host-manager/manager"
)

var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
	m       *manager.Manager
	hosts   []manager.Host
)
