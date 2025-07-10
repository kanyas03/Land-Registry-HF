package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract structure
type LandContract struct {
	contractapi.Contract
}

// Public land data
type Land struct {
	Location string   `json:"location"`
	Size     string   `json:"size"`
	Status   string   `json:"status"`
	History  []string `json:"history"`
}

// Private land ownership (Org1 & Org2)
type PrivateDetails struct {
	Owner string `json:"owner"`
}

// Private purchase requests (Org1 & Org3)
type PurchaseRequest struct {
	BuyerID string `json:"buyerId"`
}

// Org1 - Register land
func (c *LandContract) RegisterLand(ctx contractapi.TransactionContextInterface, landID, location, size, owner string) error {
	msp, _ := ctx.GetClientIdentity().GetMSPID()
	if msp != "Org1MSP" {
		return fmt.Errorf("only Org1 (Government) can register land")
	}

	exists, err := ctx.GetStub().GetState(landID)
	if err != nil {
		return fmt.Errorf("failed to check existing land: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("land %s already exists", landID)
	}

	land := Land{
		Location: location,
		Size:     size,
		Status:   "Registered",
		History:  []string{"Registered by Government"},
	}
	landBytes, _ := json.Marshal(land)
	err = ctx.GetStub().PutState(landID, landBytes)
	if err != nil {
		return fmt.Errorf("failed to store land: %v", err)
	}

	private := PrivateDetails{Owner: owner}
	privateBytes, _ := json.Marshal(private)
	return ctx.GetStub().PutPrivateData("collectionLandPrivateDetails", landID, privateBytes)
}

// Org2 - Mark land for sale
func (c *LandContract) MarkForSale(ctx contractapi.TransactionContextInterface, landID, sellerID string) error {
	msp, _ := ctx.GetClientIdentity().GetMSPID()
	if msp != "Org2MSP" {
		return fmt.Errorf("only Org2 (Seller) can mark land for sale")
	}

	landBytes, err := ctx.GetStub().GetState(landID)
	if err != nil || landBytes == nil {
		return fmt.Errorf("land not found")
	}

	privBytes, err := ctx.GetStub().GetPrivateData("collectionLandPrivateDetails", landID)
	if err != nil || privBytes == nil {
		return fmt.Errorf("private owner data not found")
	}

	var priv PrivateDetails
	_ = json.Unmarshal(privBytes, &priv)
	if priv.Owner != sellerID {
		return fmt.Errorf("only the owner can mark land for sale")
	}

	var land Land
	_ = json.Unmarshal(landBytes, &land)

	land.Status = "For Sale"
	land.History = append(land.History, "Marked for Sale by "+sellerID)

	updated, _ := json.Marshal(land)
	return ctx.GetStub().PutState(landID, updated)
}

// Org3 - Request to buy land
func (c *LandContract) RequestPurchase(ctx contractapi.TransactionContextInterface, landID, buyerID string) error {
	msp, _ := ctx.GetClientIdentity().GetMSPID()
	if msp != "Org3MSP" {
		return fmt.Errorf("only Org3 (Buyer) can request a purchase")
	}

	landBytes, err := ctx.GetStub().GetState(landID)
	if err != nil || landBytes == nil {
		return fmt.Errorf("land not found")
	}

	var land Land
	_ = json.Unmarshal(landBytes, &land)
	if land.Status != "For Sale" {
		return fmt.Errorf("land is not available for sale")
	}

	request := PurchaseRequest{BuyerID: buyerID}
	reqBytes, _ := json.Marshal(request)
	err = ctx.GetStub().PutPrivateData("collectionPurchaseRequests", landID, reqBytes)
	if err != nil {
		return fmt.Errorf("failed to store purchase request: %v", err)
	}

	land.Status = "Requested"
	land.History = append(land.History, "Purchase Requested by "+buyerID)
	updated, _ := json.Marshal(land)
	return ctx.GetStub().PutState(landID, updated)
}

// Org1 - Approve and transfer land
func (c *LandContract) ApprovePurchase(ctx contractapi.TransactionContextInterface, landID string) error {
	msp, _ := ctx.GetClientIdentity().GetMSPID()
	if msp != "Org1MSP" {
		return fmt.Errorf("only Org1 (Government) can approve the request")
	}

	landBytes, err := ctx.GetStub().GetState(landID)
	if err != nil || landBytes == nil {
		return fmt.Errorf("land not found")
	}

	reqBytes, err := ctx.GetStub().GetPrivateData("collectionPurchaseRequests", landID)
	if err != nil || reqBytes == nil {
		return fmt.Errorf("no purchase request found")
	}

	var request PurchaseRequest
	_ = json.Unmarshal(reqBytes, &request)

	var land Land
	_ = json.Unmarshal(landBytes, &land)
	land.Status = "Transferred"
	land.History = append(land.History, "Ownership transferred to "+request.BuyerID)

	updatedLand, _ := json.Marshal(land)
	err = ctx.GetStub().PutState(landID, updatedLand)
	if err != nil {
		return fmt.Errorf("failed to update land: %v", err)
	}

	newPriv := PrivateDetails{Owner: request.BuyerID}
	privBytes, _ := json.Marshal(newPriv)
	err = ctx.GetStub().PutPrivateData("collectionLandPrivateDetails", landID, privBytes)
	if err != nil {
		return fmt.Errorf("failed to update ownership: %v", err)
	}

	return nil
}

// Org1 or Org2 - View owner info
func (c *LandContract) GetOwner(ctx contractapi.TransactionContextInterface, landID string) (string, error) {
	msp, _ := ctx.GetClientIdentity().GetMSPID()
	if msp != "Org1MSP" && msp != "Org2MSP" {
		return "", fmt.Errorf("only Org1 or Org2 can view owner info")
	}

	privBytes, err := ctx.GetStub().GetPrivateData("collectionLandPrivateDetails", landID)
	if err != nil || privBytes == nil {
		return "", fmt.Errorf("owner info not found")
	}

	var priv PrivateDetails
	_ = json.Unmarshal(privBytes, &priv)
	return priv.Owner, nil
}

// Anyone - View land
func (c *LandContract) GetLand(ctx contractapi.TransactionContextInterface, landID string) (*Land, error) {
	landBytes, err := ctx.GetStub().GetState(landID)
	if err != nil || landBytes == nil {
		return nil, fmt.Errorf("land not found")
	}

	var land Land
	_ = json.Unmarshal(landBytes, &land)
	return &land, nil
}

// Org3 - View available land (status = For Sale)
func (c *LandContract) GetAvailableLand(ctx contractapi.TransactionContextInterface) ([]*Land, error) {
	msp, _ := ctx.GetClientIdentity().GetMSPID()
	if msp != "Org3MSP" {
		return nil, fmt.Errorf("only Org3 (Buyer) can view available land")
	}

	query := `{"selector":{"status":"For Sale"}}`
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()

	var availableLands []*Land
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			continue
		}
		var land Land
		err = json.Unmarshal(queryResponse.Value, &land)
		if err != nil {
			continue
		}
		availableLands = append(availableLands, &land)
	}

	return availableLands, nil
}
