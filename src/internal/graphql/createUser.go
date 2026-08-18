package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreateUserMutation send an GraphQl-Mutation to create a user
func CreateUserMutation(graphqlURL string, secret string, name string, email string, password string) error {

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
		return fmt.Errorf("failed to create GraphQL mutation payload: %w", err)
	}

	req, err := http.NewRequest("POST", graphqlURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request object: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hasura-Admin-Secret", secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send GraphQL query: %w", err)
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
		return fmt.Errorf("failed to read GraphQL response: %w", err)
	}

	if errs, exists := result["errors"]; exists {
		return fmt.Errorf("GraphQL error: %v", errs)
	}

	return nil
}
