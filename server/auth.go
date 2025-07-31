package main

import (
    "encoding/base64"
    "net/http"
    "strings"
)

func basicAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if auth == "" || !strings.HasPrefix(auth, "Basic ") {
            w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Decode credentials
        payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
        if err != nil {
            http.Error(w, "Invalid auth header", http.StatusBadRequest)
            return
        }

        pair := strings.SplitN(string(payload), ":", 2)
        if len(pair) != 2 || pair[0] != "boss" || pair[1] != "man" {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    }
}
