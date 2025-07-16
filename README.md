# 🏡 Land Registry System using Hyperledger Fabric

A simplified blockchain-based Land Registry System built using **Hyperledger Fabric** with:
- ✅ 3 Organizations: Seller (Org1), Buyer (Org2), Government Registry (Org3)
- 🔐 Private Data Collections
- 💡 Chaincode in Go
- 🖥️ Gin-based Backend API
- 🌐 HTML + JS Frontend

---

## 📦 Project Structure

Land-Registry/
│
├── Chaincode/
│   └── Contracts/
│       ├── LandContract.go              # Main chaincode logic
│       ├── collections_config.json      # Private data configuration
│       └── META-INF/
│           └── statedb/
│               └── couchdb/
│                   └── indexCardId.json # CouchDB index (if needed)
│
├── network/                             # Optional: scripts or notes for fabric test-network
│
├── ui/                                  # Backend and UI
│   ├── client.go                        # Fabric Gateway client logic
│   ├── connect.go                       # Wallet & gateway setup
│   ├── main.go                          # Gin server entrypoint
│   ├── profile.go                       # Optional helper structs
│   ├── go.mod                           # Go module file
│   ├── go.sum                           # Go module checksum file
│
│   ├── static/                          # Public static files (JS, CSS, images)
│   │   └── index.js                     # JavaScript logic for frontend
│
│   └── templates/                       # HTML templates served by Gin
│       ├── index.html                   # UI with land operations
│       ├── seller.html                  # Optional: Seller-only UI
│       ├── buyer.html                   # Optional: Buyer-only UI
│       └── registry.html                # Optional: Registry-only UI
│
├── README.md                            # Full project documentation
├── .gitignore                           # Ignore node_modules, vendor, build, etc.
└── Land-Registry.tar.gz                 # (Optional) packaged chaincode file


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
