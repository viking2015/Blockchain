package main

import (
	"log"

	"github.com/NlaakStudios/gowaf/app"
	// Import Application Models and Controllers here
)

// Register holds all the Models and Controllers you wish to register for the WebApp
func Register(a *app.App) {
	// Register Your WebApp Specific Models below
	a.Model.Register(
	//&models.MyModel{},
	)
}

func main() {
	// Create a new MVC Application object
	app, err := app.NewMVC(Version(), "./data", "webapp")
	if err != nil {
		log.Fatal(err)
	} else {
		app.Register()           // Register Core Framework Models & Controllers
		app.Run(Register, false) // Register WebApp Models & Controllers and Run
	}
}
