package dop

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"forum/db"
	"forum/structs"
	"net/http"
	"time"
)

// Create a new cookie for the user session and store it in the database
func GetCookies(user structs.User) *http.Cookie {
	sessionID, _ := GetSessionID()
	expires := time.Now().Add(96 * time.Hour)
	cookie := &http.Cookie{
		Name:     "session",
		Value:    fmt.Sprintf("%d", sessionID),
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
	}
	db.RemoveCookieData(user.ID)
	var r structs.Session
	r.UserID = user.ID
	r.Expires = expires
	r.ID = sessionID
	db.InsertCookieData(r)
	return cookie
}

// Generate a random session ID
func GetSessionID() (int, error) {
	r := make([]byte, 4)
	_, err := rand.Read(r)
	if err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint32(r)), nil
}

// Hash the password passed, and returns it as a string encoded in hexadecimal
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// Check if the password passed matches the hash value stored in the database
func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}

// Check the integrity of the cookie associated with the user session. 
// If the session is invalid, the function deletes the cookie
func CheckCookieIntegrity(w http.ResponseWriter, r *http.Request) *structs.User {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}
	session := db.GetSessionData(sessionCookie.Value)
	if session == nil {
		return nil
	}
	if session.Expires.After(time.Now()) {
		c := &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(w, c)
		return nil
	}
	user := db.GetUserDataById(session.UserID)
	return user
}

// Return user data associated with the current user session, if there is one
func GetUserInCookie(r *http.Request) *structs.User {
	cookie, err := r.Cookie("session")
	if err == nil {
		session := db.GetSessionData(cookie.Value)
		if session == nil {
			return nil
		}
		return db.GetUserDataById(session.UserID)
	}
	return nil
}