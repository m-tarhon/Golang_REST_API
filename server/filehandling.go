package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"rest_api/types"
)

const path = "./db/users.json"
const path1 = "./db/apps.json"

func LoadUsers() {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("users.json not found, starting with an empty user list.")
			Users = []types.User{}
			return
		}
		log.Fatalf("Failed to read users.json: %v", err)
	}

	err = json.Unmarshal(fileBytes, &Users)
	if err != nil {
		log.Fatalf("Failed to parse users.json: %v", err)
	}
}

func SaveUsers() {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("Failed to create db directory: %v", err)
		return
	}

	fileBytes, err := json.MarshalIndent(Users, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal users: %v", err)
		return
	}

	err = os.WriteFile(path, fileBytes, 0644)
	if err != nil {
		log.Printf("Failed to write users.json: %v", err)
	}
}

func LoadApps() {
	fileBytes,err := os.ReadFile(path1)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("apps.json not found. starting with am empty list")
			Apps = []types.App{}
			return 
		}
		log.Fatalf("failed to read apps.json: %v", err)
	}
	err = json.Unmarshal(fileBytes, &Apps)
	if err != nil {
		log.Fatalf("failed to parse apps.json: %v", err)
	}
}

func SaveApps() { 
	dir := filepath.Dir(path1)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("Failed to create db directory: %v", err)
		return 
	}

	fileBytes, err := json.MarshalIndent(Apps, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal apps: %v", err)
		return
	}

	err = os.WriteFile(path1, fileBytes, 0644)
	if err != nil {
		log.Printf("Failed to write apps.json: %v", err)
	}
}
