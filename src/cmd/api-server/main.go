package main

import (
	"fmt"
	"hasura-jwt/internal/auth"
	"hasura-jwt/internal/config"
	"log"
	"net/http"
)

func main() {
	config.LoadConfig()
	fmt.Println("Starting Server ...")
	mux := http.NewServeMux()
	mux.HandleFunc("/signup", auth.SignupHandler)
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/verify/{token}", auth.VerifyHandler)

	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
