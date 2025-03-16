package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test API endpoints
func TestAPIs(t *testing.T) {
	// 1️.Test Add Distributor
	t.Run("Add Distributor", func(t *testing.T) {
		reqBody := []byte(`{"name": "DISTRIBUTOR1"}`)
		req, err := http.NewRequest("POST", "/add-distributor", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(addDistributor)
		handler.ServeHTTP(rr, req)

		// if rr.Code != http.StatusOK {
		// 	t.Errorf("Expected status OK, got %d", rr.Code)
		// }
		if rr.Code != http.StatusCreated { // 201 instead of 200
			t.Errorf("Expected status Created, got %d", rr.Code)
		}

	})

	// 2️. Test Set Permission
	t.Run("Set Permission", func(t *testing.T) {
		reqBody := []byte(`{
			"name": "DISTRIBUTOR1",
			"includes": ["INDIA", "UNITEDSTATES"],
			"excludes": ["KARNATAKA-INDIA", "CHENNAI-TAMILNADU-INDIA"]
		}`)
		req, err := http.NewRequest("POST", "/set-permission", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(setPermissions)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rr.Code)
		}
	})

	// 3️. Test Check Permission - Should return "YES"
	t.Run("Check Permission - YES", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/check-permission?name=DISTRIBUTOR1&region=CHICAGO-ILLINOIS-UNITEDSTATES", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(checkPermission)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rr.Code)
		}

		var resp map[string]string
		json.Unmarshal(rr.Body.Bytes(), &resp)

		if resp["permission"] != "YES" {
			t.Errorf("Expected permission YES, got %s", resp["permission"])
		}
	})

	// 4️. Test Check Permission - Should return "NO"
	t.Run("Check Permission - NO", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/check-permission?name=DISTRIBUTOR1&region=CHENNAI-TAMILNADU-INDIA", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(checkPermission)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rr.Code)
		}

		var resp map[string]string
		json.Unmarshal(rr.Body.Bytes(), &resp)

		if resp["permission"] != "NO" {
			t.Errorf("Expected permission NO, got %s", resp["permission"])
		}
	})
}
