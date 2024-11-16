package tests

import (
	"STACKIT/controllers"
	"STACKIT/models"
	"STACKIT/storage"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper function to initialize the controller and clear storage for each test
func setup() controllers.GuestControllerInterface {
	storage.LoadGuests() // Make sure storage is loaded initially
	return controllers.NewGuestController()
}

func TestRegisterGuest(t *testing.T) {
	controller := setup()

	// Valid guest registration
	guest := models.Guest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}
	payload, _ := json.Marshal(guest)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()
	controller.RegisterGuest(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("expected %v, got %v", http.StatusCreated, status)
	}

	// Test invalid email
	invalidGuest := models.Guest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@invalid",
	}
	payload, _ = json.Marshal(invalidGuest)
	req, _ = http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	rr = httptest.NewRecorder()
	controller.RegisterGuest(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected %v, got %v", http.StatusBadRequest, status)
	}
}

func TestListGuests(t *testing.T) {
	controller := setup()

	req, _ := http.NewRequest("GET", "/guests?page=1&limit=10", nil)
	rr := httptest.NewRecorder()
	controller.ListGuests(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected %v, got %v", http.StatusOK, status)
	}

	var response map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("could not parse response: %v", err)
	}

	if _, ok := response["guests"]; !ok {
		t.Error("expected guests field in response")
	}
}

func TestCountGuests(t *testing.T) {
	controller := setup()

	// Simulate adding a few guests
	storage.SaveGuest(models.Guest{FirstName: "Alice", LastName: "Smith", Email: "alice@example.com"})
	storage.SaveGuest(models.Guest{FirstName: "Bob", LastName: "Brown", Email: "bob@example.com"})

	req, _ := http.NewRequest("GET", "/guests/count", nil)
	rr := httptest.NewRecorder()
	controller.CountGuests(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected %v, got %v", http.StatusOK, status)
	}

	var result map[string]int
	json.NewDecoder(rr.Body).Decode(&result)
	if result["total"] != 2 {
		t.Errorf("expected total to be 2, got %v", result["total"])
	}
}

func TestSearchGuests(t *testing.T) {
	controller := setup()

	// Simulate adding guests to search
	storage.SaveGuest(models.Guest{FirstName: "Charlie", LastName: "White", Email: "charlie@example.com"})
	storage.SaveGuest(models.Guest{FirstName: "David", LastName: "Green", Email: "david@example.com"})

	req, _ := http.NewRequest("GET", "/guests/search?name=Charlie", nil)
	rr := httptest.NewRecorder()
	controller.SearchGuests(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected %v, got %v", http.StatusOK, status)
	}

	var guests []models.Guest
	err := json.NewDecoder(rr.Body).Decode(&guests)
	if err != nil {
		t.Errorf("could not parse response: %v", err)
	}

	if len(guests) != 1 || guests[0].FirstName != "Charlie" {
		t.Error("expected to find guest named Charlie")
	}
}
