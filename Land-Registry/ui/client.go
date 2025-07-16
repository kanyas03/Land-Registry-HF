package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func submitTxnFn(
	organization string,
	channelName string,
	chaincodeName string,
	contractName string,
	txnType string,
	privateData map[string][]byte,
	txnName string,
	args ...string,
) string {

	orgProfile := profile[organization]

	clientConn := newGrpcConnection(orgProfile.TLSCertPath, orgProfile.GatewayPeer, orgProfile.PeerEndpoint)
	defer clientConn.Close()

	id := newIdentity(orgProfile.CertPath, orgProfile.MSPID)
	sign := newSign(orgProfile.KeyDirectory)

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	fmt.Printf("\n--> Submitting Transaction: %s\n", txnName)

	switch txnType {

	case "invoke":
		result, err := contract.SubmitTransaction(txnName, args...)
		if err != nil {
			panic(fmt.Errorf("submit failed: %w", err))
		}
		return fmt.Sprintf("*** Transaction Success: %s\n", result)

	case "query":
		result, err := contract.EvaluateTransaction(txnName, args...)
		if err != nil {
			panic(fmt.Errorf("query failed: %w", err))
		}
		if isByteSliceEmpty(result) {
			return string(result)
		}
		return formatJSON(result)

	case "private":
		result, err := contract.Submit(
			txnName,
			client.WithArguments(args...),
			client.WithTransient(privateData),
		)
		if err != nil {
			panic(fmt.Errorf("private submit failed: %w", err))
		}
		return fmt.Sprintf("*** Private Transaction Committed:\n%s\n", result)
	}

	return "Invalid transaction type"
}
