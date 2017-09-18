package main

import (
  "encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type CounterChaincode struct {
}

// Counter information
type Counter struct {
  Name string `json:"name"`
  Counts uint64 `json:"counts"`
}

const numOfCounters int = 3

// Initialize counter information
func (cc *CounterChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  var counters [numOfCounters]Counter
  var countersBytes [numOfCounters][]byte

  // Generate counter information
  counters[0] = Counter{Name: "Office Worker", Counts: 0}
  counters[1] = Counter{Name: "Home Worker", Counts: 0}
  counters[2] = Counter{Name: "Student", Counts: 0}

  // Add counter information to world state
  for i := 0; i < len(counters); i++ {
    // Convert to JSON
    countersBytes[i], _ = json.Marshal(counters[i])
    // Add to world state
    stub.PutState(strconv.Itoa(i), countersBytes[i])
  }

  return nil, nil
}

// Update counter information
func (cc *CounterChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  // Handle by function name
  if function == "countUp" {
    // Execute count up
    return cc.countUp(stub, args)
  }

  return nil, errors.New("Received unknown function")
}

// Query counter information
func (cc *CounterChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
  // Handle by function name
  if function == "refresh" {
    // Get counter information
    return cc.getCounters(stub, args)
  }

  return nil, errors.New("Received unknown function")
}

// Execute count up
func (cc *CounterChaincode) countUp(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  // Get selected counter information from world state
  counterId := args[0]
  counterJson, _ := stub.GetState(counterId)

  // Convert JSON format to Counter
  counter := Counter{}
  json.Unmarshal(counterJson, &counter)

  // Count up
  counter.Counts++

  // Add updated value to world state
  counterJson, _ = json.Marshal(counter)
  stub.PutState(counterId, counterJson)

  return nil, nil
}

// Get counter information
func (cc *CounterChaincode) getCounters(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var counters [numOfCounters]Counter
  var countersBytes [numOfCounters][]byte

  for i := 0; i < len(counters); i++ {
    // Get counter information from world state
    countersBytes[i], _ = stub.GetState(strconv.Itoa(i))

    // Convert JSON to Counter
    counters[i] = Counter{}
    json.Unmarshal(countersBytes[i], &counters[i])
  }

  // Convert to JSON
  return json.Marshal(counters)
}

func main() {
  err := shim.Start(new(CounterChaincode))
  if err != nil {
    fmt.Printf("Error starting chaincode: %s", err)
  }
}
