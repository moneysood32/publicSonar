package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"publicSonarAssignment/src/coordinator/counterDetails"
	masterHandler "publicSonarAssignment/src/coordinator/requestHandler"
	slaveHandler "publicSonarAssignment/src/counter/requestHandler"
	"publicSonarAssignment/src/util"
	"regexp"
	"strconv"
	"sync"
)

// /items/{tenant}/count
var tenantURL = regexp.MustCompile(`^/items\/.*\/count$`)

type CoordinatorHandler struct{}

// ServeHTTP handles all the http requests sent to Coordinator,
// spent 3+ hours on dynamic routing logic for Coordinator and CounterNodes, eventually came up with this :D
func (c CoordinatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestURL := removeSlash(r.URL.Path)
	switch {
	case requestURL == "/items":
		masterHandler.HandlePostRequest(w, r)
	case tenantURL.MatchString(requestURL):
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
	case tenantURL.MatchString(requestURL):
		slaveHandler.HandleGetRequest(w, r)
	default:
		invalidrequestHandler(w, r)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateGUID() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func generatePortID() uint16 {
	counterDetails.CurrentPortID++
	fmt.Println(counterDetails.CurrentPortID)
	if counterDetails.CurrentPortID == counterDetails.MaxPortID {
		fmt.Println("PortID has reached max limit")
	}
	return counterDetails.CurrentPortID
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
			counterDetails.Counters[GUID] = counterDetails.Node{
				GUID:   GUID,
				PortID: portID,
				Status: counterDetails.Status_WORKING,
			}
			mutex.Unlock()
			counterNode := &http.Server{
				Addr:    ":" + strconv.Itoa(int(portID)),
				Handler: CounterHandler{},
			}
			counterDetails.CurrentRequests[strconv.Itoa(int(portID))] = 0
			if err := counterNode.ListenAndServe(); err != nil {
				mutex.Lock()
				counterDetails.Counters[GUID] = counterDetails.Node{
					GUID:   GUID,
					PortID: portID,
					Status: counterDetails.Status_RESIGNED,
				}
				mutex.Unlock()
				log.Fatal(err)
			}
		}(GUID, portID)
	}
}
