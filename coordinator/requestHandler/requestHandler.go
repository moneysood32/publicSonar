// Package requestHandler contains logic to handle coordinator routes
package requestHandler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HandleGetRequest(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		url := "http://localhost:3001" + r.URL.Path
		req, err := http.NewRequest("GET", url, nil)

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
		fmt.Fprintf(w, "invalid URL for GET request, try http://localhost:8080/items/{tenant}/count")
	}
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		url := "http://localhost:3001" + r.URL.Path

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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
		fmt.Fprintf(w, "invalid URL for POST request, try http://localhost:8080/items/")
	}
}
