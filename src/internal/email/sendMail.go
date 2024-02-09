package email

import (
	"fmt"
	"hasura-jwt/internal/config"
	"net/smtp"
)

func SendVerifyMail(recipients string, token string) error {
	appConfig := config.LoadConfig()

	// Set up authentication information.
	username := appConfig.SMTPUser
	password := appConfig.SMTPPassword
	serverAddr := appConfig.SMTPHost + ":" + appConfig.SMTPPort

	// Create the verification link
	verificationLink := fmt.Sprintf("%s/verify/%s", appConfig.AppURL, token)
	fmt.Println("Verification link:", verificationLink)

	// Create the verification message
	message := []byte("To: " + recipients + "\r\n" +
		"Subject: Email verification\r\n" +
		"\r\n" +
		"Please click on the following link to verify your e-mail address:" + verificationLink)

	// Send the email
	auth := smtp.PlainAuth("", username, password, appConfig.SMTPHost)

	err := smtp.SendMail(serverAddr, auth, username, []string{recipients}, message)
	if err != nil {
		fmt.Println("can not use SendMail:", err)
	}

	return err
}
