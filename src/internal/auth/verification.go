package auth

import (
	"crypto/rand"
	"encoding/base64"
	"hasura-jwt/internal/config"
	"hasura-jwt/internal/graphql"
	"hasura-jwt/internal/model"
	"html/template"
	"net/http"
	"sync"
	"time"
)

var UserTokens = make(map[string]model.VerificationData)

// Mutex for the concurrent map access
var mu sync.Mutex

func GenerateToken(name string, email string, passwordHash string) string {
	length := 32
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	token := randomString[:length]
	mu.Lock()
	defer mu.Unlock()
	UserTokens[token] = model.VerificationData{
		Name:           name,
		Email:          email,
		PasswordHash:   passwordHash,
		ExpirationTime: time.Now().Add(10 * time.Minute),
	}

	return token
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	appConfig := config.LoadConfig()
	token := r.PathValue("token")

	mu.Lock()
	defer mu.Unlock()

	// Check if the token is valid
	verificationData, ok := UserTokens[token]
	if !ok || time.Now().After(verificationData.ExpirationTime) {
		http.NotFound(w, r)
		return
	}

	// Create the user in the database
	graphql.CreateUserMutation(appConfig.HasuraURL, appConfig.HasuraSecret, verificationData.Name, verificationData.Email, verificationData.PasswordHash)

	delete(UserTokens, token)

	tmpl, err := template.New("verification").Parse(`
	<html>
		<body>
			<h1>Successful verification</h1>
			<p>Thank you for verifying your e-mail address.</p>
		</body>
	</html>
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
