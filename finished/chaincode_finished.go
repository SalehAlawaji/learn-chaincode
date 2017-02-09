/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//==============================================================================================================================
//	Vehicle - Defines the structure for a car object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Vehicle struct {
	Make            string `json:"make"`
	Model           string `json:"model"`
	Reg             string `json:"reg"`
	VIN             int    `json:"VIN"`
	Owner           string `json:"owner"`
	Scrapped        bool   `json:"scrapped"`
	Status          int    `json:"status"`
	Colour          string `json:"colour"`
	V5cID           string `json:"v5cID"`
	LeaseContractID string `json:"leaseContractID"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

//=================================================================================================================================
//	 Create Function
//=================================================================================================================================
//	 Create Vehicle - Creates the initial JSON for the vehcile and then saves it to the ledger.
//=================================================================================================================================
func (t *SimpleChaincode) create_vehicle(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, v5cID string) ([]byte, error) {
	var v Vehicle

	v5c_ID := "\"v5cID\":\"" + v5cID + "\", " // Variables to define the JSON
	vin := "\"VIN\":0, "
	make := "\"Make\":\"UNDEFINED\", "
	model := "\"Model\":\"UNDEFINED\", "
	reg := "\"Reg\":\"UNDEFINED\", "
	owner := "\"Owner\":\"dksjnf\", "
	colour := "\"Colour\":\"UNDEFINED\", "
	leaseContract := "\"LeaseContractID\":\"UNDEFINED\", "
	status := "\"Status\":0, "
	scrapped := "\"Scrapped\":false"

	vehicle_json := "{" + v5c_ID + vin + make + model + reg + owner + colour + leaseContract + status + scrapped + "}" // Concatenates the variables to create the total JSON object

	//matched, err := regexp.Match("^[A-z][A-z][0-9]{7}", []byte(v5cID)) // matched = true if the v5cID passed fits format of two letters followed by seven digits

	/*if err != nil {
		fmt.Printf("CREATE_VEHICLE: Invalid v5cID: %s", err)
		return nil, errors.New("Invalid v5cID")
	}*/

	/*if v5c_ID == "" ||
		matched == false {
		fmt.Printf("CREATE_VEHICLE: Invalid v5cID provided")
		return nil, errors.New("Invalid v5cID provided")
	}*/

	err := json.Unmarshal([]byte(vehicle_json), &v) // Convert the JSON defined above into a vehicle object for go

	if err != nil {
		return nil, errors.New("Invalid JSON object")
	}

	record, err := stub.GetState(v.V5cID) // If not an error then a record exists so cant create a new car with this V5cID as it must be unique

	if record != nil {
		return nil, errors.New("Vehicle already exists")
	}

	/*if caller_affiliation != AUTHORITY { // Only the regulator can create a new v5c

		return nil, errors.New(fmt.Sprintf("Permission Denied. create_vehicle. %v === %v", caller_affiliation, AUTHORITY))

	}*/

	_, err = t.saveChanges(stub, v)

	if err != nil {
		fmt.Printf("CREATE_VEHICLE: Error saving changes: %s", err)
		return nil, errors.New("Error saving changes")
	}

	bytes, err := stub.GetState("v5cIDs")

	if err != nil {
		return nil, errors.New("Unable to get v5cIDs")
	}

	/*var v5cIDs V5C_Holder*/

	/*err = json.Unmarshal(bytes, &v5cIDs)

	if err != nil {
		return nil, errors.New("Corrupt V5C_Holder record")
	}

	v5cIDs.V5Cs = append(v5cIDs.V5Cs, v5cID)

	bytes, err = json.Marshal(v5cIDs)*/

	if err != nil {
		fmt.Print("Error creating V5C_Holder record")
	}

	err = stub.PutState("v5cIDs", bytes)

	if err != nil {
		return nil, errors.New("Unable to put the state")
	}

	return nil, nil

}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. name of the key and value to set")
	}
	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "create_car" {
		//return t.create_vehicle(stub, caller, caller_affiliation, args[0])

		return t.create_vehicle(stub, "", "", args[0])
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. name of the key and value to set")
	}
	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	if function == "get_vehicle_details" {

		return nil, errors.New("fuck you : " + function)
		v, err := t.retrieve_v5c(stub, args[0])
		if err != nil {
			fmt.Printf("QUERY: Error retrieving v5c: %s", err)
			return nil, errors.New("QUERY: Error retrieving v5c " + err.Error())
		}
		return t.get_vehicle_details(stub, v)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *SimpleChaincode) saveChanges(stub shim.ChaincodeStubInterface, v Vehicle) (bool, error) {

	bytes, err := json.Marshal(v)

	if err != nil {
		fmt.Printf("SAVE_CHANGES: Error converting vehicle record: %s", err)
		return false, errors.New("Error converting vehicle record")
	}

	err = stub.PutState(v.V5cID, bytes)

	if err != nil {
		fmt.Printf("SAVE_CHANGES: Error storing vehicle record: %s", err)
		return false, errors.New("Error storing vehicle record")
	}

	return true, nil
}

func (t *SimpleChaincode) retrieve_v5c(stub shim.ChaincodeStubInterface, v5cID string) (Vehicle, error) {

	var v Vehicle

	bytes, err := stub.GetState(v5cID)

	if err != nil {
		fmt.Printf("RETRIEVE_V5C: Failed to invoke vehicle_code: %s", err)
		return v, errors.New("RETRIEVE_V5C: Error retrieving vehicle with v5cID = " + v5cID)
	}

	err = json.Unmarshal(bytes, &v)

	if err != nil {
		fmt.Printf("RETRIEVE_V5C: Corrupt vehicle record "+string(bytes)+": %s", err)
		return v, errors.New("RETRIEVE_V5C: Corrupt vehicle record" + string(bytes))
	}

	return v, nil
}

//=================================================================================================================================
//	 Read Functions
//=================================================================================================================================
//	 get_vehicle_details
//=================================================================================================================================
func (t *SimpleChaincode) get_vehicle_details(stub shim.ChaincodeStubInterface, v Vehicle) ([]byte, error) {

	bytes, err := json.Marshal(v)

	if err != nil {
		return nil, errors.New("GET_VEHICLE_DETAILS: Invalid vehicle object")
	}
	return bytes, nil

}
