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

var Users = []types.User{
	{Name: "alice", Age: 42},
	{Name: "mara", Age: 28},
	{Name: "bob", Age: 26},
}

var Apps = []types.App{{Name: "Bible", Born: 1600, Price: 0.0}}

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

	switch r.Method {
		case http.MethodGet:
			if path == ""{
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
			for _, p := range Users{
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

			for _, pp := range Users{
				if pp.Name == p.Name && pp.Age== p.Age{
					http.Error(w, "Server-side User already exists", http.StatusConflict)
					return
				}
			}

			Users = append(Users, p)
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

			for i, person := range Users{
				if person.Name == p.Name && p.Age == person.Age{
					Users = slices.Delete(Users, i, i+1)
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
	fmt.Println("appsManagement called with path:", r.URL.Path)

	path := strings.TrimPrefix(r.URL.Path, "/apps")
	path = strings.Trim(path, "/")

	segments := strings.Split(path, "/")

	switch r.Method {
		case http.MethodGet:
			fmt.Println("appsManagement called with path:", r.URL.Path)
			if path == ""{
				json.NewEncoder(w).Encode(Apps)
				return
			}
			if len(segments) > 1 {
				http.Error(w, "Server-side Bad request: additional fields", http.StatusBadRequest)
				return
			}

			name := segments[0]
			var matches []types.App
			for _, app := range Apps{
				if app.Name == name {
					matches = append(matches, app)
				}
			}

			if len(matches) == 0 {
				http.Error(w, "Server-side: No such apps", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(matches)
				

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

			for _, element := range Apps{
				if element.Name == app.Name && element.Born == app.Born && element.Price == app.Price{
					http.Error(w, "Server-side: that application already exists.", http.StatusConflict)
					return
				}
			}

			Apps = append(Apps, app)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"message": "Server-side App created succesfully!"}`))
			json.NewEncoder(w).Encode(Apps) 
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

			for i, element := range Apps{
				if element.Name == app.Name && element.Born == app.Born && element.Price == app.Price{
					Apps = slices.Delete(Apps, i, i+1)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"message": "Server-side App deleted succesfully!"}`))
					json.NewEncoder(w).Encode(Apps)
					return
				}
			}

			http.Error(w, "Server-side: App doesn't exist", http.StatusNotFound)
			return

		default:
			http.Error(w, "Server-side method not allowed", http.StatusMethodNotAllowed)
	}

}

// maybe be able to parse sys args to specify a port, or check the sys for an available port
func main() {
	var port string
	flag.Parse()

	if !validators.IsFlag("port"){
		fmt.Println("Hello Server, what port would you like to connect to?")

		for{
			//an improvement could be switching to scanner but its minimal 
			reader := bufio.NewReader(os.Stdin)
			port, err0 := reader.ReadString('\n')
			if err0 != nil {
				fmt.Print("STOPPED.")
				return
			}

			port = strings.TrimSpace(port)
			if port == "" { 
				fmt.Println("Port cannot be empty, please try again.")
				continue
			}

			if !validators.IsValid(port) {
				fmt.Println("Invalid format. Pick a port from 1-65336.")
				continue
			}

			if !validators.IsAvailable(port){
				fmt.Println("This port is busy. Try something else.")
				continue
			}
			break
		}
	}
		port = *Portflag
		port = ":" + port
		mux := SetupRoutes()

		fmt.Printf("Starting server on port: %s ...\n", port)
		err := http.ListenAndServe(port, mux)
		if err != nil {
			fmt.Print("Somethig unexpected happened: ")
			fmt.Println(err)
			fmt.Println("Maybe try again.")
			
		}
	
}
