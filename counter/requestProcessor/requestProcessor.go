package requestProcessor

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
		fmt.Fprintf(w, "invalid URL for GET request, try http://localhost:8080/items/{tenantID}/count")
	}
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "in HandlePostRequest\n")

	switch r.Method {
	case "POST":
		fmt.Fprintf(w, "POST successful in HandlePostRequest\n")
	default:
		fmt.Fprintf(w, "invalid URL for POST request, try http://localhost:8080/items/")
	}
}
