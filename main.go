package main

import (
	"embed"
	"github.com/fafeitsch/private-running-journal/backend"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := backend.NewApp()

	title := "Private Running Journal"
	if app.Language() == "de" {
		title = "Privates Lauftagebuch"
	}

	// Create application with options
	err := wails.Run(
		&options.App{
			Title:  title,
			Width:  1024,
			Height: 768,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
			BackgroundColour: options.NewRGB(uint8(255), uint8(255), uint8(255)),
			OnStartup:        app.Startup,
			Bind: []interface{}{
				app,
			},
		},
	)

	if err != nil {
		println("Error:", err.Error())
	}
}
