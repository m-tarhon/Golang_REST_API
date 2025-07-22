package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// change port to make it dynamic 
func main() {
	base_url := "http://localhost:8080/users"
	fmt.Println("User Client, type help for options or exit to quit")
	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if line == "exit" || line == "quit" {
			fmt.Print("Goodbye!")
			return
		}

		args := strings.Fields(line)
		if len(args) < 1 {
			continue
		}

		cmd := args[0]
		switch cmd {
		case "help":
			fmt.Println("Commands: ")
			fmt.Println(" 1. add <a user> <their age> --> adds an user")
			fmt.Println(" 2. get <a user> --> gets information on a specific user")
			fmt.Println(" 3. list --> which returns all users")
			fmt.Println(" 4. delete <name> <age> --> delets an user's info")
			fmt.Println(" 5. exit")

		case "add":
			if len(args) != 3 {
				fmt.Println("You need only name and age!")
				continue
			}
			name := args[1]
			var age int
			fmt.Sscanf(args[2], "%d", &age)

			u := User{Name: name, Age: age} // check for type 
			body, err := json.Marshal(u)
			if err != nil {
				fmt.Println("Failed to format the user into JSON: ", err)
				continue
			}

			resp, err := http.Post(base_url, "application/json", bytes.NewReader(body))
			if err != nil {
				fmt.Println("HTTP request failed:", err)
				continue
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusCreated:
				fmt.Println("User added succesfully!!")
				continue
			case http.StatusConflict:
				fmt.Println("Error: user already exists")
				continue
			case http.StatusBadRequest:
				fmt.Println("Error: invalid user data")
				continue
			default:
				respBody, _ := io.ReadAll(resp.Body)
				fmt.Printf("Server error: %s\n", string(respBody))
				//fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
				continue
			}

		case "get":
			if len(args) > 2 {
				fmt.Println("You only need the name")
				continue
			}

			name := args[1]
			resp, err := http.Get(base_url + "/" + name)
			if err != nil {
				fmt.Println("request failed: ", err)
				continue
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK:
				var users []User
				if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
					fmt.Println("Failed to parse response:", err)
					continue
				}
				
				for _, u := range users {
					fmt.Printf("%s's age is %d\n", u.Name, u.Age)
				}
				continue

			case http.StatusNotFound:
				fmt.Println("Error: user not found.")
				continue

			case http.StatusBadRequest:
				fmt.Println("Error: bad request.")
				continue

			default:
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Server error: %s\n", string(body))
				continue
			}

		case "list":
			resp, err := http.Get(base_url)
			if err != nil {
				fmt.Println("request failed: ", err)
				continue
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK:
				var us []User
				if err := json.NewDecoder(resp.Body).Decode(&us); err != nil {
					fmt.Println("Failed to parse response:", err)
					continue
				} else {
					for _, user := range us {
						fmt.Printf("%s's age is %d\n", user.Name, user.Age)
					}
					continue
				}

			default:
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Server error: %s\n", string(body))
				continue
			}

		case "delete":
			if len(args) != 3 {
				fmt.Println("You need only name and age!")
				continue
			}
			name := args[1]
			var age int
			fmt.Sscanf(args[2], "%d", &age)

			u := User{Name: name, Age: age}

			body, err := json.Marshal(u)
			if err != nil {
				fmt.Println("Failed to format the user into JSON: ", err)
				continue
			}
			req, err := http.NewRequest(http.MethodDelete, base_url, bytes.NewReader(body))
			if err != nil {
				fmt.Println("HTTP request failed:", err)
				continue
			}

			req.Header.Set("Content-Type", "application/json" )

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("HTTP request failed:", err)
				continue	
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
				case http.StatusOK:
					fmt.Println("User deleted successfully!")
					continue
				case http.StatusNotFound:
					fmt.Println("Error: user not found.")
					continue
				case http.StatusBadRequest:
					fmt.Println("Error: bad request.")
					continue
				default:
					respBody, _ := io.ReadAll(resp.Body)
					fmt.Printf("Server error: %s\n", string(respBody))
					continue
			}

		default:
			fmt.Println("unknown command:", cmd)
			fmt.Println("Type 'help' for a list of available commands. ")
			continue
		}
	}
}
