#!/bin/bash
set -ex
LANGUAGE=${1:-"golang"}
CHAINCODE_NAME=maximo
echo "Copying Chaincode to cli container"
docker cp ../chaincode cli:/opt/gopath/src/github.com/maximo
echo "Install and Instantiate Chaincode"
docker exec cli mkdir -p /opt/gopath/src/github.com/maximo
docker exec cli go get github.com/hyperledger/fabric-chaincode-go/shim
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode install -n maximo -v 1.0 -p github.com/maximo
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n maximo -v 1.0 -c '{"Args":["101"]}'
echo "Chaincode Instantiated"
sleep 10
echo "Test Chaincode"
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n maximo -c '{"Args":["read_everything"]}'
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n maximo -c '{"Args":["init_asset", "asset1"]}'
