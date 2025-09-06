package main

import (
	"fmt"
	"net/http"

	"example.com/ccat19_go/internal"
)

func main() {
	http.HandleFunc("POST /signup", internal.SignupHandler)
	http.HandleFunc("GET /accounts/{account_id}", internal.GetAccountHandler)
	fmt.Println("Running application on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
