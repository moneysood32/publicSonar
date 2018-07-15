package main

import (
	"fmt"
	"log"
	"net/http"
	"publicSonarAssignment/src/counter"
	"publicSonarAssignment/src/counter/requestProcessor"
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

type CoordinatorHandler struct{}
type CounterHandler struct{}

func main() {

	// creates desired number of counter nodes at startup of server
	createCounters(3)

	coordinatorServer := &http.Server{
		Addr:    ":8080",
		Handler: CoordinatorHandler{},
	}
	fmt.Printf("Starting coordinator server...\n")
	log.Fatal(coordinatorServer.ListenAndServe())

}

// ServeHTTP handles all the http requests sent to Coordinator,
// spent 3+ hours on dynamic routing logic for Coordinator and CounterNodes, eventually came up with this :D
func (c CoordinatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

// ServeHTTP handles all the user request,
func (c CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestURL := removeSlash(r.URL.Path)
	switch {
	case requestURL == "/items":
		requestProcessor.HandlePostRequest(w, r)
	case tenantIDURL.MatchString(requestURL):
		requestProcessor.HandleGetRequest(w, r)
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
			counterNode := &http.Server{
				Addr:    ":" + strconv.Itoa(int(portID)),
				Handler: CounterHandler{},
			}
			if err := counterNode.ListenAndServe(); err != nil {
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
	fmt.Println(maxPortID)
	return maxPortID
}
