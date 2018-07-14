package main

import (
	"fmt"
	"log"
	"net/http"
	"publicSonarAssignment/src/handler"
	"publicSonarAssignment/src/util"
	"regexp"
)

// /items/{tenantID}/count
var tenantIDURL = regexp.MustCompile(`^/items\/.*\/count$`)

func main() {
	http.HandleFunc("/", HandleURLs)
	fmt.Printf("Starting coordinator server...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// HandleURLs handles all the user request,
// spent 3+ hours on dynamic routing logic and eventually came up with this :P
func HandleURLs(w http.ResponseWriter, r *http.Request) {
	requestURL := removeSlash(r.URL.Path)
	switch {
	case requestURL == "/items":
		handler.HandlePostRequest(w, r)
	case tenantIDURL.MatchString(requestURL):
		handler.HandleGetRequest(w, r)
	default:
		invalidHandler(w, r)
	}
}

func invalidHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, util.ExceptionMessage_InvalidURL)
}

// removeSlash removes the last character of the URL if it is "/"
func removeSlash(input string) string {
	if input[len(input)-1] == '/' {
		input = input[0 : len(input)-1]
	}
	fmt.Println(input)
	return input
}
