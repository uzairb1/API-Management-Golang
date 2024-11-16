package routes

import (
	"STACKIT/controllers"
	"net/http"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes() {
	// Create an instance of the GuestController
	guestController := controllers.NewGuestController()

	// Use the methods from the GuestController instance
	http.HandleFunc("/register", guestController.RegisterGuest)
	http.HandleFunc("/guests", guestController.ListGuests)
	http.HandleFunc("/guests/count", guestController.CountGuests)
	http.HandleFunc("/guests/search", guestController.SearchGuests)
}
