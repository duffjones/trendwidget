package trendapi

//package developed from google tutorial https://developers.google.com/identity/protocols/OAuth2WebServer
import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("HOLDING: Cannot find token - Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

func StartAPI(sheet string, srange string) (data *sheets.ValueRange) {
	fmt.Println("PASS: Calling API  ...")
	return runner(sheet, srange)
}

func runner(spreadsheetId string, readRange string) (data *sheets.ValueRange) {
	//spreadsheetId = link to sheet
	//readRange = range of columns to load ie: "A1:B" which loads the sheets A column starting at 1, and every corresponding column B element
	//This needs to be an enviromental variable (file) located outside of binary.
	b, err := ioutil.ReadFile("trendapi/client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved client_secret.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	//This runs getClient which gets the token from the web, or a file -- depending on if it exists
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit

	if spreadsheetId == "" {
		fmt.Println("FAIL: Spreadsheet null -- Using Default")
		spreadsheetId = "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
	}
	//This is the range of the table to import
	if readRange == "" {
		fmt.Println("FAIL Read Range Null -- Using Default")
		readRange = "A1:B"
		//to get the entire sheet, you just type the sheet name, ie 'Class Data'
		// readRange := "Class Data!A1:B2"
	}

	//this returns the actual spreadsheet data inside the range declared above.
	//this should be the return of Runner() -- which can then be accessed/used by Trend class
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	return resp
}
