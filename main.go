package main

import (
	"STACKIT/controllers"
	"STACKIT/routes"
	"log"
	"net/http"
)

func main() {
	// Initialize GuestController and storage
	guestController := controllers.NewGuestController()
	guestController.InitializeStorage()

	// Setup routes using the controller instance
	routes.SetupRoutes()

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
