package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Land struct {
	LandID   string `json:"landId"`
	Location string `json:"location"`
	Size     string `json:"size"`
	Owner    string `json:"owner"`
}

func main() {
	router := gin.Default()
	router.Static("/public", "./public")
	router.LoadHTMLGlob("templates/*")
	router.Use(cors.Default())

	// Home
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// 1. Org1 - Register Land
	router.POST("/api/land", func(ctx *gin.Context) {
		var req Land
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		fmt.Println("Registering land:", req)

		output := submitTxnFn("org1", "autochannel", "Land-Registry", "LandContract", "invoke", nil,
			"RegisterLand", req.LandID, req.Location, req.Size, req.Owner)

		fmt.Println(output)
		ctx.JSON(http.StatusOK, gin.H{"message": "Land registered"})
	})

	// 2. Org2 - Mark land for sale
	router.POST("/api/land/sell", func(ctx *gin.Context) {
		var req struct {
			LandID   string `json:"landId"`
			SellerID string `json:"sellerId"`
		}
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		output := submitTxnFn("org2", "autochannel", "Land-Registry", "LandContract", "invoke", nil,
			"MarkForSale", req.LandID, req.SellerID)

		fmt.Println(output)
		ctx.JSON(http.StatusOK, gin.H{"message": "Land marked for sale"})
	})

	// 3. Org3 - View Available Lands
	router.GET("/api/land/available", func(ctx *gin.Context) {
		result := submitTxnFn("org3", "autochannel", "Land-Registry", "LandContract", "query", nil,
			"GetAvailableLand")

		ctx.JSON(http.StatusOK, gin.H{"data": result})
	})

	// 4. Org3 - Request Purchase
	router.POST("/api/land/request", func(ctx *gin.Context) {
		var req struct {
			LandID  string `json:"landId"`
			BuyerID string `json:"buyerId"`
		}
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		output := submitTxnFn("org3", "autochannel", "Land-Registry", "LandContract", "invoke", nil,
			"RequestPurchase", req.LandID, req.BuyerID)

		fmt.Println(output)
		ctx.JSON(http.StatusOK, gin.H{"message": "Purchase requested"})
	})

	// 5. Org1 - Approve Purchase
	router.POST("/api/land/approve", func(ctx *gin.Context) {
		var req struct {
			LandID string `json:"landId"`
		}
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
			return
		}

		output := submitTxnFn("org1", "autochannel", "Land-Registry", "LandContract", "invoke", nil,
			"ApprovePurchase", req.LandID)

		fmt.Println(output)
		ctx.JSON(http.StatusOK, gin.H{"message": "Ownership transferred"})
	})

	// 6. Get land info by ID
	router.GET("/api/land/:landId", func(ctx *gin.Context) {
		landId := ctx.Param("landId")
		fmt.Println("Fetching land:", landId)

		result := submitTxnFn("org1", "autochannel", "Land-Registry", "LandContract", "query", nil,
			"GetLand", landId)

		fmt.Println("Result from chaincode:", result)

		if result == "" || result == "null" {
			fmt.Println("Result is empty or null.")
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Land not found"})
			return
		}
		
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(result), &parsed); err != nil {
			fmt.Println("JSON unmarshal error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse response"})
			return
		}
		
		ctx.JSON(http.StatusOK, gin.H{"data": parsed})
		
	})

	router.Run("localhost:3001")
}
