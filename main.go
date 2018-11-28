package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
	"host-manager/manager"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	//init manager instance
	m = manager.New(getCurrentDirectory())
	//open hosts file
	file, err := os.Open(manager.GetHostsFile())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	//load hosts
	go func() {
		hosts = m.GetHosts(file)
		m.WriteHosts("默认文件", hosts)
	}()
	// Init
	flag.Parse()
	astilog.FlagInit()
	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug:         *debug,
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astilectron.PtrStr("#333"),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(650),
				Width:           astilectron.PtrInt(950),
				MinHeight:       astilectron.PtrInt(650),
				MinWidth:        astilectron.PtrInt(950),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "event.name":
		// Unmarshal payload
		var s string
		if err = json.Unmarshal(m.Payload, &s); err != nil {
			payload = err.Error()
			return
		}
		payload = s + " world"
	case "list":
		payload = hosts
	}
	return
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
