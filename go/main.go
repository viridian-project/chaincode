package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(ProductChaincode))
	if err != nil {
		fmt.Printf("Error starting Product chaincode: %s", err)
	}
	err = shim.Start(new(ProducerChaincode))
	if err != nil {
		fmt.Printf("Error starting Producer chaincode: %s", err)
	}
}
