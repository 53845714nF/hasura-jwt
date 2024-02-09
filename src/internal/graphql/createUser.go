package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreateUserMutation send an GraphQl-Mutation to create a user
func CreateUserMutation(graphqlURL string, secret string, name string, email string, password string) {

	query := `
		mutation CreateUser($name: String!, $email: String!, $password: String!) {
			insert_user_one(object: {name: $name, email: $email, password: $password }) {
				name
				email
				password
			}
		}
	`

	variables := map[string]interface{}{
		"name":     name,
		"email":    email,
		"password": password,
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		fmt.Println(" failed to create GraphQL-Mutation:", err)
		return
	}

	req, err := http.NewRequest("POST", graphqlURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("failed to create to HTTP-Request object:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hasura-Admin-Secret", secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed to send GraphQl-Query", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(resp.Body)

	// read the response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("failed to read GraphQl-Response", err)
		return
	}

	// Show the result
	fmt.Printf("Created User: %s\n", result)
}
