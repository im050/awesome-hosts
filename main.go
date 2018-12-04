package main

import (
	"awesome-hosts/manager"
	"flag"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

func main() {
	//init manager instance
	m = manager.New(manager.GetUserHome() + "/.awesohosts").Init()
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
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			m.Window = ws[0]
			return nil
		},
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astilectron.PtrStr("#2d3e50"),
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
