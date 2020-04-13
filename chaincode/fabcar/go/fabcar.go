/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a did
type SmartContract struct {
	contractapi.Contract
}

// Did describes basic details of what makes up a did document
type Did struct {
	Id                          string `json:"id"`
	AuthenticationId            string `json:"authenticationId"`
	AuthenticationType          string `json:"authenticationType"`
	AuthenticationController    string `json:"authenticationController"`
	AuthenticationPublicKeyPerm string `json:"authenticationPublicKeyPerm"`
	ServiceId                   string `json:"serviceId"`
	ServiceType                 string `json:"serviceType"`
	ServiceEndPoint             string `json:"serviceEndPoint"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Did
}

// InitLedger adds a base set of dids to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	dids := []Did{
		Did{Id: "did:example:12346789abcdefghi", AuthenticationId: "did:example:12346789abcdefghi#keys-1",
			AuthenticationType: "RsaVerificationKey2018", AuthenticationController: "did:example:12346789abcdefghi",
			AuthenticationPublicKeyPerm: "-----BEGIN PUBLIC KEY...END PUBLIC KEY-----\r\n",
			ServiceId:                   "did:example:12346789abcdefghi#vcs", ServiceType: "VerifiableCredentialService",
			ServiceEndPoint: "https://example.com/vc/"},

		Did{Id: "did:example:12346789asdfghjkl", AuthenticationId: "did:example:12346789asdfghjkl#keys-1",
			AuthenticationType: "RsaVerificationKey2018", AuthenticationController: "did:example:12346789asdfghjkl",
			AuthenticationPublicKeyPerm: "-----BEGIN PUBLIC KEY...END PUBLIC KEY-----\r\n",
			ServiceId:                   "did:example:12346789aasdfghjkl#vcs", ServiceType: "VerifiableCredentialService",
			ServiceEndPoint: "https://example2.com/vc/"},
	}

	for i, did := range dids {
		didAsBytes, _ := json.Marshal(did)
		err := ctx.GetStub().PutState("DID"+strconv.Itoa(i), didAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateDid adds a new did to the world state with given details
func (s *SmartContract) CreateDid(ctx contractapi.TransactionContextInterface, didNumber string, id string, authenticationId string, authenticationType string,
	authenticationController string, authenticationPublicKeyPerm string, serviceId string, serviceType string, serviceEndPoint string) error {
	did := Did{
		Id:                          id,
		AuthenticationId:            authenticationId,
		AuthenticationType:          authenticationType,
		AuthenticationController:    authenticationController,
		AuthenticationPublicKeyPerm: authenticationPublicKeyPerm,
		ServiceId:                   serviceId,
		ServiceType:                 serviceType,
		ServiceEndPoint:             serviceEndPoint,
	}

	didAsBytes, _ := json.Marshal(did)

	return ctx.GetStub().PutState(didNumber, didAsBytes)
}

// QueryDidByKey returns the did stored in the world state with given key
func (s *SmartContract) QueryDidByKey(ctx contractapi.TransactionContextInterface, didNumber string) (*Did, error) {
	didAsBytes, err := ctx.GetStub().GetState(didNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if didAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", didNumber)
	}

	did := new(Did)
	_ = json.Unmarshal(didAsBytes, did)

	return did, nil
}

// QueryDidById returns the did stored in the world state with given id

func (s *SmartContract) QueryDidById(ctx contractapi.TransactionContextInterface, id string) (*Did, error) {
	startKey := "DID0"
	endKey := "DID99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		did := new(Did)
		_ = json.Unmarshal(queryResponse.Value, did)

		if did.Id == id {
			return did, nil
		}

		queryResult := QueryResult{Key: queryResponse.Key, Record: did}
		results = append(results, queryResult)
	}

	return nil, fmt.Errorf("%s does not exist", id)
}

// QueryAllDids returns all did documents found in world state
func (s *SmartContract) QueryAllDids(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "DID0"
	endKey := "DID99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		did := new(Did)
		_ = json.Unmarshal(queryResponse.Value, did)

		queryResult := QueryResult{Key: queryResponse.Key, Record: did}
		results = append(results, queryResult)
	}

	return results, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
