package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppConfig struct {
	HasuraURL         string
	HasuraSecret      string
	JwtKey            string
	EmailVerification bool
	AppURL            string
	SMTPHost          string
	SMTPPort          string
	SMTPUser          string
	SMTPPassword      string
}

func LoadConfig() *AppConfig {
	HasuraURL, ok := os.LookupEnv("HASURA_URL")
	if !ok {
		fmt.Println("HASURA_URL not set. Using default value...")
		HasuraURL = "http://graphql-engine:8080/v1/graphql"
	}

	HasuraSecret, ok := os.LookupEnv("HASURA_SECRET")
	if !ok {
		fmt.Println("HASURA_SECRET not set. Exiting...")
		os.Exit(1)
	}

	JwtKey, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		fmt.Println("JWT_KEY not set. Exiting...")
		os.Exit(1)
	}

	EmailVerificationEnv := os.Getenv("EMAIL_VERIFICATION")
	if EmailVerificationEnv == "" {
		fmt.Println("EMAIL_VERIFICATION not set. Using default value (true)...")
		EmailVerificationEnv = "true"
	}

	EmailVerification, err := strconv.ParseBool(EmailVerificationEnv)
	if err != nil {
		fmt.Println("Failed to parse EMAIL_VERIFICATION. Using default value (true)...")
		EmailVerification = true
	}

	AppURL := ""
	SMTPHost := ""
	SMTPPort := ""
	SMTPUser := ""
	SMTPPassword := ""

	if EmailVerification {
		AppURL, ok = os.LookupEnv("APP_URL")
		if !ok {
			fmt.Println("APP_URL not set. Exiting...")
			os.Exit(1)
		}

		SMTPHost, ok = os.LookupEnv("SMTP_HOST")
		if !ok {
			fmt.Println("SMTP_HOST not set. Exiting...")
			os.Exit(1)
		}

		SMTPPort, ok = os.LookupEnv("SMTP_PORT")
		if !ok {
			fmt.Println("SMTP_PORT not set. Use Port 587 as default...")
			SMTPPort = "587"
		}

		SMTPUser, ok = os.LookupEnv("SMTP_USER")
		if !ok {
			fmt.Println("SMTP_USER not set. Exiting...")
			os.Exit(1)
		}

		SMTPPassword, ok = os.LookupEnv("SMTP_PASSWORD")
		if !ok {
			fmt.Println("SMTP_PASSWORD not set. Exiting...")
			os.Exit(1)
		}
	}

	return &AppConfig{
		HasuraURL:         HasuraURL,
		HasuraSecret:      HasuraSecret,
		JwtKey:            JwtKey,
		EmailVerification: EmailVerification,
		AppURL:            AppURL,
		SMTPHost:          SMTPHost,
		SMTPPort:          SMTPPort,
		SMTPUser:          SMTPUser,
		SMTPPassword:      SMTPPassword,
	}
}
