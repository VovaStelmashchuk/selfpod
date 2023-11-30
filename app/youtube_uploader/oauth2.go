package youtube_uploader

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
)

func CreateClientEndPoint(w http.ResponseWriter, r *http.Request) {
	config := getConfig(youtube.YoutubeUploadScope)
	authUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	io.WriteString(w, "The auth link is: "+authUrl+"\n")
}

func Oauth2Callback(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Oauth2Callback\n")
	code := r.FormValue("code")
	token, err := exchangeToken(getConfig(youtube.YoutubeUploadScope), code)

	if err != nil {
		log.Fatalf("Unable to retrieve token %v", err)
	}
	saveToken(token)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Received code: %v\r\nYou can now safely close this browser window.", code)
}

func GetClient(scope string) (*http.Client, error) {
	token, err := tokenFromFile()
	if err != nil {
		return nil, err
	}

	return getConfig(scope).Client(context.Background(), token), nil
}

func getConfig(scope string) *oauth2.Config {
	clientSecretFile, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying the scope, delete your previously saved credentials
	// at ~/.credentials/youtube-go.json
	config, err := google.ConfigFromJSON(clientSecretFile, scope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	config.RedirectURL = "http://localhost:5000/oauth2callback"

	return config
}

func exchangeToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token %v", err)
	}
	return tok, nil
}

func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(
		tokenCacheDir,
		url.QueryEscape("youtube-go.json"),
	), err
}

func tokenFromFile() (*oauth2.Token, error) {
	file, err := tokenCacheFile()
	open, err := os.Open(file)

	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(open).Decode(t)
	defer open.Close()
	return t, err
}

func saveToken(token *oauth2.Token) {
	file, err := tokenCacheFile()
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
