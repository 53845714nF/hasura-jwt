package main

import (
	"fmt"
	"hasura-jwt/internal/auth"
	"hasura-jwt/internal/config"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		res, err := http.Get("http://127.0.0.1:3000/health")
		if err != nil || res.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	}

	config.LoadConfig()
	fmt.Println("Starting Server ...")
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/signup", auth.SignupHandler)
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/verify/{token}", auth.VerifyHandler)

	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
