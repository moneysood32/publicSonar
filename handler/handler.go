// handler contains logic to handle routes
package handler

import (
	"fmt"
	"net/http"
)

func HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "in HandleGetRequest")

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "GET successful in HandleGetRequest")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "in HandlePostRequest")

	switch r.Method {
	case "POST":
		fmt.Fprintf(w, "POST successful in HandlePostRequest")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
