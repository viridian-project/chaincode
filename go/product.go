package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type ProductChaincode struct {
}

type ProductLocaleData struct {
	Lang        string   `json:"lang"`        // regex=/^[a-z]{2}$/ // ISO language code according to https://en.wikipedia.org/wiki/ISO_639-1, there should be only one locale data for each language
	Name        string   `json:"name"`        // product 'short name'
	Price       string   `json:"price"`       // optional
	Currency    string   `json:"currency"`    // optional
	Description string   `json:"description"` // optional
	Quantities  []string `json:"quantities"`
	Ingredients string   `json:"ingredients"` // optional
	Packaging   []string `json:"packaging"`
	Categories  []string `json:"categories"`
	ImageURL    string   `json:"imageUrl"` // optional
	URL         string   `json:"url"`      // optional
}

// ex:
// &ProductLocaleData{
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

type Product struct {
	ScorableAsset
	DocType           string              `json:"docType"` // docType is used to distinguish the various types of objects in state database
	GTIN              string              `json:"gtin"`    // optional
	Producer          string              `json:"producer"`
	ContainedProducts []string            `json:"containedProducts"`
	Labels            []string            `json:"labels"`
	Locales           []ProductLocaleData `json:"locales"`
}

// ex:
// &Product{
//   ObjectType: "product",
//   ID: 1,
//   GTIN: "7612100055557",
//   CreatedBy: "user123",
//   CreatedAt: "2018-12-24 12:11:54 UTC",
//   UpdatedBy: "user123",
//   UpdatedAt: "2018-02-10 18:33:39 UTC",
//   Producer: "Wander AG",
//   ContainedProducts: []string{},
//   Labels: []string{"H2892sKSksksdkwops9"},
//   Locale: []ProductLocaleData{
//     &ProductLocaleData{
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
//   Score: &Score{
//     Environment: -34,
//     Climate: -46,
//     Society: -7,
//     Health: -78,
//     Economy: 21,
//   },
//   Status: "active",
// }

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(ProductChaincode))
	if err != nil {
		fmt.Printf("Error starting Product chaincode: %s", err)
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
	if function == "initProduct" { // create a new product
		return t.initProduct(stub, args)
		// } else if function == "delete" { // delete a product
		// 	return t.delete(stub, args)
		// } else if function == "readProduct" { //read a product
		// 	return t.readProduct(stub, args)
	} else if function == "queryProductsByGTIN" { // find product for GTIN X using rich query
		return t.queryProductsByGTIN(stub, args)
		// } else if function == "queryProducts" { // find products based on an ad hoc rich query
		// 	return t.queryProducts(stub, args)
		// } else if function == "getHistoryForProduct" { // get history of values for a product
		// 	return t.getHistoryForProduct(stub, args)
		// } else if function == "getMarblesByRange" { //get marbles based on range query
		// 	return t.getMarblesByRange(stub, args)
		// } else if function == "getMarblesByRangeWithPagination" {
		// 	return t.getMarblesByRangeWithPagination(stub, args)
		// } else if function == "queryMarblesWithPagination" {
		// 	return t.queryMarblesWithPagination(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ===============================================================
// initProduct - create a new product, store into chaincode state
// ===============================================================
func (t *ProductChaincode) initProduct(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error

	// TODO: create incremental ID
	id := uint64(1)
	createdBy := "user123"
	createdAt := "2019-03-17 22:45:35 UTC"
	updatedBy := ""
	updatedAt := ""

	//   0                  1              2                  3         4
	// GTIN,             Producer,      ContainedProducts, Labels,    Locale
	// "7612100055557", "Wander AG",    "[]",              `["UTZ"]`, `["lang": "de", ...]`
	// or ""
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init product")

	// === Arg 0: GTIN ===
	gtin := args[0]
	if len(gtin) > 0 {
		fmt.Println("GTIN: " + gtin)
	} else {
		fmt.Println("GTIN not provided")
	}

	// === Arg 1: Producer ===
	producer := args[1]
	if len(producer) > 0 {
		fmt.Println("Producer: " + gtin)
	} else {
		fmt.Println("Producer not provided")
	}

	// === Arg 2: ContainedProducts ===
	var containedProducts []string
	err = json.Unmarshal([]byte(args[2]), &containedProducts)
	if err != nil {
		return shim.Error("3rd argument 'containedProducts' must be a string with " +
			"a JSON list of IDs of contained products: [\"123\", \"456\", ...] " +
			"(or an empty list: [])")
	}
	if len(containedProducts) > 0 {
		fmt.Printf("ContainedProducts: %v", containedProducts)
	} else {
		fmt.Println("ContainedProducts not provided")
	}

	// === Arg 3: Labels ===
	var labels []string
	err = json.Unmarshal([]byte(args[3]), &labels)
	if err != nil {
		return shim.Error("4th argument 'labels' must be a string with " +
			"a JSON list of labels labelling this product: [\"Fairtrade\", \"GOTS\", ...]" +
			"(or an empty list: [])")
	}
	if len(labels) > 0 {
		fmt.Printf("Labels: %v", labels)
	} else {
		fmt.Println("Labels not provided")
	}

	// === Arg 4: Locale ===
	var locale []ProductLocaleData
	err = json.Unmarshal([]byte(args[4]), &locale)
	if err != nil {
		return shim.Error("5th argument 'locale' must be a string with " +
			"a JSON list of objects with keys 'lang', 'name', 'price', 'currency', " +
			"'description', 'quantity', 'ingredients', 'packaging', 'categories', " +
			"'image', 'productUrl', where each contains a string, except 'packaging' " +
			"and 'categories' contain a list of strings.")
	}
	if len(locale) > 0 {
		fmt.Printf("Locale: %+v", locale)
	} else {
		fmt.Println("Locale not provided")
	}

	// ==== Check if product with this GTIN already exists ====
	// TODO: use stub.GetQueryResult(query string)
	// productAsBytes, err := stub.GetState(productName)
	// if err != nil {
	// 	return shim.Error("Failed to get product: " + err.Error())
	// } else if productAsBytes != nil {
	// 	fmt.Println("This product already exists: " + productName)
	// 	return shim.Error("This product already exists: " + productName)
	// }

	// Create new initial score
	score := Score{Environment: 0, Climate: 0, Society: 0, Health: 0, Economy: 0}

	// ==== Create product object and marshal to JSON ====
	objectType := "product"
	product := &Product{objectType, id, gtin, createdBy, createdAt, updatedBy, updatedAt,
		producer, containedProducts, labels, locale, score, "active"}
	productJSONasBytes, err := json.Marshal(product)
	if err != nil {
		return shim.Error(err.Error())
	}
	//Alternatively, build the product json string manually if you don't want to use struct marshalling
	//productJSONasString := `{"docType":"product",  "name": "` + productName + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + owner + `"}`
	//productJSONasBytes := []byte(str)

	// === Save product to state ===
	err = stub.PutState(string(id), productJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// //  ==== Index the product to enable color-based range queries, e.g. return all blue products ====
	// //  An 'index' is a normal key/value entry in state.
	// //  The key is a composite key, with the elements that you want to range query on listed first.
	// //  In our case, the composite key is based on indexName~color~name.
	// //  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	// indexName := "color~name"
	// colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{product.Color, product.Name})
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the product.
	// //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	// value := []byte{0x00}
	// stub.PutState(colorNameIndexKey, value)

	// ==== Product saved and indexed. Return success ====
	fmt.Println("- end init product")
	return shim.Success(nil)
}

// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// =======Rich queries =========================================================================
// Two examples of rich queries are provided below (parameterized query and ad hoc query).
// Rich queries pass a query string to the state database.
// Rich queries are only supported by state database implementations
//  that support rich query (e.g. CouchDB).
// The query string is in the syntax of the underlying state database.
// With rich queries there is no guarantee that the result set hasn't changed between
//  endorsement time and commit time, aka 'phantom reads'.
// Therefore, rich queries should not be used in update transactions, unless the
// application handles the possibility of result set changes between endorsement and commit time.
// Rich queries can be used for point-in-time queries against a peer.
// ============================================================================================

// ===== Example: Parameterized rich query =================================================
// queryProductsByGTIN queries for products based on a passed in GTIN number (barcode).
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (GTIN).
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *ProductChaincode) queryProductsByGTIN(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	//       0
	// "7612100055557"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	gtin := args[0]

	queryString := fmt.Sprintf("{\"selector\": {\"docType\": \"product\", \"gtin\": \"%s\"}}", gtin)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
