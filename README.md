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

## The model

A heritage of the now inactive project Hyperledger Composer is the modeling language [Concerto](https://github.com/hyperledger/composer-concerto).

We use the Concerto language to model our application data in a file `model/org.viridian.cto`.

With the [Concerto Tools](https://github.com/hyperledger/composer-concerto-tools), we can convert this file to UML, Go or other languages.

```
npm install -g composer-concerto-tools # need to install Node.js which also ships the npm package manager
node ~/.nvm/versions/node/v8.16.0/lib/node_modules/composer-concerto-tools/cli.js generate --ctoFiles model/org.viridian.cto --format PlantUML
node ~/.nvm/versions/node/v8.16.0/lib/node_modules/composer-concerto-tools/cli.js generate --ctoFiles model/org.viridian.cto --format Go
```

The conversions are written to files in the `output` folder.

Convert the UML file to a PNG image with a UML diagram:

```
sudo apt install plantuml
cd output
plantuml
```

## Deployment

### Prerequisites

First install Hyperledger Fabric's prerequisites, Docker and Go.

#### For Ubuntu Linux (tested on 16.04)

##### Install Docker

```
sudo apt-get update

sudo apt-get install apt-transport-https ca-certificates curl gnupg-agent software-properties-common

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

sudo apt-key fingerprint 0EBFCD88

sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli containerd.io

# See if it works:

sudo docker run hello-world

# For default user to be able to interact with Docker daemon, add it to the group `docker`: (see https://docs.docker.com/install/linux/linux-postinstall/)

sudo usermod -aG docker $USER

# Log out of the OS and in again
# If the following command returns no error message, it worked:

docker ps

# Now docker run can also run without sudo:

docker run hello-world

# Install docker-compose: (needed later to start up the blockchain)

sudo pip3 install docker-compose
```

##### Install Go

```
# Download:
curl -sL -o go1.11.5.linux-amd64.tar.gz https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz

# Compare checksum with the one on https://golang.org/dl/:
sha256sum go1.11.5.linux-amd64.tar.gz

# Untar to install:
sudo tar -C /usr/local -xzf go1.11.5.linux-amd64.tar.gz

# Set environment variable to make go executables executable
# Put this lines at the end of the file `~/.bashrc`:
export PATH=$PATH:/usr/local/go/bin

# Set a gopath variable where you will later put your
# chaincode.
# Also downloaded go packages are put there (see below).
# Add lines in file `~/.bashrc`:
export GOPATH=$HOME/path/to/my/chaincode_dir
export PATH=$PATH:$GOPATH/bin

# In terminal type commands:
mkdir $GOPATH/bin
mkdir $GOPATH/src

# Optional: install some extra go packages:
# A REPL:
go get -u github.com/motemen/gore/cmd/gore
# For code-completion in the REPL:
go get -u github.com/mdempsky/gocode
# Sth. else for REPL:
go get -u github.com/k0kubun/pp
# For using :doc in REPL:
go get -u golang.org/x/tools/cmd/godoc
```

#### Install Hyperledger Fabric

Download the Hyperledger Fabric samples.

```
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/master/scripts/bootstrap.sh | bash -s 1.4.0
# Shorter version:
# curl -sSL http://bit.ly/2ysbOFE | bash -s 1.4.0

cd fabric-samples

echo $(pwd)/bin

# Take the result of the last command ($(pwd)/bin) and put
# it into a new line in `~/.bashrc`:
# New line in `~/.bashrc`:
export PATH=<output of echo $(pwd)/bin>:$PATH

# If you haven't done so before, set the $GOPATH variable to a location
# where you will later store your chaincode:
export GOPATH=$HOME/path/to/my/chaincode_dir
export PATH=$PATH:$GOPATH/bin

mkdir $GOPATH/bin
mkdir $GOPATH/src
```

### Download the code

Place the content of this repository under Fabric's sample directory in `fabric-samples/chaincode`,
e.g. in a directory `fabric-samples/chaincode/viridian`.

```
cd fabric-samples/chaincode
git clone https://github.com/viridian-project/chaincode.git viridian
```

### Add external dependencies ("vendoring")

To add the external dependency "client identity" (CID) library, which is needed to access the client's, i.e. user's, certificate, use `govendor`:

Install govendor:

```
go get -u github.com/kardianos/govendor
```

In your chaincode directory (i.e. inside the `go` directory):

(Note: For this to work, the `go` directory must be somewhere under the `$GOPATH`. You can create a symbolic link `viridian` inside `$GOPATH/src/github.com/chaincode` that points to the `fabric-samples/chaincode/viridian` directory and then `cd` into `$GOPATH/src/github.com/chaincode/viridian/go`.) With `$GOPATH/src/github.com/chaincode`, you obtain the same directory structure as on a Fabric node.

```
govendor init
govendor add github.com/hyperledger/fabric/core/chaincode/shim/ext/cid
govendor add github.com/golang/protobuf/proto # Dep of cid
govendor add github.com/pkg/errors # Dep of cid
```

This creates a `vendor` directory that is accessible to the chaincode when the code is installed on the peers.

### Compile the chaincode for testing

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

### Test the chaincode in development mode

Following `fabric-samples/chaincode-docker-devmode/README.rst`. This method is a bit faster than the one described under "Deploy the chaincode for testing using CouchDB". Unfortunately, this development mode uses leveldb and not CouchDB, so the Viridian chaincode is not fully compatible.

Open three terminal windows with `fabric-samples/chaincode-docker-devmode/` as working directory.

In terminal 1:
```
docker-compose -f docker-compose-simple.yaml up
```

In terminal 2:
```
docker exec -it chaincode bash
cd viridian/go
go build -o viridian
# If it has worked, an executable (green color upon `ls`) file `viridian` was created.
CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=viridian:0 ./viridian
```

In terminal 3:
```
docker exec -it cli bash
peer chaincode install -p chaincodedev/chaincode/viridian/go -n viridian -v 0
peer chaincode instantiate -C myc -n viridian -v 0 -c '{"Args":["init"]}'

# Insert first test product:
peer chaincode invoke -C myc -n viridian -c '{"Args":["addProduct","1fcc2c43-12a1-4451-ac56-dd73099b3f34","7612100055557","producer-84a234b7-c9d8-43b2-93c9-90f83d8773fb","[]","[\"label-31d3a05e-fb10-483c-8c8b-0c7079e5bc95\"]", "[{\"lang\": \"de\", \"name\": \"Ovomaltine crunchy cream - 400 g\",\"price\": \"4.99\",\"currency\": \"EUR\",\"description\": \"Brotaufstrich mit malzhaltigem Getraenkepulver Ovomaltine\",\"quantities\": [\"400 g\"]}]"]}'

# Insert the first test producer:
peer chaincode invoke -C myc -n viridian -c '{"Args":["initProducer","84a234b7-c9d8-43b2-93c9-90f83d8773fb","Wander AG","CH-3176 Neuenegg, Switzerland","https://www.wander.ch/","[]"]}'
```

#### Shut down and start again

Hit Ctrl-C in terminal 2 to stop the chaincode. Hit Ctrl-D in terminals 2 and 3 to log out. Hit Ctrl-C in terminal 1 to stop the containers.

Remove the containers to be able to start fresh:

```
# Look for the container IDs with
docker ps -a
# Remove the containers `hyperledger/fabric-ccenv`, `hyperledger/fabric-tools`, `hyperledger/fabric-peer` and `hyperledger/fabric-orderer` by entering their IDs, e.g.:
docker rm db8f289923ea 9edbb221665c 7070bf7e8409 32de20667557
```

Now you have a clean state and could start fresh with `docker-compose -f docker-compose-simple.yaml up`.

### Deploy the chaincode for testing using CouchDB (not for production)

```
cd fabric-samples/first-network

# Make sure previous networks are removed so that we have a clean statedb
./byfn.sh down
# Remove any remnant docker containers: Look for images like `dev-peer0.org1.example.com-viridian-1.0-c1c88edc790...`:
docker ps -a
# If so, remove that container, e.g. using its names:
docker rm dev-peer0.org1.example.com-viridian-1.0
# Look if there is a remnant chaincode image like `dev-peer0.org1.example.com-viridian-1.0-c1c88edc790...`:
docker images
# If so, remove that image, e.g. using its ID:
docker rmi 0d6a11ebeee0

# Start up BYFN network with COUCHDB
./byfn.sh up -c mychannel -s couchdb

# Login to docker container named 'cli'
docker exec -it cli bash

# Install chaincode:
peer chaincode install -n viridian -v 1.0 -p github.com/chaincode/viridian/go/

# Instantiate chaincode:
export CAFILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $CAFILE -C mychannel -n viridian -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"
```

#### Optional:

Verify that it worked by looking in logs if CouchDB index was created:

```
docker logs peer0.org1.example.com  2>&1 | grep "CouchDB index"
```

It should return

```
[couchdb] CreateIndex -> INFO 089 Created CouchDB index [indexProductGTIN] in state database [mychannel_viridian] using design document [_design/indexProductGTINDoc]
```

#### Insert first test product

Inside the `cli` docker container:

```
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $CAFILE -C mychannel -n viridian -c '{"Args":["addProduct","1fcc2c43-12a1-4451-ac56-dd73099b3f34","7612100055557","producer-84a234b7-c9d8-43b2-93c9-90f83d8773fb","[]","[\"label-31d3a05e-fb10-483c-8c8b-0c7079e5bc95\"]", "[{\"lang\": \"de\", \"name\": \"Ovomaltine crunchy cream - 400 g\",\"price\": \"4.99\",\"currency\": \"EUR\",\"description\": \"Brotaufstrich mit malzhaltigem Getraenkepulver Ovomaltine\",\"quantities\": [\"400 g\"]}]"]}'
```

#### Query for product by GTIN

Inside the `cli` docker container:

```
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $CAFILE -C mychannel -n viridian -c '{"Args":["queryProductsByGTIN","7612100055557"]}'
```

#### Insert the first test producer

```
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile $CAFILE -C mychannel -n viridian -c '{"Args":["initProducer","84a234b7-c9d8-43b2-93c9-90f83d8773fb","Wander AG","CH-3176 Neuenegg, Switzerland","https://www.wander.ch/","[]"]}'
```

#### Install new version of chaincode

```
peer chaincode install -n viridian -v 1.1 -p github.com/chaincode/viridian/go/
peer chaincode upgrade -o orderer.example.com:7050 --tls --cafile $CAFILE -C mychannel -n viridian -v 1.1 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"
```

#### Remove old version of chaincode

See https://stackoverflow.com/questions/51015655/hyperledger-fabric-how-to-remove-a-chaincode-on-peer

```
# On host system:
# Stop and remove the docker container (and image) that belongs to the chaincode
docker ps     # Look for the ID/name if necessary
docker images # Look for the ID/name if necessary
docker stop dev-peer0.org1.example.com-viridian-1.0
docker rm dev-peer0.org1.example.com-viridian-1.0
docker rmi dev-peer0.org1.example.com-viridian-1.0-c1c88edc79...
# Enter the peer on which chaincode was installed:
docker exec -it peer0.org1.example.com bash
# Remove the chaincode file:
> rm /var/hyperledger/production/chaincodes/viridian.1.0
```

### Write a unit test

From: https://blogs.sap.com/2019/01/11/how-to-write-unit-tests-for-hyperledger-fabric-go-chaincode/
See also: https://medium.com/coinmonks/test-driven-hyperledger-fabric-golang-chaincode-development-dbec4cb78049

```
go get -u github.com/onsi/ginkgo/ginkgo
go get -u github.com/onsi/gomega/...
# maybe not needed, maybe yes: go get -u github.com/s7techlab/cckit
cd viridian/go
mkdir viridian_test
ginkgo bootstrap # generates the test suite file (which runs a suite of tests)
ginkgo generate product # generates a test file `product_test.go` alongside `product.go`
```
