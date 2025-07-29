package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"rest_api/types"
	"rest_api/validators"
	"slices"
	"strings"
)

var Users = []types.User{}
var Apps = []types.App{}
var (
	Portflag = flag.String("port", "8080", "specifies which port to use to connect the server to")
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "status is: available")
}

func userManagement(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users")
	path = strings.Trim(path, "/")

	segments := strings.Split(path, "/")
	LoadUsers()

	switch r.Method {
	case http.MethodGet:
		if path == "" {
			json.NewEncoder(w).Encode(Users)
			return
		}
		if len(segments) > 1 {
			http.Error(w, "Server-side Bad request: additional fields", http.StatusBadRequest)
			//log.Println("Server-side Bad request: additional fields")
			return
		}

		name := segments[0]
		var matches []types.User
		//fmt.Printf("Searching for user: '%s'\n", name)
		for _, p := range Users {
			//fmt.Printf("Comparing with: '%s'\n", p.Name)
			if p.Name == name {
				matches = append(matches, p)
			}
		}

		if len(matches) == 0 {
			http.Error(w, "Server-side: User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(matches)

	case http.MethodPost:
		var p types.User

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if p.Name == "" || p.Age == 0 {
			http.Error(w, "Server-side invalid payload: missing required fields", http.StatusBadRequest)
			return
		}

		for _, pp := range Users {
			if pp.Name == p.Name && pp.Age == p.Age {
				http.Error(w, "Server-side User already exists", http.StatusConflict)
				return
			}
		}

		Users = append(Users, p)
		SaveUsers()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "Server-side User created succesfully!"}`))
		json.NewEncoder(w).Encode(Users) // or: Encode(Users) to return all users
		return

	case http.MethodDelete:
		var p types.User

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for i, person := range Users {
			if person.Name == p.Name && p.Age == person.Age {
				Users = slices.Delete(Users, i, i+1)
				SaveUsers()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "Server-side User deleted succesfully!"}`))
				json.NewEncoder(w).Encode(Users)
				return
			}
		}

		http.Error(w, "Server-side: User doesn't exist", http.StatusNotFound)
		return

	default:
		http.Error(w, "Server-side method not allowed", http.StatusMethodNotAllowed)
	}
}

func appsManagement(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("appsManagement called with path:", r.URL.Path)
	path := strings.TrimPrefix(r.URL.Path, "/apps")
	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")

	switch r.Method {
	case http.MethodGet:
		if len(segments) > 1 {
			http.Error(w, "Server-side Bad request: additional fields", http.StatusBadRequest)
			return
		}

		name := ""
		if len(segments)==1{
			name = segments[0]
		}

		apps, err := LoadApps(name)
		if err != nil {
			http.Error(w, "Failed to load apps", http.StatusInternalServerError)
			return
		}

		if len(apps) == 0 && name != "" {
			http.Error(w, "No such app found", http.StatusNotFound)
			return
		}		

		w.Header().Set("Content-Type", "application/x-ndjson")
		json.NewEncoder(w).Encode(apps)

	case http.MethodPost:
		var app types.App
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&app)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if app.Name == "" || app.Born == 0 || app.Price == 0 {
			http.Error(w, "Server-side invalid payload: missing required fields", http.StatusBadRequest)
			return
		}

		Apps, err := LoadApps(app.Name)
		if err != nil {
			http.Error(w, "Error when checking for duplicates", http.StatusInternalServerError)
			return
		}

		for _, element := range Apps {
			if element.Name == app.Name && element.Born == app.Born && element.Price == app.Price {
				http.Error(w, "Server-side: that application already exists.", http.StatusConflict)
				return
			}
		}

		if err := AppendApp(app); err != nil {
			http.Error(w, "Failed to write app", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "Server-side App created succesfully!"}`))
		return

	case http.MethodDelete:
		var app types.App
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&app)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		found, err := DeleteApp(app)
		if err != nil {
			http.Error(w, "Failed to delete", http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, "App not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "App deleted successfully!"}`))

	default:
		http.Error(w, "Server-side method not allowed", http.StatusMethodNotAllowed)
	}

}

func main() {
	var port string
	flag.Parse()

	if !validators.IsFlag("port") {
		fmt.Println("Hello Server, what port would you like to connect to?")

		for {
			//an improvement could be switching to scanner but its minimal
			reader := bufio.NewReader(os.Stdin)
			port1, err0 := reader.ReadString('\n')
			if err0 != nil {
				fmt.Print("STOPPED.")
				return
			}

			port1 = strings.TrimSpace(port1)
			if port1 == "" {
				fmt.Println("Port cannot be empty, please try again.")
				continue
			}

			if !validators.IsValid(port1) {
				fmt.Println("Invalid format. Pick a port from 1-65336.")
				continue
			}

			if !validators.IsAvailable(port1) {
				fmt.Println("This port is busy. Try something else.")
				continue
			}
			port = port1
			break
		}
	} else {
		port = *Portflag
	}
	port = ":" + port
	
	mux := SetupRoutes()

	fmt.Printf("Starting server on port %s ...\n", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Print("Somethig unexpected happened: ")
		fmt.Println(err)
		fmt.Println("Maybe try again.")

	}

}
