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
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - Array of strings
//    0   ,    1
//   key  ,  value
//  "abc" , "test"
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
//
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs - Array of strings
//      0      ,         1
//     id      ,  authed_by_company
// "m999999999", "united marbles"
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
//
// Shows off building a key's JSON value manually
//
// Inputs - Array of strings

//      0      ,    1  ,     				2  ,      																			3          			 ,    4     ,     5					, 	6
//     id      ,  loan amount ,  borrower_info , 															, state								 , interest ,  balance due	, grade
// "m999999999", "545,000"    ,   object																				deliquent/in payment , 	3.0    , 		520,000			,  BBB
															// credit/income verification/debt to income,

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


// update_asset

// args.add("workorder1")
// args.add("failing")
// args.add("vendor1")
// args.add("asset1")


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
// Init Owner - create a new owner aka end user, store into chaincode state
//
// Shows off building key's value from GoLang Structure
//
// Inputs - Array of Strings
//           0     ,     1   ,   2
//      owner id   , username, company
// "o9999999999999",     bob", "united marbles"
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

/*
func init_asset_listing(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_asset_listing")

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}
	// user.ObjectType = "asset_user"
  asset_listing_id := args[0]
  supplier_id := args[1]
  // asset_ids := args[2] // Expecting JSON array of asset ids. TODO, csv would probably be easier?

  assetListing := ProductListingContract{}
  assetListing.Id = asset_listing_id
  assetListing.Status = "INITIALREQUEST"
  assetListing.Owner = supplier_id
  assetListing.Supplier = supplier_id
  assetListing.OwnerType = "Supplier"

  numProducts := len(args) - 2
  var assets = make([]string, numProducts)
  // all array elements after first 2 are parsed as asset ids
  for i := 2; i < len(args); i++ {
    assets[i - 2] = args[i]
  }
  assetListing.Products = assets //csv.NewReader(asset_ids) //json.Unmarshal(asset_ids)

  // TODO? update asset location to same as supplier
  // supplierAsBytes, err := stub.GetState(supplier_id)
  // supplier := Supplier{}
	// err = json.Unmarshal(supplierAsBytes, &supplier)           //un stringify it aka JSON.parse()
	// if err != nil {
	// 	return shim.Error("Error loading supplier")
	// }
  assetListingAsBytes, _ := json.Marshal(assetListing)                         //convert to array of bytes
  err = stub.PutState(asset_listing_id, assetListingAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store asset listing")
		return shim.Error(err.Error())
	}
	fmt.Println("- end init_asset_listing")
	return shim.Success(nil)
}

func init_regulator(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_regulator")

	//input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}
  // regulator_id := args[0]
  regulator := Regulator{}
  regulator.Id = args[0]
  regulator.countryId = args[1]
  regulatorAsBytes, _ := json.Marshal(regulator)                         //convert to array of bytes
  err = stub.PutState(regulator.Id, regulatorAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store regulator")
		return shim.Error(err.Error())
	}
	fmt.Println("- end init_regulator")
	return shim.Success(nil)
}

func transfer_asset_listing(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  var err error
	fmt.Println("-starting transfer_asset_listing")
  //input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}
  asset_listing_id := args[0]
  new_owner_id := args[1]
  // user_type := args[1]
  // user_id := args[1]
  // retailer_id := args[2]

  assetListingAsBytes, err := stub.GetState(asset_listing_id)
  assetListing := ProductListingContract{}
	err = json.Unmarshal(assetListingAsBytes, &assetListing)           //un stringify it aka JSON.parse()
	if err != nil {
		fmt.Println(string(assetListingAsBytes))
		return shim.Error(err.Error())
	}
  assetListing.Owner = new_owner_id
  if (strings.ToLower(assetListing.OwnerType) == "supplier") {
    assetListing.OwnerType = "Importer"
    assetListing.Status = "EXEMPTCHECKREQ"
    assetListing.Owner = new_owner_id

  } else if ( strings.ToLower(assetListing.OwnerType) == "importer" ) {
    assetListing.OwnerType = "Retailer"
    if ( assetListing.Status == "EXEMPTCHECKREQ" ){
      return shim.Error("Products in listing need to be checked by regulator.")
    } else if ( assetListing.Status == "HAZARDANALYSISCHECKREQ" ){
      return shim.Error("Products cannot be transferred as they've been flagged by regulator.")
    }
    retailerAsBytes, err := stub.GetState(new_owner_id)
    retailer := Retailer{}
  	err = json.Unmarshal(retailerAsBytes, &retailer)           //un stringify it aka JSON.parse()
  	if err != nil {
  		fmt.Println(string(retailerAsBytes))
  		return shim.Error(err.Error())
  	}
    // _, assets := json.Marshal(assetListing.Products)
    for _, asset := range assetListing.Products {
      retailer.Products = append(retailer.Products, asset)
    }
    retailerAsBytes, _ = json.Marshal(retailer)           //convert to array of bytes
    err = stub.PutState(new_owner_id, retailerAsBytes)     //rewrite the marble with id as key
    if err != nil {
      return shim.Error(err.Error())
    }
  } else {
      return shim.Error("Invalid user type provided.")
  }
  assetListingAsBytes, _ = json.Marshal(assetListing)                         //convert to array of bytes
  err = stub.PutState(asset_listing_id, assetListingAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store asset listing")
		return shim.Error(err.Error())
	}
  fmt.Println("- end transfer_asset_listing")
	return shim.Success(nil)
}

func update_exempted_list(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  fmt.Println("- start update_exempted_list")
  // ["regulator1", "org", "org1", "org2"..]
  // add list of exempted orgs or assets to regulator
  regulator_id := args[0]
  exempted_type := args[1]
  // all remaining args are list of ids


  regulatorAsBytes, err := stub.GetState(regulator_id)
  regulator := Regulator{}
	err = json.Unmarshal(regulatorAsBytes, &regulator)           //un stringify it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

  numIds := len(args) - 2
  ids := make([]string, numIds)
  // all array elements after first 2 are parsed as asset ids
  for i := 2; i < len(args); i++ {
    ids[i - 2] = args[i]
  }

  // TODO, should probably append list
  if (exempted_type == "org") {
    regulator.ExemptedOrgIds = ids
  } else if (exempted_type == "asset") {
    regulator.ExemptedProductIds = ids
  }

  regulatorAsBytes, _ = json.Marshal(regulator)                         //convert to array of bytes
  err = stub.PutState(regulator_id, regulatorAsBytes)                    //store owner by its Id
	if err != nil {
		fmt.Println("Could not store regulator")
		return shim.Error(err.Error())
	}
  fmt.Println("- end update_exempted_list")
	return shim.Success(nil)

}

func check_assets(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  asset_listing_id := args[0]
  regulator_id := args[1]

  assetListingAsBytes, err := stub.GetState(asset_listing_id)
  assetListing := ProductListingContract{}
	err = json.Unmarshal(assetListingAsBytes, &assetListing)           //un stringify it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

  supplierAsBytes, err := stub.GetState(assetListing.Supplier)
  supplier := Supplier{}
	err = json.Unmarshal(supplierAsBytes, &supplier)           //un stringify it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

  regulatorAsBytes, err := stub.GetState(regulator_id)
  regulator := Regulator{}
	err = json.Unmarshal(regulatorAsBytes, &regulator)           //un stringify it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

  if (assetListing.Status != "EXEMPTCHECKREQ" && assetListing.Status != "HAZARDANALYSISCHECKREQ"){
    return shim.Error("Invalid state, listing cannot be checked");
  }

  check := true
  if (assetListing.Status=="EXEMPTCHECKREQ"){
    // check if the supplier org is exempted by regulator

    for _, orgId := range regulator.ExemptedOrgIds {
      if ( supplier.orgId == orgId ) {
        check = false
      }
    }
  }

  if (check) {
    assetListing.Status="CHECKCOMPLETED"
  } else {
    assetListing.Status="HAZARDANALYSISCHECKREQ"
  }

  assetListingAsBytes, _ = json.Marshal(assetListing)           //convert to array of bytes
  err = stub.PutState(asset_listing_id, assetListingAsBytes)     //rewrite the marble with id as key
  if err != nil {
    return shim.Error(err.Error())
  }
  fmt.Println("- end check_assets")
	return shim.Success(nil)
}
*/
