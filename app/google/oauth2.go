package google

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"log"
	appconfig "main/app/config"
	"main/app/notifications"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func LoginToGoogle(w http.ResponseWriter, r *http.Request) {
	queryToken := r.URL.Query()["token"][0]

	if appconfig.ACastHookToken != queryToken {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	config, err := getConfig(youtube.YoutubeUploadScope)
	if err != nil {
		http.Error(w, "Error getting config: "+err.Error(), http.StatusInternalServerError)
		return
	}
	authUrl := config.AuthCodeURL(generateRandomString(10), oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprintf(w,
		"<html><body><a href=\""+removeTrailingSlash(authUrl)+"\"><p>Auth link. </a></body></html>",
	)
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	rand.Read(bytes)

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func removeTrailingSlash(s string) string {
	if strings.HasSuffix(s, "/") {
		return s[:len(s)-1]
	}
	return s
}

func Oauth2Callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "Code not found in the request", http.StatusBadRequest)
		return
	}

	config, err := getConfig(youtube.YoutubeUploadScope)
	if err != nil {
		log.Printf("Error getting config: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	token, err := exchangeToken(config, code)
	if err != nil {
		log.Printf("Error exchanging token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = saveToken(token)
	if err != nil {
		log.Printf("Error saving token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	notifications.SendDiscordNotification("Google authorization successfully updated")

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body>∂<p>Authorization successful. You can now close this window.</p></body></html>")
}

func GetClient(scope string) (*http.Client, error) {
	config, err := getConfig(scope)
	if err != nil {
		log.Printf("Error getting config: %v", err)
		return nil, err
	}

	tok, err := tokenFromFile()
	if err != nil {
		log.Printf("Error getting token from file: %v", err)
		return nil, err
	}

	tokenSource := config.TokenSource(context.Background(), tok)
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Printf("Error getting new token: %v", err)
		return nil, err
	}
	if newToken.AccessToken != tok.AccessToken || newToken.RefreshToken != "" {
		if err := saveToken(newToken); err != nil {
			log.Printf("Error saving refreshed token: %v", err)
			return nil, err
		}
	}

	return config.Client(context.Background(), newToken), nil
}

func getConfig(scope string) (*oauth2.Config, error) {
	clientSecretFile, err := os.ReadFile("files/client_secret.json")
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
		return nil, err
	}

	config, err := google.ConfigFromJSON(clientSecretFile, scope)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}

	googleRedirectUrl, _ := url.JoinPath(appconfig.GoogleRedirectHost, appconfig.GoogleRedirectPath)
	config.RedirectURL = googleRedirectUrl

	return config, nil
}

func exchangeToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	tok, err := config.Exchange(context.Background(), code)

	fmt.Println(tok)

	if err != nil {
		log.Printf("Unable to retrieve token: %v", err)
		return nil, err
	}
	return tok, nil
}

func tokenCacheFile() (string, error) {
	tokenCacheDir := filepath.Join("files")
	err := os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(
		tokenCacheDir,
		url.QueryEscape("youtube-go.json"),
	), err
}

func tokenFromFile() (*oauth2.Token, error) {
	file, err := tokenCacheFile()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func saveToken(token *oauth2.Token) error {
	file, err := tokenCacheFile()
	if err != nil {
		log.Printf("Unable to get token cache file path: %v", err)
		return err
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("Unable to cache oauth token: %v", err)
		return err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Printf("Unable to encode oauth token: %v", err)
	}
	return err
}
