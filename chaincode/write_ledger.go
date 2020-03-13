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
	"encoding/json"
  // "encoding/csv"
	"fmt"
	// "strconv"
	// "strings"
	// "reflect"
	// "math"
	"math/rand"
	"time"
	"strconv"
	// "encoding/gob"
	// "bytes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// ============================================================================================================================
// write() - genric write variable into ledger
// ============================================================================================================================
func write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	var err error
	fmt.Println("starting write")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2. key of the variable and value to set")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	key = args[0]                                   //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value))         //write the variable into the ledger
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end write")
	return shim.Success(nil)
}

// ============================================================================================================================
// delete_marble() - remove a marble from state and from marble index
// ============================================================================================================================
func delete(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	fmt.Println("starting delete")

	id := args[0]

	err := stub.DelState(id)                                                 //remove the key from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Println("- end delete")
	return shim.Success(nil)
}

// ============================================================================================================================
// Init Product - create a new asset, store into chaincode state
// ============================================================================================================================
func init_asset(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
	fmt.Println("starting init_asset")
	var asset Asset
	asset.Id = args[0]
	// asset.Site =  args[1]
	// asset.CountryId = args[2]
	// check if asset already exists
	// TODO, uncomment
	// _, err = get_asset(stub, asset.Id)
	// if err == nil {
	// 	fmt.Println("This asset already exists - " + asset.Id)
	// 	return shim.Error("This asset already exists - " + asset.Id)
	// }
	//store asset
	assetAsBytes, _ := json.Marshal(asset)                         //convert to array of bytes
	fmt.Println("writing asset to ledger state")
	fmt.Println(string(assetAsBytes))
	err = stub.PutState(asset.Id, assetAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store asset")
		return shim.Error(err.Error())
	}
	fmt.Println("- end init_asset")

	return shim.Success(nil)
}

// this should be called by automation scripts from maximo
func init_work_order(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
	fmt.Println("starting init_work_order")
	var workorder WorkOrder
	// json.Marshal(args)
	workorder.Id = args[0]
	workorder.Status = args[1]
	workorder.LastModifiedBy = "Maximo" // TODO, perhaps add actual name/maximo_user?
	workorder.Vendor = args[2] // TODO, lookup a way to prevent duplicate vendors
	if (len(args) > 3) {
		workorder.Asset = args[3]
	} else {
		workorder.Asset = ""
	}
	workOrderAsBytes, _ := json.Marshal(workorder)                         //convert to array of bytes
	fmt.Println("writing workorder to ledger state")
	fmt.Println(string(workOrderAsBytes))
	err = stub.PutState(workorder.Id, workOrderAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store workorder")
		return shim.Error(err.Error())
	}

	// register vendor
	fmt.Println("registering wo vendor")

	var user User
  // allow user to provide id from UI instead of randomly generating one TODO
	rand.Seed(time.Now().UnixNano())
	id_num := strconv.Itoa(rand.Intn(10000))
	user.Id = "user" + id_num
	user.Company = workorder.Vendor
	user.Type = "HAZMAT_VENDOR"
	fmt.Println("user info")
	fmt.Println(user.Id)
	fmt.Println(user.Company)
	fmt.Println(user.Type)
	userAsBytes, _ := json.Marshal(user)                         //convert to array of bytes
  err = stub.PutState(user.Id, userAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store user")
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_work_order")
	return shim.Success(nil)
}

// called by third-party UI
func update_work_order(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
	fmt.Println("starting update_work_order")
	workorder_id := args[0]
	workorderAsBytes, err := stub.GetState(workorder_id)
	fmt.Println("work_order loaded")
	workorder := WorkOrder{}
	err = json.Unmarshal(workorderAsBytes, &workorder)           //un stringify it aka JSON.parse()
	if err != nil {
		return shim.Error("Error loading workorder")
	}
	workorder.Status = args[1] // INPRG, APPR, COMPLETED, WAPPR (should be called after)
	workorder.LastModifiedBy = args[2] // user_id --- HAZMAT_INSPECTOR, HAZMAT_VENDOR, BUILDING_INSPECTOR

	if (len(args) > 3 ) {
		workorder.Priority = args[3]
	}
	if (workorder.Status == "WAPPR" ) {
		workorder.Status = "APPR"
	} else if (workorder.Status == "APPR" ) {
		workorder.Status = "INPRG"
	} else if (workorder.Status == "INPRG" ) {
		workorder.Status = "COMP"
	}
	workorder.DateUpdated = time.Now().Format("2006-01-02 15:04:05")
	workOrderAsBytes, _ := json.Marshal(workorder)                         //convert to array of bytes
	err = stub.PutState(workorder.Id, workOrderAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store updated workorder")
		return shim.Error(err.Error())
	}
	fmt.Println("- end update_work_order")
	return shim.Success(nil)
}
// ============================================================================================================================
// Init Meter - register a new meter, store into chaincode state
// ============================================================================================================================

func init_meter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_meter")

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

  // TODO, should probably allow JSON object as third arg
	// meter1, asset1
  id := args[0]
  assetId := args[1]

  var meter Meter

	// reading type (F, Volts, etc)
	// t := args[2]
	// meter.Type = t

	meter.Id = id
  meter.Asset = assetId

	meterAsBytes, _ := json.Marshal(meter)                         //convert to array of bytes
	err = stub.PutState(id, meterAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store meter")
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_meter")
	return shim.Success(nil)
}


func add_meter_reading(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// add reading to meter
	// meter1, 50
	meter_id := args[0]
	meter_reading := args[1]
	meterAsBytes, err := stub.GetState(meter_id)
  meter := Meter{}
	err = json.Unmarshal(meterAsBytes, &meter)           //un stringify it aka JSON.parse()

	reading := Reading{}
	reading.Value = meter_reading
	reading.Time = time.Now().Format("2006-01-02 15:04:05")
	reading.Meter = meter_id
	// reading.Meter =
	// readingAsBytes, _ = json.Marshal(meter)                         //convert to array of bytes

	// append reading
	meter.Readings = append(meter.Readings, reading)

	meterAsBytes, _ = json.Marshal(meter)                         //convert to array of bytes
  err = stub.PutState(meter_id, meterAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store updated meter reading")
		return shim.Error(err.Error())
	}
	fmt.Println("- end add_meter_reading")
	return shim.Success(nil)
}


func init_user(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_user")

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}
  // regulator_id := args[0]
  user := User{}
  user.Id = args[0]
  user.Company = args[1]
	user.Type = args[2]
  userAsBytes, _ := json.Marshal(user)                         //convert to array of bytes
  err = stub.PutState(user.Id, userAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store user")
		return shim.Error(err.Error())
	}
	fmt.Println("- end init_user")
	return shim.Success(nil)
}
