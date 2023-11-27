package page

import (
	"fmt"
	"net/http"
	"text/template"
)

type errorResponse struct {
	Code    int
	Message string
}

// Handles the error, generates an error message and sends it to the user
func HandleError(w http.ResponseWriter, statusCode int, Massage string) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, fmt.Sprint("Error parsing:", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	var errResp errorResponse
	errResp.Message = Massage
	errResp.Code = statusCode
	err = tmpl.Execute(w, errResp)
	if err != nil {
		http.Error(w, fmt.Sprint("Error executing template:", err), http.StatusInternalServerError)
	}
}

