package requestProcessor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"publicSonarAssignment/src/counter"
)

func HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "in HandleGetRequest")

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "GET successful in HandleGetRequest")
	default:
		fmt.Fprintf(w, "invalid URL for GET request, try http://localhost:3001/items/{tenantID}/count")
	}
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "in HandlePostRequest\n")
	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		userInfo, err := counter.ParseAndSave(body)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(w, userInfo)
	default:
		fmt.Fprintf(w, "invalid URL for POST request, try http://localhost:3001/items/")
	}
}
