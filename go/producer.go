package main

import (
	"encoding/json"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/hyperledger/fabric/protos/peer"
)

// ProducerChaincode is the chaincode associated with producers
type ProducerChaincode struct {
}

// Producer is the asset associated with bringing a product to market, so being responsible for it
type Producer struct {
	ScorableAsset
	DocType string   `json:"docType"` // docType is used to distinguish the various types of objects in state database
	Name    string   `json:"name"`
	Address string   `json:"address"` // optional
	URL     string   `json:"url"`     // optional
	Labels  []string `json:"labels"`
}

func (c *ProducerChaincode) initProducer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// Arguments:
	//  0                     1             2                              3                          4
	// Key,                 Name,        Address,                         URL                       Labels
	// "8a259c61-6825-...", "Wander AG", "CH-3176 Neuenegg, Switzerland", "https://www.wander.ch/", []
	var err error
	createdBy, err := cid.GetID(stub)
	if err != nil {
		return shim.Error("Access denied. There is a problem with the client certificate.")
	}
	createdAt := time.Now()
	updatedBy := ""
	updatedAt, _ := time.Parse("2005-12-31", "1776-03-09")
	supersedes := ""
	supersededBy := ""
	changeReason := ""
	score := Score{Environment: 0, Climate: 0, Society: 0, Health: 0, Economy: 0}

	key := args[0]
	name := args[1]
	address := args[2]
	url := args[3]
	var labels []string
	err = json.Unmarshal([]byte(args[4]), &labels)
	if err != nil {
		return shim.Error("5th argument 'labels' must be a string with " +
			"a JSON list of label Keys labelling this producer: [\"label-bd80e824-938c-...\", \"label-127cc795-3a20-...\", ...]" +
			"(or an empty list: [])")
	}

	docType := "producer"
	producer := &Producer{
		ScorableAsset{
			UpdatableAsset{
				ReviewableAsset{createdBy, createdAt, Preliminary},
				updatedBy, updatedAt, supersedes, supersededBy, changeReason},
			score},
		docType, name, address, url, labels}
	jsonAsBytes, err := json.Marshal(producer)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(docType+"-"+key, jsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
