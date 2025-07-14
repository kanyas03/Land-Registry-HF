// SPDX-License-Identifier: Apache-2.0
package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type LandContract struct {
	contractapi.Contract
}

type Land struct {
	LandID       string `json:"landID"`
	Location     string `json:"location"`
	Size         string `json:"size"`
	Type         string `json:"type"`
	SoilQuality  string `json:"soilQuality"`
	WaterSource  string `json:"waterSource"`
	NearbyRoad   string `json:"nearbyRoad"`
	NearbyCity   string `json:"nearbyCity"`
	Coordinates  string `json:"coordinates"`
	SellingPrice string `json:"sellingPrice"`
	Status       string `json:"status"` // For Sale, Sold
}

type BuyerOwnership struct {
	OwnerID      string `json:"ownerID"`
	BuyerName    string `json:"buyerName"`
	Aadhar       string `json:"aadhar"`
	DocumentHash string `json:"documentHash"`
	TransferDate string `json:"transferDate"`
	LandID       string `json:"landID"`
	Location     string `json:"location"`
	Size         string `json:"size"`
	Type         string `json:"type"`
	Coordinates  string `json:"coordinates"`
	SellingPrice string `json:"sellingPrice"`
}

// Org1 Seller lists land to public ledger
func (c *LandContract) ListLand(ctx contractapi.TransactionContextInterface, landID string, location string, size string, landType string, soilQuality string, waterSource string, nearbyRoad string, nearbyCity string, coordinates string, sellingPrice string) error {
	msp, _ := ctx.GetClientIdentity().GetMSPID()
	if msp != "Org1MSP" {
		return fmt.Errorf("only Seller (Org1) can list land")
	}

	existing, err := ctx.GetStub().GetState(landID)
	if err != nil {
		return fmt.Errorf("failed to read land from world state: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("land with ID %s already exists", landID)
	}

	land := Land{
		LandID:       landID,
		Location:     location,
		Size:         size,
		Type:         landType,
		SoilQuality:  soilQuality,
		WaterSource:  waterSource,
		NearbyRoad:   nearbyRoad,
		NearbyCity:   nearbyCity,
		Coordinates:  coordinates,
		SellingPrice: sellingPrice,
		Status:       "For Sale",
	}

	landJSON, err := json.Marshal(land)
	if err != nil {
		return fmt.Errorf("failed to marshal land: %v", err)
	}

	err = ctx.GetStub().PutState(landID, landJSON)
	if err != nil {
		return fmt.Errorf("failed to write land to public ledger: %v", err)
	}

	return nil
}

// Anyone (e.g., Org1, Org2, Org3) can get public land info
func (c *LandContract) GetLandByID(ctx contractapi.TransactionContextInterface, landID string) (*Land, error) {
	landBytes, err := ctx.GetStub().GetState(landID)
	if err != nil {
		return nil, fmt.Errorf("failed to read land from world state: %v", err)
	}
	if landBytes == nil {
		return nil, fmt.Errorf("land with ID %s does not exist", landID)
	}

	var land Land
	err = json.Unmarshal(landBytes, &land)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling land data: %v", err)
	}

	return &land, nil
}

// Buyer (Org2) views lands that are For Sale
func (c *LandContract) GetAvailableLands(ctx contractapi.TransactionContextInterface) ([]*Land, error) {
	msp, _ := ctx.GetClientIdentity().GetMSPID()
	if msp != "Org2MSP" {
		return nil, fmt.Errorf("only Buyer (Org2) can view available lands")
	}

	query := `{"selector":{"status":"For Sale"}}`
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query available lands: %v", err)
	}
	defer resultsIterator.Close()

	var lands []*Land
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var land Land
		err = json.Unmarshal(queryResponse.Value, &land)
		if err != nil {
			return nil, err
		}
		lands = append(lands, &land)
	}

	return lands, nil
}

// Buyer (Org2) sends private request to buy land
func (c *LandContract) RequestToBuy(ctx contractapi.TransactionContextInterface, offerID string) error {
	msp, _ := ctx.GetClientIdentity().GetMSPID()
	if msp != "Org2MSP" {
		return fmt.Errorf("only Buyer (Org2) can send requests")
	}

	transient, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient data: %v", err)
	}
	privateData, ok := transient["buyerRequest"]
	if !ok {
		return fmt.Errorf("buyerRequest key missing in transient data")
	}

	err = ctx.GetStub().PutPrivateData("collectionBuyerSeller", offerID, privateData)
	if err != nil {
		return fmt.Errorf("failed to store buyer request: %v", err)
	}

	return nil
}

// Land Registry (Org3) assigns land to buyer and stores private ownership
func (c *LandContract) RegisterToBuyer(ctx contractapi.TransactionContextInterface, landID string) (string, error) {
	msp, _ := ctx.GetClientIdentity().GetMSPID()
	if msp != "Org3MSP" {
		return "", fmt.Errorf("only LandRegistry (Org3) can register land to buyer")
	}

	landBytes, err := ctx.GetStub().GetState(landID)
	if err != nil || landBytes == nil {
		return "", fmt.Errorf("landID not found or error: %v", err)
	}

	var land Land
	err = json.Unmarshal(landBytes, &land)
	if err != nil {
		return "", fmt.Errorf("error parsing land data: %v", err)
	}

	land.Status = "Sold"

	updatedJSON, err := json.Marshal(land)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated land: %v", err)
	}

	err = ctx.GetStub().PutState(landID, updatedJSON)
	if err != nil {
		return "", fmt.Errorf("failed to update land status: %v", err)
	}

	transient, err := ctx.GetStub().GetTransient()
	if err != nil {
		return "", fmt.Errorf("error getting transient data: %v", err)
	}
	privateOwnerData, ok := transient["buyerOwnership"]
	if !ok {
		return "", fmt.Errorf("buyerOwnership key missing in transient")
	}

	err = ctx.GetStub().PutPrivateData("collectionBuyerLandRegistry", landID, privateOwnerData)
	if err != nil {
		return "", fmt.Errorf("failed to store private ownership data: %v", err)
	}

	var cert BuyerOwnership
	err = json.Unmarshal(privateOwnerData, &cert)
	if err != nil {
		return "", fmt.Errorf("failed to parse ownership certificate: %v", err)
	}

	certificate := fmt.Sprintf(`--- Ownership Certificate ---
Land ID: %s
Location: %s
Size: %s
Type: %s
Price: %s
Owner: %s (%s)
Transfer Date: %s
----------------------------`,
		cert.LandID, cert.Location, cert.Size, cert.Type, cert.SellingPrice, cert.BuyerName, cert.OwnerID, cert.TransferDate)

	return certificate, nil
}
