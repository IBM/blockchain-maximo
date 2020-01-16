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
	"bytes"
	"encoding/json"
	"fmt"
  // "reflect"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// ============================================================================================================================
// Read - read a generic variable from ledger
//
// Shows Off GetState() - reading a key/value from the ledger
//
// Inputs - Array of strings
//  0
//  key
//  "abc"
//
// Returns - string
// ============================================================================================================================
func read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, jsonResp string
	var err error
	fmt.Println("starting read")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting key of the var to query")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)           //get the var from ledger
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("Object keys")

	fmt.Println(string(valAsbytes))
	// fmt.Printf("%+v\n", security.Investor)


	fmt.Println("- end read")
	return shim.Success(valAsbytes)                  //send it onward
}

// ============================================================================================================================
// Get everything we need (owners + marbles + companies)
//
// Inputs - none
//
// Returns:
// {
//	"owners": [{
//			"id": "o99999999",
//			"company": "United Marbles"
//			"username": "alice"
//	}],
//	"marbles": [{
//		"id": "m1490898165086",
//		"color": "white",
//		"docType" :"marble",
//		"owner": {
//			"company": "United Marbles"
//			"username": "alice"
//		},
//		"size" : 35
//	}]
// }
// ============================================================================================================================
func read_everything(stub shim.ChaincodeStubInterface) pb.Response {
	type Everything struct {
		// Users   []User   `json:"users"`
		Assets  		[]Asset     `json:"assets"`
		Meters			[]Meter				 `json:"meters"`
		Readings		[]Reading		 `json:"readings"`
		Users				[]User		 `json:"users"`
		WorkOrders	[]WorkOrder	 `json:"workorders"`
	}

	var everything Everything

	// ---- Get All Marbles ---- //
	assetsIterator, err := stub.GetStateByRange("asset0", "asset9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer assetsIterator.Close()

	for assetsIterator.HasNext() {
		aKeyValue, err := assetsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryValAsBytes := aKeyValue.Value
		var asset Asset
		json.Unmarshal(queryValAsBytes, &asset)                  //un stringify it aka JSON.parse()
		everything.Assets = append(everything.Assets, asset)   //add this marble to the list
	}

	usersIterator, err := stub.GetStateByRange("user0", "user9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer usersIterator.Close()

	for usersIterator.HasNext() {
		aKeyValue, err := usersIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryValAsBytes := aKeyValue.Value
		var user User
		json.Unmarshal(queryValAsBytes, &user)                  //un stringify it aka JSON.parse()
		everything.Users = append(everything.Users, user)   //add this marble to the list
	}

	metersIterator, err := stub.GetStateByRange("meter0", "meter9999999999999999999")
	if err != nil {
	  return shim.Error(err.Error())
	}
	defer metersIterator.Close()

	for metersIterator.HasNext() {
	  aKeyValue, err := metersIterator.Next()
	  if err != nil {
	    return shim.Error(err.Error())
	  }
	  queryValAsBytes := aKeyValue.Value
	  var meter Meter
	  json.Unmarshal(queryValAsBytes, &meter)                  //un stringify it aka JSON.parse()
	  everything.Meters = append(everything.Meters, meter)   //add this marble to the list
	}

	workordersIterator, err := stub.GetStateByRange("workorder0", "workorder9999999999999999999")
	if err != nil {
	  return shim.Error(err.Error())
	}
	defer metersIterator.Close()

	for workordersIterator.HasNext() {
	  aKeyValue, err := workordersIterator.Next()
	  if err != nil {
	    return shim.Error(err.Error())
	  }
	  queryValAsBytes := aKeyValue.Value
	  var workorder WorkOrder
	  json.Unmarshal(queryValAsBytes, &workorder)                  //un stringify it aka JSON.parse()
	  everything.WorkOrders = append(everything.WorkOrders, workorder)   //add this marble to the list
	}
	fmt.Println("result", everything)

	//change to array of bytes
	everythingAsBytes, _ := json.Marshal(everything)              //convert to array of bytes
	return shim.Success(everythingAsBytes)
}

// ============================================================================================================================
// Get history of asset
//
// Shows Off GetHistoryForKey() - reading complete history of a key/value
//
// Inputs - Array of strings
//  0
//  id
//  "m01490985296352SjAyM"
// ============================================================================================================================
func getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type AuditHistory struct {
		TxId    string   `json:"txId"`
		Value   WorkOrder   `json:"value"`
	}
	var history []AuditHistory;
	var workorder WorkOrder

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	workOrderId := args[0]
	fmt.Printf("- start getHistoryForWorkOrder: %s\n", workOrderId)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(workOrderId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AuditHistory
		tx.TxId = historyData.TxId                     //copy transaction id over
		json.Unmarshal(historyData.Value, &workorder)     //un stringify it aka JSON.parse()
		if historyData.Value == nil {                  //marble has been deleted
			var emptyWO WorkOrder
			tx.Value = emptyWO                 //copy nil marble
		} else {
			json.Unmarshal(historyData.Value, &workorder) //un stringify it aka JSON.parse()
			tx.Value = workorder                      //copy marble over
		}
		history = append(history, tx)              //add this tx to the list
	}
	fmt.Printf("- getHistoryForWorkOrder returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history)     //convert to array of bytes
	return shim.Success(historyAsBytes)
}

// ============================================================================================================================
// Get history of asset - performs a range query based on the start and end keys provided.
//
// Shows Off GetStateByRange() - reading a multiple key/values from the ledger
//
// Inputs - Array of strings
//       0     ,    1
//   startKey  ,  endKey
//  "marbles1" , "marbles5"
// ============================================================================================================================
func getMarblesByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		aKeyValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryResultKey := aKeyValue.Key
		queryResultValue := aKeyValue.Value

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResultKey)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResultValue))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getMarblesByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
