package main

import (
	"encoding/json"
	"gowebrtc/app/route"
	"gowebrtc/app/shared/database"
	"gowebrtc/app/shared/email"
	"gowebrtc/app/shared/jsonconfig"
	"gowebrtc/app/shared/recaptcha"
	"gowebrtc/app/shared/server"
	"gowebrtc/app/shared/session"
	"gowebrtc/app/shared/view"
	"gowebrtc/app/shared/view/plugin"
	"log"
	"os"
	"runtime"
	// "gowebrtc/server"
)

// configuration contains the application settings
type configuration struct {
	Database  database.Info   `json:"Database"`
	Email     email.SMTPInfo  `json:"Email"`
	Recaptcha recaptcha.Info  `json:"Recaptcha"`
	Server    server.Server   `json:"Server"`
	Session   session.Session `json:"Session"`
	Template  view.Template   `json:"Template"`
	View      view.View       `json:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// config the settings variable
var config = &configuration{}

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// server.Run()

	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)

	// Configure the session cookie store
	session.Configure(config.Session)

	// Connect to database
	database.Connect(config.Database)

	// Configure the Google reCAPTCHA prior to loading view plugins
	recaptcha.Configure(config.Recaptcha)

	// Setup the views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		recaptcha.Plugin())

	// Start the listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)

	// sigs := make(chan os.Signal, 1)
	// done := make(chan bool, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	// 	sig := <-sigs
	// 	log.Println(sig)
	// 	done <- true
	// }()

	// log.Println("Server Start Awaiting Signal")
	// <-done
	// log.Println("Exiting")
}
