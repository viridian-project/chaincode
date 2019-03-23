# Viridian backend chaincode

The Viridian Project plans to employ a [Hyperledger Fabric](https://hyperledger-fabric.readthedocs.io)
blockchain network to manage its assets: data and information on the
sustainability of products, producers and product/producer labels.

The code in this repository is the so-called 'chaincode': software installed on
the peers participating in the blockchain network that is used to both access
the data from the distributed ledger as well as make modifications to it.

The main purposes of the chaincode are to define the data structure, provide an
API to access data from the distributed ledger and make sure that modifications
to the ledger are legitimate. When a request to change the ledger is sent via
the network, a certain number of peers need to validate the request and agree
that the proposed change is legitimate by adding their signature.

This chaincode follows along the two tutorials
[Chaincode for Developers](https://hyperledger-fabric.readthedocs.io/en/release-1.4/chaincode4ade.html) and
[Using CouchDB](https://hyperledger-fabric.readthedocs.io/en/release-1.4/couchdb_tutorial.html).
The code is mainly based on the marbles02 example provided with a Hyperledger
Fabric installation in `fabric-samples/chaincode/marbles02/go/marbles_chaincode.go`.

## Deployment

### Prerequisites

First install Hyperledger Fabric and Fabric's prerequisites, Docker and Go.

Download the Hyperledger Fabric samples.

### Download the code

Place the content of this repository under Fabric's sample directory in `fabric-samples/chaincode`,
e.g. in a directory `fabric-samples/chaincode/viridian`.

```
cd fabric-samples/chaincode
git clone https://github.com/viridian-project/chaincode.git viridian
```

### Compile the chaincode

See https://stackoverflow.com/questions/37433618/how-to-use-a-chaincode-thats-not-on-github?rq=1.

First install go package dependencies:

```
go get -u github.com/hyperledger/fabric/core/chaincode/shim
# go get -u github.com/hyperledger/fabric/protos/peer # seems unnecessary after first command
```

Then you can build:

```
cd fabric-samples/chaincode/viridian/go
# go build
go test -run BuildImage_Peer
```

Is this step really necessary? But it was useful for debugging.


### Start a test network

```
cd fabric-samples/first-network

# Make sure previous networks are removed so that we have a clean statedb
./byfn.sh down

# Start up BYFN network with COUCHDB
./byfn.sh up -c mychannel -s couchdb

# Login to docker container named 'cli'
docker exec -it cli bash

# Install chaincode:
peer chaincode install -n viridian -v 1.0 -p github.com/chaincode/viridian/go

# Instantiate chaincode:
export CHANNEL_NAME=mychannel
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n viridian -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org0MSP.peer','Org1MSP.peer')"
```

Verify that it worked by looking in logs if CouchDB index was created:

```
docker logs peer0.org1.example.com  2>&1 | grep "CouchDB index"
```

It should return

```
[couchdb] CreateIndex -> INFO 089 Created CouchDB index [indexProductGTIN] in state database [mychannel_viridian] using design document [_design/indexProductGTINDoc]
```

### Insert first test product

```
docker exec -it cli bash
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n viridian -c '{"Args":["initProduct","7612100055557","Wander AG","[]","[\"UTZ\"]", "[{\"lang\": \"de\", \"name\": \"Ovomaltine crunchy cream - 400 g\",\"price\": \"4.99\",\"currency\": \"EUR\",\"description\": \"Brotaufstrich mit malzhaltigem Getraenkepulver Ovomaltine\",\"quantity\": \"400 g\"}]"]}'
```

### Query for product by gtin

```
docker exec -it cli bash
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n viridian -c '{"Args":["queryProductsByGTIN","7612100055557"]}'
```
