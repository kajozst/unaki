package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"log"
	"time"

	"github.com/kajozst/unaki/nak"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {

	app := application.New(application.Options{
		Name:        "unaki",
		Description: "A demo of using raw HTML & CSS",
		Services:    []application.Service{},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})
	app.OnEvent("command", func(e *application.CustomEvent) {
		app.Logger.Info("---EVENT RECEIVED !!!!!!")
		var args []string
		err := json.Unmarshal([]byte(e.Data.(string)), &args)
		if err != nil {
			// Handle the error when unmarshaling the JSON data
			app.Logger.Info("---NOT-OK")
			return
		}

		output, err := nak.Nak(args)
		if err != nil {
			// Handle the error returned by the Nak() function
			app.Logger.Info("---ERR IS NOT NIL !!!")

			return
		}

		//app.Logger.Info("output type = %T\n", output)
		app.Logger.Info(string(output))
		app.EmitEvent("commandRes", string(output))

	})

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.EmitEvent("time", now)
			app.Logger.Info("TICK")
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
