package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetUserRoles send an GraphQl-Query to get the roles of a user
func GetUserRoles(graphqlURL string, secret string, id string) ([]string, error) {

	query := `
		query GetUserRoles($id: uuid!) {
			assigned_user_roles(where: {user_id: {_eq: $id}}) {
				user_role_name
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to creating GraphQl query")
	}

	// GraphQL-Query create
	req, err := http.NewRequest("POST", graphqlURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create to HTTP-Request object")
	}

	// Header added
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Hasura-Admin-Secret", secret)

	// GraphQL-query send
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed while sending GraphQL-Query")
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
		return nil, fmt.Errorf("failed to read GraphQL-Query")
	}

	data, dataExists := result["data"].(map[string]interface{})
	if dataExists {
		userRoles, userRolesExist := data["assigned_user_roles"].([]interface{})

		if userRolesExist && len(userRoles) > 0 {
			roles := make([]string, len(userRoles))
			for i, role := range userRoles {
				roles[i] = role.(map[string]interface{})["user_role_name"].(string)
			}

			return roles, nil
		} else {
			return nil, fmt.Errorf("failed no user roles found")
		}
	} else {
		return nil, fmt.Errorf("failed no data found")
	}
}
