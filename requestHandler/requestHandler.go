// requestHandler contains logic to handle routes
package requestHandler

import (
	"bytes"
	"fmt"
	"net/http"
)

func HandleGetRequest(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		url := "http://localhost:3001" + r.URL.Path
		var jsonStr = []byte(`{}`)
		req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		err = resp.Write(w)
		if err != nil {
			fmt.Fprintln(w, err)
		}
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
