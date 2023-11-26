package main

import (
	"encoding/json"
	"fmt"
	"route_chaincode/dijkstra"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}
type outputEvent struct {
	EventName string
}
type RouteMap struct {
	HostID        string
	HostConneting map[string]int
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Printf("init...")
	graph := dijkstra.Graph{
		"a": {"b": 20, "c": 80},
		"b": {"a": 20, "c": 20},
		"c": {"a": 80, "b": 20},
	}
	var routeMap RouteMap

	for k, v := range graph {
		routeMap.HostID = k
		routeMap.HostConneting = v
		value, err := json.Marshal(v)
		if err != nil {
			return shim.Error("fail to json ")

		}
		err1 := stub.PutState(k, []byte(value))
		if err1 != nil {
			return shim.Error("fail to set asset")
		}
	}

	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new Key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else if fn == "path" {
		result, err = path(stub, args)
	} else {

		// assume 'get' even if fn is nil
		result, err = get(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Set stores the asset (both Key and value) on the ledger. If the Key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("incorrect arguments. Expecting a Key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("failed to set asset: %s", args[0])
	}
	event := outputEvent{
		EventName: "set",
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	err = stub.SetEvent("chaincode-event", payload)
	if err != nil {
		fmt.Println("ss")
	}
	return args[1], nil
}

// Get returns the value of the specified asset Key
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("incorrect arguments. Expecting a Key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("asset not found: %s", args[0])
	}
	return string(value), nil
}

func path(stub shim.ChaincodeStubInterface, args []string) (pathres1 string, err error) {
	var g dijkstra.Graph
	var Cost int
	Cost = 0
	fmt.Println(Cost)
	var pathres []string

	start := args[0]
	target := args[1]
	if len(args) != 3 {
		return "", fmt.Errorf("incorrect arguments. Expecting a Key and a value")
	}
	json.Unmarshal([]byte(args[2]), &g)
	if err != nil {
		return "", fmt.Errorf("json to map failed")
	}

	if len(g) == 0 {
		err = fmt.Errorf("cannot find pathres in empty map")
		return "", err
	}

	// ensure start and target are part of the graph
	if _, ok := g[start]; !ok {
		err = fmt.Errorf("cannot find start %v in graph", start)
		return "", err
	}
	if _, ok := g[target]; !ok {
		err = fmt.Errorf("cannot find target %v in graph", target)
		return "", err
	}

	explored := make(map[string]bool)   // set of Nodes we already explored
	frontier := dijkstra.NewQueue()     // queue of the Nodes to explore
	previous := make(map[string]string) // previously visited Node

	// add starting point to the frontier as it'll be the first Node visited
	frontier.Set(start, 0)

	// run until we visited every Node in the frontier
	for !frontier.IsEmpty() {
		// get the Node in the frontier with the lowest Cost (or priority)
		aKey, aPriority := frontier.Next()
		n := dijkstra.Node{Key: aKey, Cost:aPriority}

		// when the Node with the lowest Cost in the frontier is target, we can
		// compute the Cost and pathres and exit the loop
		if n.Key == target {
			Cost = n.Cost
			nKey := n.Key
			for nKey != start {
				pathres = append(pathres, nKey)
				nKey = previous[nKey]
			}

			break
		}

		// add the current Node to the explored set
		explored[n.Key] = true

		// loop all the neighboring Nodes
		for nKey, nCost := range g[n.Key] {
			// skip alreadt-explored Nodes
			if explored[nKey] {
				continue
			}

			// if the Node is not yet in the frontier add it with the Cost
			if _, ok := frontier.Get(nKey); !ok {
				previous[nKey] = n.Key
				frontier.Set(nKey, n.Cost+nCost)
				continue
			}

			frontierCost, _ := frontier.Get(nKey)
			NodeCost := n.Cost + nCost

			// only update the Cost of this Node in the frontier when
			// it's below what's currently set
			if NodeCost < frontierCost {
				previous[nKey] = n.Key
				frontier.Set(nKey, NodeCost)
			}
		}
	}

	// add the origin at the end of the pathres
	pathres = append(pathres, start)

	// reverse the pathres because it was popilated
	// in reverse, form target to start
	for i, j := 0, len(pathres)-1; i < j; i, j = i+1, j-1 {
		pathres[i], pathres[j] = pathres[j], pathres[i]
	}
	pathtag, _ := json.Marshal(pathres)
	pathres1 = string(pathtag) + strconv.Itoa(Cost)

	return pathres1, nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
