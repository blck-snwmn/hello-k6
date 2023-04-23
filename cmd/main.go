package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type payload struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("recived request")
		if r.Method == http.MethodPost {
			fmt.Println("recived post request")
			decoder := json.NewDecoder(r.Body)
			var payload payload
			err := decoder.Decode(&payload)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "Received payload: %+v", payload)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
