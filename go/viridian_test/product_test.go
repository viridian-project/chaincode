package viridian_test

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/chaincode/viridian/go/viridian"
)

var _ = Describe("Product", func() {
	stub := shim.NewMockStub("testingStub", new(viridian.Chaincode))
	status200 := int32(200)

	BeforeSuite(func() {
		stub.MockInit("000", nil)
	})

	Describe("Checking product lifecycle", func() {
		It("Should be possible to add a new product", func() {
			args := [][]byte{
				[]byte("addProduct"),
				[]byte("1fcc2c43-12a1-4451-ac56-dd73099b3f34"),          // key
				[]byte("7612100055557"),                                 // GTIN
				[]byte("producer-84a234b7-c9d8-43b2-93c9-90f83d8773fb"), // producer key
				[]byte("[]"), // contained product keys
				[]byte("[\"label-31d3a05e-fb10-483c-8c8b-0c7079e5bc95\"]"), // label keys
				[]byte("[{\"lang\": \"de\", \"name\": \"Ovomaltine crunchy cream - 400 g\",\"price\": \"4.99\",\"currency\": \"EUR\",\"description\": \"Brotaufstrich mit malzhaltigem Getraenkepulver Ovomaltine\",\"quantities\": [\"400 g\"]}]"), // locales
			}
			response := stub.MockInvoke("000", args)
			fmt.Println(response)
			Expect(response.Status).Should(Equal(status200))
		})
	})
})
