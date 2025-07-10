package main

import (
	contracts "landcontract/Contracts"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	LandContract := new(contracts.LandContract)

	chaincode, err := contractapi.NewChaincode(LandContract)

	if err != nil {
		log.Panicf("Could not create chaincode : %v", err)
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode : %v", err)
	}
}
