package page

import (
	"context"
	"encoding/json"
	"fmt"
	"forum/db"
	"forum/dop"
	"forum/structs"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleRegisterOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8083/registercallback",
		ClientID:     "901065605224-bostfjer52o13jopkhdt43ogau25c6oa.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"}, // permitions we want to ask user
		Endpoint:     google.Endpoint,
	}

	googleLoginOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8083/logincallback",
		ClientID:     "901065605224-bostfjer52o13jopkhdt43ogau25c6oa.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"}, // permitions we want to ask user
		Endpoint:     google.Endpoint,
	}
	randomState = "random"
)

func LoginGoogle(w http.ResponseWriter, r *http.Request) {
	url := googleLoginOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func LoginCallBack(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != randomState {
		fmt.Println("state is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}
	code := r.URL.Query()["code"][0]
	//exchange where we can provide context and code and retrieve the token
	token, err := googleLoginOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("could not get token")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("could not create get request")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	//receive the userinfo in json and decode
	var userData struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		log.Fatal(err)
	}

	user := db.GetUserDataByEmail(userData.Email) //check user by email
	if user != nil {                              //user already exists
		cookie := dop.GetCookies(*user)
		http.SetCookie(w, cookie)
	} else {
		HandleError(w, http.StatusConflict, http.StatusText(http.StatusConflict))
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterGoogle(w http.ResponseWriter, r *http.Request) {
	url := googleRegisterOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func RegisterCallBack(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != randomState {
		fmt.Println("state is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}
	code := r.URL.Query()["code"][0]
	//exchange where we can provide context and code and retrieve the token
	token, err := googleRegisterOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("could not get token")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("could not create get request")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	//receive the userinfo in json and decode
	var userData struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		log.Fatal(err)
	}

	user := db.GetUserDataByEmail(userData.Email) //check user by email

	if user != nil { //user already exists
		HandleError(w, http.StatusCreated, http.StatusText(http.StatusCreated))
		return
	} else {
		newUser := structs.User{
			Username: userData.Name,
			Email:    userData.Email,
		}
		db.InsertUserData(newUser) //add user to DB
		userFromDb := db.GetUserDataByEmail(userData.Email)
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
