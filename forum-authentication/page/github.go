package page

import (
	"bytes"
	"encoding/json"
	"fmt"
	"forum/db"
	"forum/dop"
	"forum/structs"
	"io"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type GithubUserInfo struct {
	Email  string `json:"email"`
	Name   string `json:"login"`
	Avatar string `json:"avarat_url"`
}

type GetEmailInfo struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Veryfied   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

var oauthConf = &oauth2.Config{
	ClientID:     "d4b5c606c18cf127a0ec",
	ClientSecret: "e9c4c54470d389cecd58343384f7ffbe183303be",
	RedirectURL:  "http://localhost:8083/login/github/callback",
	// Scopes:       []string{"user:email"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
}

func GithubLogin(w http.ResponseWriter, r *http.Request) { // redirect to github link
	url := oauthConf.AuthCodeURL("state", oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GithubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"][0]
	//exchange where we can provide context and code and retrieve the token
	token := GithubAccessToken(code)

	name, email, err := GithubData(token)

	if err != nil {
		HandleError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	user := db.GetUserDataByEmail(email) //check user by email
	if user != nil {                              //user already exists
		cookie := dop.GetCookies(*user)
		http.SetCookie(w, cookie)
	} else {
		newUser := structs.User{
			Username: name,
			Email:    email,
		}
		db.InsertUserData(newUser) //add user to DB
		userFromDb := db.GetUserDataByEmail(email)
		if userFromDb != nil { //error when we get user from DB
			cookie := dop.GetCookies(*userFromDb)
			http.SetCookie(w, cookie)
		} else {
			HandleError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GithubAccessToken(code string) string {
	requestBodyMap := map[string]string{
		"client_id":     oauthConf.ClientID,
		"client_secret": oauthConf.ClientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)
	req, reqerr := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}
	respbody, _ := io.ReadAll(resp.Body)
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)
	return ghresp.AccessToken
}

func GithubData(accessToken string) (string, string, error) {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user", nil)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := io.ReadAll(resp.Body)
	req, reqerr = http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	req.Header.Set("Authorization", authorizationHeaderValue)
	resp, resperr = http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	emailResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var githubUser GithubUserInfo
	var githubEmail []GetEmailInfo

	json.Unmarshal(respbody, &githubUser)
	json.Unmarshal(emailResp, &githubEmail)

	return githubUser.Name, githubEmail[0].Email, nil
}