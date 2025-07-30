package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"rest_api/types"
)

const path = "./db/users.json"
const path1 = "./db/apps.jsonl"

// try getting logs and see if you can measure latency and memory usage betweeen the 2 implementations

func LoadUsers() {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		// if the file doesnt exist create it cuz i dont think im doing that here
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

func LoadApps(name string) ([]types.App, error) {
	file, err := os.Open(path1)
	if err != nil {
		if os.IsNotExist(err) {
			return []types.App{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var matches []types.App
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var app types.App
		if err := json.Unmarshal(scanner.Bytes(), &app); err != nil {
			continue
		}
		if name == "" || app.Name == name {
			matches = append(matches, app)
		}
	}
	return matches, scanner.Err()
}

func AppendApp(app types.App) error {
	file, err := os.OpenFile(path1, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(app)
}

func DeleteApp(target types.App) (bool, error) {
	input, err := os.Open(path1)
	if err != nil {
		return false, err
	}
	defer input.Close()

	tempPath := path1 + ".tmp"
	output, err := os.Create(tempPath)
	if err != nil {
		return false, err
	}
	defer output.Close()

	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	found := false

	for scanner.Scan() {
		var app types.App
		if err := json.Unmarshal(scanner.Bytes(), &app); err != nil {
			continue
		}

		if app.Name == target.Name && app.Born == target.Born && app.Price == target.Price {
			found = true
			continue // skip writing this one
		}

		writer.Write(scanner.Bytes())
		writer.Write([]byte("\n"))
	}
	writer.Flush()
	input.Close()
	output.Close()

	if found {
		os.Rename(tempPath, path1)
	} else {
		os.Remove(tempPath)
	}

	return found, scanner.Err()

}
