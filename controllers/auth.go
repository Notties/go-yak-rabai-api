// controllers/auth.go
package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"yak.rabai/config"
	"yak.rabai/models"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     "YOUR_GOOGLE_CLIENT_ID",     // Replace with your Client ID
		ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET", // Replace with your Client Secret
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
)

// GoogleLogin initiates the OAuth login
func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the OAuth callback and user creation
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")

	// Exchange the authorization code for an access token
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	var userInfo map[string]interface{}
	json.Unmarshal(data, &userInfo)

	// Save user info to database
	googleID := userInfo["id"].(string)
	var user models.User
	config.DB.Where("google_id = ?", googleID).FirstOrCreate(&user, models.User{
		GoogleID: googleID,
		Name:     userInfo["name"].(string),
		Email:    userInfo["email"].(string),
	})

	// Store session or token as needed
	c.JSON(http.StatusOK, gin.H{"message": "User authenticated", "user": user})
}
