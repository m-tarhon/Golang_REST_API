package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"rest_api/types"
	"rest_api/validators"
	"strings"
)
var ( 
	Portflag = flag.String("port", "8080", "specifies which port to use to connect the server to")
)

func userSide(port, username, password string) {
	base_url := "http://localhost:" + port + "/users"
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
		case "switch":
			fmt.Println("you're going to the app endpoint bye")
			appSide(port)

		case "help":
			fmt.Println("Commands: ")
			fmt.Println(" 1. add <a user> <their age> --> adds an user")
			fmt.Println(" 2. get <a user> --> gets information on a specific user")
			fmt.Println(" 3. list --> which returns all users")
			fmt.Println(" 4. delete <name> <age> --> delets an user's info")
			fmt.Println(" 5. exit")
			fmt.Println(" 6. switch --> takes you to the app endpoint")

		case "add":
			if len(args) != 3 {
				fmt.Println("You need only name and age!")
				continue
			}
			name := args[1]
			var age int
			fmt.Sscanf(args[2], "%d", &age)

			u := types.User{Name: name, Age: age} // check for type
			body, err := json.Marshal(u)
			if err != nil {
				fmt.Println("Failed to format the user into JSON: ", err)
				continue
			}

			//resp, err := http.Post(base_url, "application/json", bytes.NewReader(body))
			resp, err := Helper4Auth(http.MethodPost, base_url, username, password, bytes.NewReader(body))
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
			if len(args) != 2 {
				fmt.Println("You only need the name")
				continue
			}

			name := args[1]
			resp, err := Helper4Auth(http.MethodGet, base_url+"/"+name, username, password, nil)
			if err != nil {
				fmt.Println("request failed: ", err)
				continue
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK:
				var users []types.User
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
			//resp, err := http.Get(base_url)
			resp, err := Helper4Auth(http.MethodGet, base_url, username, password, nil)
			if err != nil {
				fmt.Println("request failed: ", err)
				continue
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK:
				var us []types.User
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

			u := types.User{Name: name, Age: age}

			body, err := json.Marshal(u)
			if err != nil {
				fmt.Println("Failed to format the user into JSON: ", err)
				continue
			}
			// req, err := http.NewRequest(http.MethodDelete, base_url, bytes.NewReader(body))
			resp, err := Helper4Auth(http.MethodDelete, base_url, username, password, bytes.NewReader(body))
			if err != nil {
				fmt.Println("HTTP request failed:", err)
				continue
			}

			// req.Header.Set("Content-Type", "application/json")

			// client := &http.Client{}
			// resp, err := client.Do(req)
			// if err != nil {
			// 	fmt.Println("HTTP request failed:", err)
			// 	continue
			// }
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

func appSide(port string) {
	base_url := "http://localhost:" + port + "/apps"

	fmt.Println("App Client, type help for options or exit to quit")
	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("§§§ ")
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
			fmt.Println(" 1. add <name> <year> <price> --> adds an app")
			fmt.Println(" 2. get <an app> --> gets information on a specific app")
			fmt.Println(" 3. list --> which returns all apps")
			fmt.Println(" 4. delete <name> <DoB> <price> --> delets an app's info")
			fmt.Println(" 5. exit")

		case "add":
			if len(args) != 4 {
				fmt.Println("You need the name, year of creation, and price of the app!")
				continue
			}
			name := args[1]
			var dob int
			var price float32
			fmt.Sscanf(args[2], "%d", &dob)
			fmt.Sscanf(args[3], "%f", &price)

			a := types.App{Name: name, Born: dob, Price: price} // check for type
			body, err := json.Marshal(a)
			if err != nil {
				fmt.Println("Failed to format the app's info into JSON: ", err)
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
				fmt.Println("App added succesfully!!")
				continue
			case http.StatusConflict:
				fmt.Println("Error: App already exists")
				continue
			case http.StatusBadRequest:
				fmt.Println("Error: invalid app data")
				continue
			default:
				respBody, _ := io.ReadAll(resp.Body)
				fmt.Printf("Server error: %s\n", string(respBody))
				//fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
				continue
			}

		case "get":
			if len(args) != 2 {
				fmt.Println("You only need the app's name")
				continue
			}

			name := args[1]
			fmt.Println("Client requesting:", base_url+"/"+name)
			resp, err := http.Get(base_url + "/" + name)
			if err != nil {
				fmt.Println("request failed: ", err)
				continue
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK:
				var apps []types.App
				if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
					fmt.Println("Failed to parse response:", err)
					continue
				}

				for _, element := range apps {
					fmt.Printf("%s was created in %d and its price is %f\n", element.Name, element.Born, element.Price)
				}
				continue

			case http.StatusNotFound:
				fmt.Println("Error: app not found.")
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
				var apps []types.App
				if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
					fmt.Println("Failed to parse response:", err)
					continue
				} else {
					for _, element := range apps {
						fmt.Printf("%s was created in %d and its price is %f\n", element.Name, element.Born, element.Price)
					}
					continue
				}

			default:
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Server error: %s\n", string(body))
				continue
			}

		case "delete":
			if len(args) != 4 {
				fmt.Println("You need the name, year of creation, and price of the app!")
				continue
			}
			name := args[1]
			var dob int
			var price float32
			fmt.Sscanf(args[2], "%d", &dob)
			fmt.Sscanf(args[3], "%f", &price)

			a := types.App{Name: name, Born: dob, Price: price} // check for type

			body, err := json.Marshal(a)
			if err != nil {
				fmt.Println("Failed to format the user into JSON: ", err)
				continue
			}
			req, err := http.NewRequest(http.MethodDelete, base_url, bytes.NewReader(body))
			if err != nil {
				fmt.Println("HTTP request failed:", err)
				continue
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("HTTP request failed:", err)
				continue
			}
			defer resp.Body.Close()

			switch resp.StatusCode {
			case http.StatusOK:
				fmt.Println("App deleted successfully!")
				continue
			case http.StatusNotFound:
				fmt.Println("Error: app not found.")
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

// focus on making this less monolithic, bcs its horrible
// add database man
// try adding besides http a way to get into https

func main() {
	var port, endpoint string
	flag.Parse()
	if !validators.IsFlag("port"){
			
		fmt.Println("Hello Client, which port do you want?")
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if !scanner.Scan() {
				fmt.Print("STOPPED")
				return
			}
			port = scanner.Text()
			var port_int int
			fmt.Sscanf(port, "%d", &port_int)
			if !validators.IsValid(port){
				fmt.Println("Nonexistent port, format is xxxx. Try again.")
				continue
			}
			break
		}
	} else{
		port = *Portflag
	}

	fmt.Println("Which endpoint do you want to vist: type 'users' or 'apps'")
	for {
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			fmt.Print("STOPPED")
			return
		}
		endpoint = scanner.Text()

		if endpoint == "users" || endpoint == "apps" {
			break
		}
		fmt.Println("Nonexistent endpoint, try again.")
		continue
	}

	//if you dont have the boss:man credentials wont be allowed to get into users
	if endpoint == "users" {
		var username, password string
		for {
			fmt.Println("This endpoint requires authentication")
			fmt.Println("please give your username: ")
			fmt.Scanln(&username)

			fmt.Println("please give your password: ")
			fmt.Scanln(&password)

			err := AuthReq(username, password, port)
			if err != nil {
				fmt.Printf("Auth error: %v\n", err)
				fmt.Println("try again")
				continue
			}
			break
		}
		userSide(port, username, password)
	} else {
		appSide(port)
	}
}
