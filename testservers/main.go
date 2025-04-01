package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	r1 := http.NewServeMux()

	r1.HandleFunc("/one", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("GooD pYthoN hello what bad GOpher, Bad GOpHer, Good Python") // Insert your response data here
	})

	r1.HandleFunc("/two", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("") // Insert your response data here
	})

	r1.HandleFunc("/three", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]string{"bad GOpher hello what bad GOpher", "hello good python bad GOpher good python", "good python"}) // Insert your response data here
	})

	log.Print("Test servers are up and running")
	http.ListenAndServe(":8000", r1)
}
