package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Distributor struct {
	Name     string   `json:"name"`
	Includes []string `json:"includes"`
	Excludes []string `json:"excludes"`
}

var (
	distributors = make(map[string]*Distributor)
	dataLock     sync.Mutex
	cityMapping  = make(map[string]string) // Mapping city codes to full names
)

// Load city mapping from CSV
func loadCityMapping() error {
	file, err := os.Open("regions.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip header
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		cityCode := record[0]
		fullName := strings.ToUpper(fmt.Sprintf("%s-%s-%s", record[3], record[4], record[5]))
		cityMapping[cityCode] = fullName
	}
	return nil
}

// Add a new distributor
func addDistributor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var d Distributor
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dataLock.Lock()
	distributors[d.Name] = &d
	dataLock.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Distributor added"})
}

// Set permissions for an existing distributor
func setPermissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var d Distributor
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dataLock.Lock()
	defer dataLock.Unlock()

	if existing, found := distributors[d.Name]; found {
		existing.Includes = d.Includes
		existing.Excludes = d.Excludes
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Permissions updated"})
	} else {
		http.Error(w, "Distributor not found", http.StatusNotFound)
	}
}

func checkPermission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	region := strings.ToUpper(r.URL.Query().Get("region"))

	dataLock.Lock()
	defer dataLock.Unlock()

	d, exists := distributors[name]
	if !exists {
		http.Error(w, "Distributor not found", http.StatusNotFound)
		return
	}

	//  Debugging: Print permissions to verify correctness
	// fmt.Println("Checking permissions for:", name)
	// fmt.Println("Includes:", d.Includes)
	// fmt.Println("Excludes:", d.Excludes)
	// fmt.Println("Requested Region:", region)

	for _, excl := range d.Excludes {
		if region == excl || strings.HasSuffix(region, excl) {
			json.NewEncoder(w).Encode(map[string]string{"permission": "NO"})
			return
		}
	}

	for _, incl := range d.Includes {
		if region == incl || strings.HasSuffix(region, incl) {
			json.NewEncoder(w).Encode(map[string]string{"permission": "YES"})
			return
		}
	}

	// Default case: NO
	json.NewEncoder(w).Encode(map[string]string{"permission": "NO"})
}

func main() {
	err := loadCityMapping()
	if err != nil {
		fmt.Println("Failed to load city mapping:", err)
		return
	}

	http.HandleFunc("/add-distributor", addDistributor)
	http.HandleFunc("/set-permission", setPermissions)
	http.HandleFunc("/check-permission", checkPermission)

	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
