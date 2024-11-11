// services/auth.go
package services

import (
	"context"
	"encoding/json"

	"yak.rabai/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Set up Google OAuth configuration
var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     "YOUR_GOOGLE_CLIENT_ID",     // Replace with your Client ID
		ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET", // Replace with your Client Secret
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
)

// GetGoogleOAuthURL returns the Google OAuth URL for login
func GetGoogleOAuthURL() string {
	return googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

// GetGoogleUser retrieves Google user information and parses it into a User model
func GetGoogleUser(code string) (*models.User, error) {
	// Exchange the code for an access token
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)

	// Parse the information into a User model
	user := &models.User{
		GoogleID: userInfo["id"].(string),
		Name:     userInfo["name"].(string),
		Email:    userInfo["email"].(string),
	}

	return user, nil
}
