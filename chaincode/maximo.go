/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Product Definitions - The ledger will store marbles and owners
// ============================================================================================================================

// TODO, add
// type Status struct{
//   o INITIALREQUEST
//   o EXEMPTCHECKREQ
//   o HAZARDANALYSISCHECKREQ
//   o CHECKCOMPLETED
// }


// Concept
type Asset struct {
	// ObjectType string        `json:"docType"` //field for couchdb
	// productId       string          `json:"productId"`      //the fieldtags are needed to keep case from bouncing around
  Id       string          `json:"id"`
	Site       string          `json:"site"`
	Meters		 []string						`json:"meters"`
	DateUpdated string				`json:"dateupdated"`
	// DateCreated string				`json:"datecreated"`
	// Temperature 			string
	// Owner      OwnerRelation `json:"owner"`
}

// Participants
// TODO, inheriting from User might be unnecessary
type Meter struct {
	Id				string			`json:"id"`
	Asset			string 			`json:"asset"`
	Readings  []Reading 	`json:"readings"`
	DateUpdated string		`json:"dateupdated"`
}

type WorkOrder struct {
	Id				string			`json:"id"`
	Asset			string 			`json:"asset"`
	LastModifiedBy  string			`json:"lastmodifiedby"`
	Vendor  	string			`json:"vendor"`
	Priority		string			`json:"priority"`
	Status		string			`json:"status"` // trigger api call to maximo when status is complete
	Description string		`json:"description"`
	DateUpdated string		`json:"dateupdated"`
	// RiskAssessment string	`json:"riskassessment"`
}

type User struct {
	Id				string			`json:"id"`
	Company		string			`json:"company"` // city, private_hazmat_co
	Type			string 			`json:"type"` // vendor, inspector
}

type Reading struct {
	Time			string		`json:"time"`
	Value			string		`json:"value"`
	Type			string 		`json:"type"`
	Meter			string 		`json:"meter"` // id
}

// Here as reference for related records
// // ----- Owners ----- //
// type Owner struct {
// 	ObjectType string `json:"docType"`     //field for couchdb
// 	Id         string `json:"id"`
// 	Username   string `json:"username"`
// 	Company    string `json:"company"`
// 	Enabled    bool   `json:"enabled"`     //disabled owners will not be visible to the application
// }
//
// type OwnerRelation struct {
// 	Id         string `json:"id"`
// 	Username   string `json:"username"`    //this is mostly cosmetic/handy, the real relation is by Id not Username
// 	Company    string `json:"company"`     //this is mostly cosmetic/handy, the real relation is by Id not Company
// }

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}


// ============================================================================================================================
// Init - initialize the chaincode
//
// Marbles does not require initialization, so let's run a simple test instead.
//
// Shows off PutState() and how to pass an input argument to chaincode.
// Shows off GetFunctionAndParameters() and GetStringArgs()
// Shows off GetTxID() to get the transaction ID of the proposal
//
// Inputs - Array of strings
//  ["314"]
//
// Returns - shim.Success or error
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Maximo Chaincode Is Starting Up")
	funcName, args := stub.GetFunctionAndParameters()
	var number int
	var err error
	txId := stub.GetTxID()

	fmt.Println("Init() is running")
	fmt.Println("Transaction ID:", txId)
	fmt.Println("  GetFunctionAndParameters() function:", funcName)
	fmt.Println("  GetFunctionAndParameters() args count:", len(args))
	fmt.Println("  GetFunctionAndParameters() args found:", args)

	// expecting 1 arg for instantiate or upgrade
	if len(args) == 1 {
		fmt.Println("  GetFunctionAndParameters() arg[0] length", len(args[0]))

		// expecting arg[0] to be length 0 for upgrade
		if len(args[0]) == 0 {
			fmt.Println("  Uh oh, args[0] is empty...")
		} else {
			fmt.Println("  Great news everyone, args[0] is not empty")

			// convert numeric string to integer
			number, err = strconv.Atoi(args[0])
			if err != nil {
				return shim.Error("Expecting a numeric string argument to Init() for instantiate")
			}

			// this is a very simple test. let's write to the ledger and error out on any errors
			// it's handy to read this right away to verify network is healthy if it wrote the correct value
			err = stub.PutState("selftest", []byte(strconv.Itoa(number)))
			if err != nil {
				return shim.Error(err.Error())                  //self-test fail
			}
		}
	}

	// showing the alternative argument shim function
	alt := stub.GetStringArgs()
	fmt.Println("  GetStringArgs() args count:", len(alt))
	fmt.Println("  GetStringArgs() args found:", alt)

	// store compatible marbles application version
	err = stub.PutState("maximo", []byte("4.0.1"))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Ready for action")                          //self-test pass
	return shim.Success(nil)
}


// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("invoking function - " + function)

	// Handle different functions
	if function == "init" {                    //initialize the chaincode state, used as reset
		return t.Init(stub)
	} else if function == "read" {             //generic read ledger
		return read(stub, args)
	} else if function == "write" {            //generic writes to ledger
		return write(stub, args)
	} else if function == "init_work_order" {      //create a new marble
		return init_work_order(stub, args)
	} else if function == "init_user" {      //create a new marble
		return init_user(stub, args)
	} else if function == "update_work_order" {      //create a new marble
		return update_work_order(stub, args)
	} else if function == "init_asset" {      //create a new marble
		return init_asset(stub, args)
	} else if function == "init_meter" {      //create a new marble
		return init_meter(stub, args)
	} else if function == "add_meter_reading"{        //create a new marble owner
		return add_meter_reading(stub, args)
	} else if function == "read_everything"{   //read everything, (owners + marbles + companies)
		return read_everything(stub)
	} else if function == "getHistory"{        //read history of a marble (audit)
		return getHistory(stub, args)
  }
	// } else if function == "getMarblesByRange"{ //read a bunch of marbles by start and stop id
	// 	return getMarblesByRange(stub, args)
	// } else if function == "disable_owner"{     //disable a marble owner from appearing on the UI
	// 	return disable_owner(stub, args)



  // create product
  // create user (type required...retailer, importer, or supplier)
  // create regulator (not inherited from user)
  // transfer product listing
  // check products
  // update exempted list


	// error out
	fmt.Println("Received unknown invoke function name - " + function)
	return shim.Error("Received unknown invoke function name - '" + function + "'")
}


// ============================================================================================================================
// Query - legacy function
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call - Query()")
}
