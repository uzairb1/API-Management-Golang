package controllers

import (
	"STACKIT/models"
	"STACKIT/storage"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// GuestControllerInterface defines the methods for guest operations
type GuestControllerInterface interface {
	InitializeStorage()
	RegisterGuest(w http.ResponseWriter, r *http.Request)
	ListGuests(w http.ResponseWriter, r *http.Request)
	CountGuests(w http.ResponseWriter, r *http.Request)
	SearchGuests(w http.ResponseWriter, r *http.Request)
}

// GuestController implements the GuestControllerInterface
type GuestController struct{}

// NewGuestController creates a new instance of GuestController
func NewGuestController() GuestControllerInterface {
	return &GuestController{}
}

// InitializeStorage loads the guest data from persistent storage
func (gc *GuestController) InitializeStorage() {
	storage.LoadGuests()
}

func (gc *GuestController) RegisterGuest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var guest models.Guest
	if err := json.NewDecoder(r.Body).Decode(&guest); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(guest.FirstName) == "" {
		http.Error(w, "First name is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(guest.LastName) == "" {
		http.Error(w, "Last name is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(guest.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Validate email format
	if !isValidEmail(guest.Email) {
		http.Error(w, "Invalid email address format", http.StatusBadRequest)
		return
	}

	// Check if the guest already exists
	if gc.guestExists(guest) {
		http.Error(w, "A guest with the same name and email already exists", http.StatusConflict)
		return
	}

	// Save the guest to storage
	storage.SaveGuest(guest)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(guest)
}

// Helper function to validate email format using regex
func isValidEmail(email string) bool {
	// Simple regex for validating email format
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// Helper function to check if a guest with the same name and email already exists
func (gc *GuestController) guestExists(newGuest models.Guest) bool {
	guests := storage.GetGuests()
	for _, guest := range guests {
		if strings.EqualFold(guest.FirstName, newGuest.FirstName) &&
			strings.EqualFold(guest.LastName, newGuest.LastName) &&
			strings.EqualFold(guest.Email, newGuest.Email) {
			return true
		}
	}
	return false
}

func (gc *GuestController) ListGuests(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	allGuests := storage.GetGuests()
	start := (page - 1) * limit
	end := start + limit

	if start > len(allGuests) {
		start = len(allGuests)
	}
	if end > len(allGuests) {
		end = len(allGuests)
	}

	guests := allGuests[start:end]

	response := map[string]interface{}{
		"page":   page,
		"limit":  limit,
		"total":  len(allGuests),
		"guests": guests,
	}
	json.NewEncoder(w).Encode(response)
}

func (gc *GuestController) CountGuests(w http.ResponseWriter, r *http.Request) {
	guests := storage.GetGuests()
	json.NewEncoder(w).Encode(map[string]int{"total": len(guests)})
}

func (gc *GuestController) SearchGuests(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("name")
	if query == "" {
		http.Error(w, "Name query parameter is required", http.StatusBadRequest)
		return
	}

	var matched []models.Guest
	for _, guest := range storage.GetGuests() {
		if strings.Contains(strings.ToLower(guest.FirstName), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(guest.LastName), strings.ToLower(query)) {
			matched = append(matched, guest)
		}
	}
	json.NewEncoder(w).Encode(matched)
}
