package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"rest_api/metrics"
	"rest_api/types"
)

const path = "./db/users.json"
const path1 = "./db/apps.jsonl"

func LoadUsers() {
	err := metrics.MeasureFileOperation("load", "users", func() error {
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			// if the file doesnt exist create it cuz i dont think im doing that here
			if os.IsNotExist(err) {
				log.Println("users.json not found, starting with an empty user list.")
				Users = []types.User{}
				return nil // Not an error case
			}
			return err
		}

		return json.Unmarshal(fileBytes, &Users)
	})

	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}
}

func SaveUsers() {
	log.Printf("SaveUsers() called - about to measure file operation")

	err := metrics.MeasureFileOperation("save", "users", func() error {
		log.Printf("Inside SaveUsers() file operation function")

		// Ensure directory exists
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		fileBytes, err := json.MarshalIndent(Users, "", "  ")
		if err != nil {
			return err
		}

		return os.WriteFile(path, fileBytes, 0644)
	})

	if err != nil {
		log.Printf("Failed to save users: %v", err)
	} else {
		log.Printf("SaveUsers() completed successfully")
	}
}

func LoadApps(name string) ([]types.App, error) {
	var matches []types.App

	err := metrics.MeasureFileOperation("load", "apps", func() error {
		file, err := os.Open(path1)
		if err != nil {
			if os.IsNotExist(err) {
				matches = []types.App{}
				return nil
			}
			return err
		}
		defer file.Close()

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
		return scanner.Err()
	})

	return matches, err
}

func AppendApp(app types.App) error {
	return metrics.MeasureFileOperation("append", "apps", func() error {
		file, err := os.OpenFile(path1, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		return encoder.Encode(app)
	})
}

func DeleteApp(target types.App) (bool, error) {
	var found bool

	err := metrics.MeasureFileOperation("delete", "apps", func() error {
		input, err := os.Open(path1)
		if err != nil {
			return err
		}
		defer input.Close()

		tempPath := path1 + ".tmp"
		output, err := os.Create(tempPath)
		if err != nil {
			return err
		}
		defer output.Close()

		scanner := bufio.NewScanner(input)
		writer := bufio.NewWriter(output)

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

		return nil
	})

	return found, err
}
