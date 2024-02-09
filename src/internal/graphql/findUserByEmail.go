package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// UserByEmail send an GraphQl-Query to find a user by email
func UserByEmail(graphqlURL string, secret string, email string) (string, string, error) {

	query := `
		 query UserByEmail($email: String!) {
            user(where: {email: {_eq: $email}}, limit: 1) {
                password,
				id
            }
        }
	`

	variables := map[string]interface{}{
		"email": email,
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to creating GraphQl query")
	}

	// GraphQL-Query create
	req, err := http.NewRequest("POST", graphqlURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", "", fmt.Errorf("failed to create to HTTP-Request object")
	}

	// Header added
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hasura-Admin-Secret", secret)

	// GraphQL-Query send
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed while sending GraphQL-Query")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(resp.Body)

	// read response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", "", fmt.Errorf("failed to read GraphQL-Query")
	}

	data, dataExists := result["data"].(map[string]interface{})
	if dataExists {
		users, usersExist := data["user"].([]interface{})
		if usersExist && len(users) > 0 {
			user := users[0].(map[string]interface{})
			password, passwordExists := user["password"].(string)
			id := user["id"].(string)
			if passwordExists {
				return password, id, nil
			} else {
				return "", "", fmt.Errorf("failed no password found")
			}
		} else {
			return "", "", fmt.Errorf("failed no user found")
		}
	} else {
		return "", "", fmt.Errorf("failed no data found")
	}
}
