package main

import (
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello World! </hr"))
}
