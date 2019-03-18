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
[https://hyperledger-fabric.readthedocs.io/en/release-1.4/chaincode4ade.html](Chaincode
for Developers) and
[https://hyperledger-fabric.readthedocs.io/en/release-1.4/couchdb_tutorial.html](Using CouchDB).
The code is mainly based on the marbles02 example provided with a Hyperledger
Fabric installation in `fabric-samples/chaincode/marbles02/go/marbles_chaincode.go`.

## Deployment

The content of this directory should be placed under `$GOPATH/src`, e.g. in a
directory `$GOPATH/src/viridian`.

<!--
### Compile the chaincode

```
go get -u github.com/hyperledger/fabric/core/chaincode/shim
go build
```
-->

### Start a test newtwork

```
cd fabric-samples/first-network
# Make sure previous networks are removed so that we have a clean statedb
./byfn.sh down
# Start up BYFN network with COUCHDB
./byfn.sh up -c mychannel -s couchdb
```

### Install and instantiate chaincode

```
docker exec -it cli bash
```
