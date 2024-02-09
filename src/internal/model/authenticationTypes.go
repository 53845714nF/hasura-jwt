package model

import "time"

type ActionPayloadSignup struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            SignupArgs             `json:"input"`
}

type ActionPayloadLogin struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            LoginArgs              `json:"input"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

type CreateUserOutput struct {
	Status string
}

type JsonWebToken struct {
	Token string
}

type Mutation struct {
	Login  *JsonWebToken
	Signup *CreateUserOutput
}

type LoginArgs struct {
	Email    string
	Password string
}

type SignupArgs struct {
	Name     string
	Email    string
	Password string
}

type VerificationData struct {
	Name           string
	Email          string
	PasswordHash   string
	ExpirationTime time.Time
}
