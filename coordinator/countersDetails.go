package main

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
