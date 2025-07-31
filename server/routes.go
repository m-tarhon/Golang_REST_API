package main

import "net/http"

func SetupRoutes() *http.ServeMux {
		mux := http.NewServeMux()
		
		mux.HandleFunc("/healthcheck", healthcheck)
		mux.HandleFunc("/users", basicAuth(userManagement))
		mux.HandleFunc("/users/", basicAuth(userManagement))
		mux.HandleFunc("/apps", appsManagement)
		mux.HandleFunc("/apps/", appsManagement)

		return mux
}
