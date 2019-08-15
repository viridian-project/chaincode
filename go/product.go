package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/hyperledger/fabric/protos/peer"
)

// ProductChaincode is the chaincode associated with products
type ProductChaincode struct {
}

// ProductLocaleData is the locale-specific (language-specific) part of a product
type ProductLocaleData struct {
	Lang        string   `json:"lang"`        // regex=/^[a-z]{2}$/ // ISO language code according to https://en.wikipedia.org/wiki/ISO_639-1, there should be only one locale data for each language
	Name        string   `json:"name"`        // product 'short name'
	Price       string   `json:"price"`       // optional
	Currency    string   `json:"currency"`    // optional
	Description string   `json:"description"` // optional
	Quantities  []string `json:"quantities"`
	Ingredients string   `json:"ingredients"` // optional
	Packagings  []string `json:"packagings"`
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
//   Description: "Brotaufstrich mit malzhaltigem Getränkepulver Ovomaltine",
//   Quantities: []string{"400 g"},
//   Ingredients: "33% malzhaltiges Getränkepulver: Ovomaltine (Gerstenmalzextrakt, kondensierte Magermilch, kondensiertes Milchserum, fettarmer Kakao, Zucker, Fruktose, Magnesiumcarbonat, Calciumphosphat, Rapsöl, Vitamine [A, E, B1, B2, Pantothensäure, B6, Folsäure, B12, C, Biotin, Niacin], Kochsalz, Aroma Vanillin), Zucker, Pflanzenöle (Raps- und Palmöl), 2.6% Haselnüsse, Calciumphosphat, fettarmer Kakao, Emulgator Sonnenblumenlecithin, Aroma Vanillin.",
//   Packagings: []string{"Glas", "Plastik"},
//   Categories: []string{"Brotaufstriche", "Frühstück", "Nougatcremes"},
//   ImageURL: "ipfs://jf3f03-kf30-fk3-kf3-fk3.png",
//   URL: "http://www.ovomaltine.de/produkte/ovomaltine-crunchy-cream-1/",
// }

// Product is the asset representing a product
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
//   ScorableAsset{
//	   UpdatableAsset{
//		   ReviewableAsset{
//         CreatedBy: "user123",
//         CreatedAt: "2018-12-24 12:11:54 UTC",
//         Status: Active,
//       },
//       UpdatedBy: "user123",
//       UpdatedAt: "2019-02-10 18:33:39 UTC",
//       Supersedes: "product-bc36d43e-c40c-40c0-8086-dc19bc000fe1",
//       SupersededBy: "",
//       ChangeReason: "Wrong quantity information.",
//		 },
//		 Score: &Score{
//       Environment: -34,
//       Climate: -46,
//       Society: -7,
//       Health: -78,
//       Economy: 21,
//     },
//   },
//   DocType: "product",
//   GTIN: "7612100055557",
//   Producer: "producer-afd05a40-4ed6-4ae5-8120-eb7daebc336c",
//   ContainedProducts: []string{},
//   Labels: []string{"label-42c2f586-a893-485f-8995-8639446bb6b8"},
//   Locale: []ProductLocaleData{
//     &ProductLocaleData{
//       Lang: "de",
//       Name: "Ovomaltine crunchy cream — 400 g",
//       Price: "4.99",
//       Currency: "€",
//       Description: "Brotaufstrich mit malzhaltigem Getränkepulver Ovomaltine",
//       Quantities: []string{"400 g"},
//       Ingredients: "33% malzhaltiges Getränkepulver: Ovomaltine (Gerstenmalzextrakt, kondensierte Magermilch, kondensiertes Milchserum, fettarmer Kakao, Zucker, Fruktose, Magnesiumcarbonat, Calciumphosphat, Rapsöl, Vitamine [A, E, B1, B2, Pantothensäure, B6, Folsäure, B12, C, Biotin, Niacin], Kochsalz, Aroma Vanillin), Zucker, Pflanzenöle (Raps- und Palmöl), 2.6% Haselnüsse, Calciumphosphat, fettarmer Kakao, Emulgator Sonnenblumenlecithin, Aroma Vanillin.",
//       Packagings: []string{"Glas", "Plastik"},
//       Categories: []string{"Brotaufstriche", "Frühstück", "Nougatcremes"},
//       ImageURL: "ipfs://jf3f03-kf30-fk3-kf3-fk3.png",
//       URL: "http://www.ovomaltine.de/produkte/ovomaltine-crunchy-cream-1/",
//     },
//   },
// }

// ===============================================================
// initProduct - create a new product, store into chaincode state
// ===============================================================
func (c *ProductChaincode) initProduct(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// Arguments:
	//  0                     1                  2                                  3            4                                   5
	// Key,                 GTIN,            Producer,                     ContainedProducts, Labels,                             Locales
	// "8a259c61-6825-...", "7612100055557", "producer-a3006838-bdf2-...", "[]",              `["label-31d3a05e-fb10-...", ...]`, `[{"lang": "de", ...}]`
	// or ""
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6.")
	}

	var err error

	// Create initial values
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

	// ==== Input sanitation ====
	fmt.Println("- start init product")

	// === Arg 0: Key ===
	key := args[0]
	if len(key) > 0 {
		fmt.Println("Key: " + key)
	} else {
		fmt.Println("Key not provided")
	}

	// === Arg 1: GTIN ===
	gtin := args[1]
	if len(gtin) > 0 {
		fmt.Println("GTIN: " + gtin)
	} else {
		fmt.Println("GTIN not provided")
	}

	// === Arg 2: Producer ===
	producer := args[2]
	if len(producer) > 0 {
		fmt.Println("Producer: " + gtin)
	} else {
		fmt.Println("Producer not provided")
	}

	// === Arg 3: ContainedProducts ===
	var containedProducts []string
	err = json.Unmarshal([]byte(args[3]), &containedProducts)
	if err != nil {
		return shim.Error("4th argument 'containedProducts' must be a string with " +
			"a JSON list of Keys of contained products, e.g.: [\"product-123\", \"product-456\", ...] " +
			"(or an empty list: [])")
	}
	if len(containedProducts) > 0 {
		fmt.Printf("ContainedProducts: %v", containedProducts)
	} else {
		fmt.Println("ContainedProducts not provided")
	}

	// === Arg 4: Labels ===
	var labels []string
	err = json.Unmarshal([]byte(args[4]), &labels)
	if err != nil {
		return shim.Error("5th argument 'labels' must be a string with " +
			"a JSON list of label Keys labelling this product: [\"label-bd80e824-938c-...\", \"label-127cc795-3a20-...\", ...]" +
			"(or an empty list: [])")
	}
	if len(labels) > 0 {
		fmt.Printf("Labels: %v", labels)
	} else {
		fmt.Println("Labels not provided")
	}

	// === Arg 5: Locale ===
	var locale []ProductLocaleData
	err = json.Unmarshal([]byte(args[5]), &locale)
	if err != nil {
		return shim.Error("6th argument 'locale' must be a string with " +
			"a JSON list of objects with keys 'lang', 'name', 'price', 'currency', " +
			"'description', 'quantities', 'ingredients', 'packagings', 'categories', " +
			"'imageUrl', 'url', where each contains a string, except 'quantities', " +
			"'packagings' and 'categories' contain a list of strings.")
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

	// ==== Create product object and marshal to JSON ====
	docType := "product"
	product := &Product{
		ScorableAsset{
			UpdatableAsset{
				ReviewableAsset{createdBy, createdAt, Preliminary},
				updatedBy, updatedAt, supersedes, supersededBy, changeReason},
			score},
		docType, gtin, producer, containedProducts, labels, locale}
	jsonAsBytes, err := json.Marshal(product)
	if err != nil {
		return shim.Error(err.Error())
	}
	//Alternatively, build the product json string manually if you don'`t` want to use struct marshalling
	//productJSONasString := `{"docType":"product",  "name": "` + productName + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + owner + `"}`
	//jsonAsBytes := []byte(str)

	// === Save product to state ===
	err = stub.PutState(docType+"-"+key, jsonAsBytes)
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
func (c *ProductChaincode) queryProductsByGTIN(stub shim.ChaincodeStubInterface, args []string) peer.Response {
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
