# 🏡 Land Registry System using Hyperledger Fabric

A simplified blockchain-based Land Registry System built using **Hyperledger Fabric** with:
- ✅ 3 Organizations: Seller (Org1), Buyer (Org2), Government Registry (Org3)
- 🔐 Private Data Collections
- 💡 Chaincode in Go
- 🖥️ Gin-based Backend API
- 🌐 HTML + JS Frontend
---

## 🧠 Project Logic

### Roles:
- **Seller (Org1):** Lists land for sale.
- **Buyer (Org2):** Requests to buy land.
- **Government Registry (Org3):** Finalizes ownership transfer.

### Data Handling:
- **Public Ledger:** Stores land ID, status, etc.
- **Private Data:**
  - `collectionSellerLandRegistry`: Between Seller & Registry
  - `collectionBuyerSeller`: Between Buyer & Seller
  - `collectionBuyerLandRegistry`: Between Buyer & Registry (ownership transfer)

---

## ⚙️ Setup Instructions

### 1. Clone & Setup Network
```bash
    git clone https://github.com/hyperledger/fabric-samples.git
    cd fabric-samples/test-network
```
---
### Ensure you have:

1. Docker

2. Go (v1.20+)

3. Fabric binaries

## Backend Function setup

### First step up the network
```bash

./network.sh up createChannel -c autochannel -ca -s couchdb


## Up the org 3 
    
    cd addOrg3
    
    ./addOrg3.sh up -c autochannel -ca -s couchdb
    
    cd  ..


## For deploy The chainCode

    ./network.sh deployCC -ccn Land-Registry -ccp ../../Land-Registry/Chaincode/ -ccl go -c autochannel -cccg ../../Land-Registry/Chaincode/collections.json 

```
---
### Frontend Set up 

From the Land-Registry/ui/ folder:
```bash
    go mod tidy
    go run main.go
```
