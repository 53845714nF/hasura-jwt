package auth

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"hasura-jwt/internal/config"
	"hasura-jwt/internal/graphql"
	"hasura-jwt/internal/model"
	"io"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// read request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// parse the body as action payload
	var actionPayload model.ActionPayloadLogin
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// Send the request params to the Action's generated handler function
	result, err := Login(actionPayload.Input)

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

// Login Auto-generated function that takes the Action parameters and must return it is response type
func Login(args model.LoginArgs) (response model.JsonWebToken, err error) {

	appConfig := config.LoadConfig()
	currentTime := time.Now().Unix()

	oldHashedPassword, id, err := graphql.UserByEmail(appConfig.HasuraURL, appConfig.HasuraSecret, args.Email)
	if err != nil {
		return model.JsonWebToken{}, fmt.Errorf("no user found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(oldHashedPassword), []byte(args.Password)); err != nil {
		return model.JsonWebToken{}, fmt.Errorf("error when comparing the password")
	}

	var (
		key         []byte
		userRole    []string
		token       *jwt.Token
		tokenString string
	)

	key = []byte(appConfig.JwtKey)

	// Check if user is admin with a database query
	userRole, err = graphql.GetUserRoles(appConfig.HasuraURL, appConfig.HasuraSecret, id)
	if err != nil {
		fmt.Println("Error when getting the user roles:", err)
		return model.JsonWebToken{}, fmt.Errorf("error when getting the user roles")
	}

	type YourClaimStruct struct {
		HasuraClaims map[string]interface{} `json:"https://hasura.io/jwt/claims"`
	}

	claims := YourClaimStruct{
		HasuraClaims: map[string]interface{}{
			"x-hasura-allowed-roles": userRole,
			"x-hasura-default-role":  "anonymous",
			"x-hasura-user-id":       id,
		},
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":                          id,
		"iat":                          currentTime,
		"https://hasura.io/jwt/claims": claims.HasuraClaims,
	})

	tokenString, _ = token.SignedString(key)

	response = model.JsonWebToken{
		Token: tokenString,
	}
	return response, nil
}
