package main

import (
  "fmt"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  "github.com/hyperledger/fabric/protos/peer"
)

type ProductChaincode struct {
}

type Product struct {
  ObjectType        string `json:"docType"` // docType is used to distinguish the various types of objects in state database
  ID                string `json:"id"`      // with the field tags (`json:...`), we set the names used in JSON, Go needs upper case
  GTIN              string `json:"gtin"` // this is the barcode, see https://en.wikipedia.org/wiki/Global_Trade_Item_Number
  CreatedBy         string `json:"createdBy"`
  CreatedAt         string `json:"createdAt"`
  UpdatedBy         string `json:"updatedBy"`
  UpdatedAt         string `json:"updatedAt"`
  Producer          string `json:"producer"`
  ContainedProducts []string `json:"containedProducts"`
  Labels            []string `json:"labels"`
  Locale            []ProductLocaleData `json:"locale"`
  Score             Score `json:"score"`
  Status            string `json:"status"`
}
// ex:
// Product{
//   ObjectType: "product",
//   ID: "1",
//   GTIN: "7612100055557",
//   CreatedBy: "user123",
//   CreatedAt: "2018-12-24 12:11:54 UTC",
//   UpdatedBy: "user123",
//   UpdatedAt: "2018-02-10 18:33:39 UTC",
//   Producer: "Wander AG",
//   ContainedProducts: []string{},
//   Labels: []string{"H2892sKSksksdkwops9"},
//   Locale: []ProductLocaleData{
//     ProductLocaleData{
//       Lang: "de",
//       Name: "Ovomaltine crunchy cream — 400 g",
//       Price: "4.99",
//       Currency: "€",
//       Description: "Brotaufstrich mit malzhaltigem Getränkepulver Ovomaltine", Quantity: "400 g",
//       Ingredients: "33% malzhaltiges Getränkepulver: Ovomaltine (Gerstenmalzextrakt, kondensierte Magermilch, kondensiertes Milchserum, fettarmer Kakao, Zucker, Fruktose, Magnesiumcarbonat, Calciumphosphat, Rapsöl, Vitamine [A, E, B1, B2, Pantothensäure, B6, Folsäure, B12, C, Biotin, Niacin], Kochsalz, Aroma Vanillin), Zucker, Pflanzenöle (Raps- und Palmöl), 2.6% Haselnüsse, Calciumphosphat, fettarmer Kakao, Emulgator Sonnenblumenlecithin, Aroma Vanillin.",
//       Packaging: []string{"Glas", "Plastik"},
//       Categories: []string{"Brotaufstriche", "Frühstück", "Nougatcremes"},
//       Image: "products/1/de_1.png",
//       ProductUrl: "http://www.ovomaltine.de/produkte/ovomaltine-crunchy-cream-1/",
//     },
//   },
//   Score: Score{
//     Environment: -34,
//     Climate: -46,
//     Society: -7,
//     Health: -78,
//     Economy: 21,
//   },
//   Status: "active",
// }

type ProductLocaleData struct {
  Lang        string `json:"lang"`
  Name        string `json:"name"`
  Price       string `json:"price"`
  Currency    string `json:"currency"`
  Description string `json:"description"`
  Quantity    string `json:"quantity"`
  Ingredients string `json:"ingredients"`
  Packaging   []string `json:"packaging"`
  Categories  []string `json:"categories"`
  Image       string `json:"image"`
  ProductUrl  string `json:"productUrl"`
}
// ex:
// ProductLocaleData{
//   Lang: "de",
//   Name: "Ovomaltine crunchy cream — 400 g",
//   Price: "4.99",
//   Currency: "€",
//   Description: "Brotaufstrich mit malzhaltigem Getränkepulver Ovomaltine", Quantity: "400 g",
//   Ingredients: "33% malzhaltiges Getränkepulver: Ovomaltine (Gerstenmalzextrakt, kondensierte Magermilch, kondensiertes Milchserum, fettarmer Kakao, Zucker, Fruktose, Magnesiumcarbonat, Calciumphosphat, Rapsöl, Vitamine [A, E, B1, B2, Pantothensäure, B6, Folsäure, B12, C, Biotin, Niacin], Kochsalz, Aroma Vanillin), Zucker, Pflanzenöle (Raps- und Palmöl), 2.6% Haselnüsse, Calciumphosphat, fettarmer Kakao, Emulgator Sonnenblumenlecithin, Aroma Vanillin.",
//   Packaging: []string{"Glas", "Plastik"},
//   Categories: []string{"Brotaufstriche", "Frühstück", "Nougatcremes"},
//   Image: "products/1/de_1.png",
//   ProductUrl: "http://www.ovomaltine.de/produkte/ovomaltine-crunchy-cream-1/",
// }

type Score struct {
  Environment int `json:"environment"`
  Climate     int `json:"climate"`
  Society     int `json:"society"`
  Health      int `json:"health"`
  Economy     int `json:"economy"`
}
// ex:
// Score{
//   Environment: -34,
//   Climate: -46,
//   Society: -7,
//   Health: -78,
//   Economy: 21,
// }

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(ProductChaincode))
	if err != nil {
		fmt.Printf("Error starting Viridian chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *ProductChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *ProductChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initMarble" { //create a new marble
		return t.initMarble(stub, args)
	} else if function == "transferMarble" { //change owner of a specific marble
		return t.transferMarble(stub, args)
	} else if function == "transferMarblesBasedOnColor" { //transfer all marbles of a certain color
		return t.transferMarblesBasedOnColor(stub, args)
	} else if function == "delete" { //delete a marble
		return t.delete(stub, args)
	} else if function == "readMarble" { //read a marble
		return t.readMarble(stub, args)
	} else if function == "queryMarblesByOwner" { //find marbles for owner X using rich query
		return t.queryMarblesByOwner(stub, args)
	} else if function == "queryMarbles" { //find marbles based on an ad hoc rich query
		return t.queryMarbles(stub, args)
	} else if function == "getHistoryForMarble" { //get history of values for a marble
		return t.getHistoryForMarble(stub, args)
	} else if function == "getMarblesByRange" { //get marbles based on range query
		return t.getMarblesByRange(stub, args)
	} else if function == "getMarblesByRangeWithPagination" {
		return t.getMarblesByRangeWithPagination(stub, args)
	} else if function == "queryMarblesWithPagination" {
		return t.queryMarblesWithPagination(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}
