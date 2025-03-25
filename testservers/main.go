package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	r1 := http.NewServeMux()

	r1.HandleFunc("/one", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("hello what bad GOpher")
	})

	r1.HandleFunc("/two", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("hello python is good bad GOpher")
	})

	r1.HandleFunc("/three", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]string{"hello what bad GOpher", "hello good python bad GOpher"})
	})

	http.ListenAndServe(":8000", r1)
}
