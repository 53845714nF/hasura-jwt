package auth

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"hasura-jwt/internal/config"
	"hasura-jwt/internal/email"
	"hasura-jwt/internal/graphql"
	"hasura-jwt/internal/model"
	"io"
	"net/http"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// read request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// parse the body as action payload
	var actionPayload model.ActionPayloadSignup
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Send the request params to the Action's generated handler function
	result, err := Signup(actionPayload.Input)

	// throw if an error happens
	if err != nil {
		errorObject := model.GraphQLError{
			Message: err.Error(),
		}
		errorBody, _ := json.Marshal(errorObject)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errorBody)
		return
	}

	// Write the response as JSON
	data, _ := json.Marshal(result)
	_, _ = w.Write(data)

}

// Signup function that takes the Action parameters and must return its response type
func Signup(args model.SignupArgs) (response model.CreateUserOutput, err error) {

	appConfig := config.LoadConfig()

	// Check if user already exists
	_, _, err = graphql.UserByEmail(appConfig.HasuraURL, appConfig.HasuraSecret, args.Email)
	if err == nil {
		return model.CreateUserOutput{}, fmt.Errorf("user already exists")
	}

	// The Password must be at least 12 characters long
	if len(args.Password) < 12 {
		return model.CreateUserOutput{}, fmt.Errorf("password must be at least 12 characters long")
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error when hashing the password: %s", err)
	}

	var status string

	if appConfig.EmailVerification {
		// Generate a token for the user
		token := GenerateToken(args.Name, args.Email, string(passwordHash))

		// Send the verification email
		err = email.SendVerifyMail(args.Email, token)
		if err != nil {
			fmt.Printf("Error when sending verification email: %s\n", err)
			return model.CreateUserOutput{}, fmt.Errorf("failed to send verification email: %s\n", err)
		}
		status = "Verification email sent"
	} else {
		// Create User in Hasura
		graphql.CreateUserMutation(appConfig.HasuraURL, appConfig.HasuraSecret, args.Name, args.Email, string(passwordHash))
		status = "User created"
	}

	// Create a response only for hasura, the user is not created yet
	response = model.CreateUserOutput{
		Status: status,
	}
	fmt.Println("Signup response:", response)
	return response, nil
}
