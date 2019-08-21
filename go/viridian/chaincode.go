package viridian

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Chaincode is the one central chaincode that handles all assets
//   (Because we can only start one chaincode per chaincode install.
//    Alternatively, one could install a separate chaincode for each
//    asset.)
// Chaincode is the chaincode that implements the general API
//   methods `Init` and `Invoke`, from which specialised methods are called.
//   The specialized methods belong to other chaincodes like
//   ProductChaincode, ProducerChaincode, etc.
type Chaincode struct {
	Product  *ProductChaincode
	Producer *ProducerChaincode
}

// Init initializes the chaincode
// ==============================
func (c *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	c.Product = new(ProductChaincode)
	c.Producer = new(ProducerChaincode)
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (c *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle the product functions
	if function == "addProduct" { // create a new product
		return c.Product.AddProduct(stub, args)
		// } else if function == "delete" { // delete a product
		// 	return c.delete(stub, args)
		// } else if function == "readProduct" { //read a product
		// 	return c.readProduct(stub, args)
	} else if function == "queryProductsByGTIN" { // find product for GTIN X using rich query
		return c.Product.QueryProductsByGTIN(stub, args)
		// } else if function == "queryProducts" { // find products based on an ad hoc rich query
		// 	return c.queryProducts(stub, args)
		// } else if function == "getHistoryForProduct" { // get history of values for a product
		// 	return c.getHistoryForProduct(stub, args)
		// } else if function == "getMarblesByRange" { //get marbles based on range query
		// 	return c.getMarblesByRange(stub, args)
		// } else if function == "getMarblesByRangeWithPagination" {
		// 	return c.getMarblesByRangeWithPagination(stub, args)
		// } else if function == "queryMarblesWithPagination" {
		// 	return c.queryMarblesWithPagination(stub, args)
	} else if function == "queryProductsByName" {
		return c.Product.QueryProductsByName(stub, args)
	}

	// Handle the producer functions
	if function == "initProducer" {
		return c.Producer.InitProducer(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}
