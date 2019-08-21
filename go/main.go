package main

import (
	"fmt"

	"github.com/chaincode/viridian/go/viridian"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(viridian.Chaincode))
	if err != nil {
		fmt.Printf("Error starting Product chaincode: %s", err)
	}
}
