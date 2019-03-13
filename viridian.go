package main

import (
  "fmt"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  "github.com/hyperledger/fabric/protos/peer"
)

type ViridianChaincode struct {
}

type marble struct {
  ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
  Name       string `json:"name"`    //the field tags are needed to keep case from bouncing around
  Color      string `json:"color"`
  Size       int    `json:"size"`
  Owner      string `json:"owner"`
}
