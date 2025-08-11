package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"  // Add this import
	"net/http"
	"rest_api/types"
	"rest_api/validators"
	"rest_api/metrics"
	"slices"
	"strings"
)

var Users = []types.User{}
var Apps = []types.App{}
var (
	HTTPflag = flag.String("http", "8080", "specifies which port to use to connect the server to")
	HTTPSflag = flag.String("https", "8081", "specifies which port to use to connect the server to")
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
			if person.Name == p.Name && person.Age == p.Age {
				Users = slices.Delete(Users, i, i+1)
				
				// Add debug logging
				log.Printf("About to call SaveUsers() for deletion of %s", p.Name)
				SaveUsers()
				log.Printf("SaveUsers() completed for deletion")
				
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
	metrics.Init()

	var port, ports string
	flag.Parse()

	if !validators.IsFlag("http") {
		validators.SMTH("http", &port)
	} else {
		port = *HTTPflag
	}
	if !validators.IsFlag("https") {
		validators.SMTH("https", &ports)
	} else {
		ports = *HTTPSflag
	}

	port = ":" + port
	ports = ":" + ports

	mux := SetupRoutes()

	fmt.Printf("Starting HTTP server on port %s ...\n", port)
	go func(){
		err := http.ListenAndServe(port, mux)
		if err != nil {
			fmt.Print("Something unexpected happened: ")
			fmt.Println(err)
			fmt.Println("Maybe try again.")
		}
	}()

	fmt.Printf("Starting HTTPS server on port %s ...\n", ports)
	errs := http.ListenAndServeTLS(ports, "cert.pem", "key.pem", mux)
	if errs != nil{
		fmt.Print("Somethig unexpected happened: ")
		fmt.Println(errs)
		fmt.Println("Maybe try again.")
	}
}
