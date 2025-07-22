package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
)

type User struct {
	Name string	`json:"name"`
	Age  int	`json:"age"`
}

var Users = []User{
	{Name: "alice", Age: 42},
	{Name: "mara", Age: 28},
	{Name: "bob", Age: 26},
}

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
				return
			}

			name := segments[0]
			var matches []User
			for _, p := range Users{
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
			var p User

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
			w.Write([]byte(`{"message": "Server-side User created succesfully"!}`))
			json.NewEncoder(w).Encode(Users) // or: Encode(Users) to return all users
			return
		
		case http.MethodDelete:
			var p User

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
					w.Write([]byte(`{"message": "Server-side User deleted succesfully"!}`))
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

}

func main() {

	fmt.Println("Hello Server, what port would you like to connect to?")

	for{

		//an improvement could be switching to scanner but its minimal 
		reader := bufio.NewReader(os.Stdin)
		line, err0 := reader.ReadString('\n')
		if err0 != nil {
			fmt.Print("STOPPED")
			return
		}

		line = strings.TrimSpace(line)
		line = ":" + line
		//fmt.Println(line)
		
		mux := http.NewServeMux()
		mux.HandleFunc("/healthcheck", healthcheck)
		mux.HandleFunc("/users", userManagement)
		mux.HandleFunc("/apps", appsManagement)

		err := http.ListenAndServe(line, mux)
		if err != nil {
			fmt.Println(err)
			fmt.Println("That port doesnt exist, the format is xxxx, where x is from [0,9]. Try again!")
			continue
		}
	}
}
