package main

import (
	"fmt"
	"log"
	"net/http"
	"publicSonarAssignment/src/counter"
	"publicSonarAssignment/src/requestHandler"
	"publicSonarAssignment/src/util"
	"regexp"
	"strconv"
	"sync"
)

// /items/{tenantID}/count
var tenantIDURL = regexp.MustCompile(`^/items\/.*\/count$`)

// Counters holds the information about all the Counter Nodes in the system
var Counters = make(map[string]counter.Node)

var maxPortID uint16 = 3000

func main() {

	// creates desired number of counter nodes at startup of server
	createCounters(3)

	http.HandleFunc("/", HandleURLs)

	fmt.Printf("Starting coordinator server...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// HandleURLs handles all the user request,
// spent 3+ hours on dynamic routing logic and eventually came up with this :D
func HandleURLs(w http.ResponseWriter, r *http.Request) {
	requestURL := removeSlash(r.URL.Path)
	switch {
	case requestURL == "/items":
		requestHandler.HandlePostRequest(w, r)
	case tenantIDURL.MatchString(requestURL):
		requestHandler.HandleGetRequest(w, r)
	default:
		invalidrequestHandler(w, r)
	}
}

func invalidrequestHandler(w http.ResponseWriter, r *http.Request) {
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

func createCounters(count int) {
	mutex := &sync.Mutex{}
	for i := 0; i < count; i++ {
		GUID := generateGUID()
		portID := generatePortID()

		go func(GUID string, portID uint16) {
			mutex.Lock()
			Counters[GUID] = counter.Node{
				GUID:   GUID,
				PortID: portID,
				Status: counter.Status_WORKING,
			}
			mutex.Unlock()
			if err := http.ListenAndServe(":"+strconv.Itoa(int(portID)), nil); err != nil {
				mutex.Lock()
				Counters[GUID] = counter.Node{
					GUID:   GUID,
					PortID: portID,
					Status: counter.Status_RESIGNED,
				}
				mutex.Unlock()
				log.Fatal(err)
			}
		}(GUID, portID)
	}
}

func generateGUID() string {
	return "CounterNode1"
}

func generatePortID() uint16 {
	maxPortID++
	return maxPortID
}
