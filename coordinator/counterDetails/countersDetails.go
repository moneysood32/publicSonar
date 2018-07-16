package counterDetails

// Status provides the current status of a CounterNode
type Status uint16

const (
	// Status_WORKING is set when counter node is available to receive requests from user.
	Status_WORKING Status = iota

	// Status_UNAVAILABLE is set when counter node temporarily stops working
	Status_UNAVAILABLE

	// Status_RESIGNED is set when counter node permanently halts
	// and we can safely reclaim its resources
	Status_RESIGNED
)

// Node contains information about a Counter Node (Counter Server)
type Node struct {
	GUID   string
	PortID uint16
	Status Status
}

// Counters holds the information about all the Counter Nodes in the system, map[GUID]Node
var Counters = make(map[string]Node)

var CurrentPortID uint16 = 3000
var MaxPortID uint16 = 10000

// CurrentRequests holds the number of requests currently sent to each node map[portID]requestCount
var CurrentRequests = make(map[string]int)
