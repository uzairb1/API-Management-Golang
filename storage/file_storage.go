package storage

import (
	"STACKIT/models"
	"encoding/json"
	"os"
	"sync"
)

var (
	guests   []models.Guest
	dataFile = "guests.json"
	mu       sync.Mutex
)

func LoadGuests() {
	file, err := os.ReadFile(dataFile)
	if err == nil {
		json.Unmarshal(file, &guests)
	}
}

func SaveGuest(guest models.Guest) {
	mu.Lock()
	defer mu.Unlock()
	guests = append(guests, guest)
	data, _ := json.Marshal(guests)
	os.WriteFile(dataFile, data, 0644)
}

func GetGuests() []models.Guest {
	return guests
}
