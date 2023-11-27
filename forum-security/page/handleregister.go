package page

import (
	"forum/db"
	"forum/helpers"
	"forum/structs"
	"net/http"
	"text/template"
	"time"
)

// Handle new user registration requests
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	tmpl,_ := template.ParseFiles("templates/register.html")
	if r.Method == "GET" {
		if r.URL.Path != "/register" {
			HandleError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		tmpl.Execute(w, "")

	} else if r.Method == "POST" {

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if db.GetUserDataByEmail(email) != nil {
			errormes := "Email is already in use"
			tmpl.Execute(w, errormes)
			return
		}

		if db.GetUserDataByUsername(username) != nil {
			errormes := "Username already in use"
			tmpl.Execute(w, errormes)
			return
		}

		hashedPassword := helpers.HashPassword(password)

		newUser := structs.User{
			Username:  username,
			Email:     email,
			Password:  hashedPassword,
			CreatedAt: time.Now(),
		}

		db.InsertUserData(newUser)
		userFromDb := db.GetUserDataByEmail(email)
		if userFromDb != nil && helpers.CheckPasswordHash(password, userFromDb.Password) {
			cookie := helpers.GetCookies(*userFromDb)
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			HandleError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}
	}
}