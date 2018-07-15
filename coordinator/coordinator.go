package main

import (
	"fmt"
	"log"
	"net/http"
	masterHandler "publicSonarAssignment/src/coordinator/requestHandler"
	slaveHandler "publicSonarAssignment/src/counter/requestHandler"
	"publicSonarAssignment/src/util"
	"regexp"
	"strconv"
	"sync"
)

// /items/{tenantID}/count
var tenantIDURL = regexp.MustCompile(`^/items\/.*\/count$`)

// Counters holds the information about all the Counter Nodes in the system
var Counters = make(map[string]Node)

var maxPortID uint16 = 3000

type CoordinatorHandler struct{}

// ServeHTTP handles all the http requests sent to Coordinator,
// spent 3+ hours on dynamic routing logic for Coordinator and CounterNodes, eventually came up with this :D
func (c CoordinatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestURL := removeSlash(r.URL.Path)
	switch {
	case requestURL == "/items":
		masterHandler.HandlePostRequest(w, r)
	case tenantIDURL.MatchString(requestURL):
		masterHandler.HandleGetRequest(w, r)
	default:
		invalidrequestHandler(w, r)
	}
}

type CounterHandler struct{}

// ServeHTTP handles all the user request,
func (c CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestURL := removeSlash(r.URL.Path)
	switch {
	case requestURL == "/items":
		slaveHandler.HandlePostRequest(w, r)
	case tenantIDURL.MatchString(requestURL):
		slaveHandler.HandleGetRequest(w, r)
	default:
		invalidrequestHandler(w, r)
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

// removeSlash removes the last character of the URL if it is "/"
func removeSlash(input string) string {
	if input[len(input)-1] == '/' {
		input = input[0 : len(input)-1]
	}
	fmt.Println(input)
	return input
}

func invalidrequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, util.ExceptionMessage_InvalidURL)
}

func createCounters(count int) {
	mutex := &sync.Mutex{}
	for i := 0; i < count; i++ {
		GUID := generateGUID()
		portID := generatePortID()

		go func(GUID string, portID uint16) {
			mutex.Lock()
			Counters[GUID] = Node{
				GUID:   GUID,
				PortID: portID,
				Status: Status_WORKING,
			}
			mutex.Unlock()
			counterNode := &http.Server{
				Addr:    ":" + strconv.Itoa(int(portID)),
				Handler: CounterHandler{},
			}
			if err := counterNode.ListenAndServe(); err != nil {
				mutex.Lock()
				Counters[GUID] = Node{
					GUID:   GUID,
					PortID: portID,
					Status: Status_RESIGNED,
				}
				mutex.Unlock()
				log.Fatal(err)
			}
		}(GUID, portID)
	}
}
