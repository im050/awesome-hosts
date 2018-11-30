package main

import (
	"awesome-hosts/manager"
	"flag"
	"github.com/asticode/go-astilectron"
)

var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
	m       *manager.Manager
)
