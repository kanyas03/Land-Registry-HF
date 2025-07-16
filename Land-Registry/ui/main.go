package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Allow requests from browser frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5500"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Serve static JS (index.js) and HTML
	router.Static("/static", "./static")           // for /static/index.js
	router.Static("/ui", "./templates")            // for /ui/index.html
	router.GET("/", func(c *gin.Context) {         // redirect to index.html
		c.Redirect(http.StatusMovedPermanently, "/ui/index.html")
	})

	// ========== API ENDPOINTS ==========

	// Org1 - List Land
	router.POST("/api/list-land", func(c *gin.Context) {
		var land struct {
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
		}

		if err := c.BindJSON(&land); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		result := submitTxnFn("org1", "autochannel", "Land-Registry", "LandContract", "invoke",
			map[string][]byte{},
			"ListLand",
			land.LandID, land.Location, land.Size, land.Type, land.SoilQuality,
			land.WaterSource, land.NearbyRoad, land.NearbyCity, land.Coordinates, land.SellingPrice,
		)

		c.String(http.StatusOK, result)
	})

	// Org2 - Get Available Lands
	router.GET("/api/get-available-lands", func(c *gin.Context) {
		result := submitTxnFn("org2", "autochannel", "Land-Registry", "LandContract", "query",
			map[string][]byte{}, "GetAvailableLands")

		var parsed []map[string]interface{}
		if err := json.Unmarshal([]byte(result), &parsed); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse lands"})
			return
		}

		c.JSON(http.StatusOK, parsed)
	})

	// Org2 - Request to Buy
	router.POST("/api/request-buy", func(c *gin.Context) {
		var body struct {
			OfferID      string            `json:"offerID"`
			BuyerRequest map[string]string `json:"buyerRequest"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		privateData := map[string][]byte{
			"buyerRequest": encodeJSONBytes(body.BuyerRequest),
		}

		result := submitTxnFn("org2", "autochannel", "Land-Registry", "LandContract", "private",
			privateData, "RequestToBuy", body.OfferID)

		c.String(http.StatusOK, result)
	})

	// Org3 - Register to Buyer
	router.POST("/api/register-buyer", func(c *gin.Context) {
		var body struct {
			LandID         string            `json:"landID"`
			BuyerOwnership map[string]string `json:"buyerOwnership"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		privateData := map[string][]byte{
			"buyerOwnership": encodeJSONBytes(body.BuyerOwnership),
		}

		result := submitTxnFn("org3", "autochannel", "Land-Registry", "LandContract", "private",
			privateData, "RegisterToBuyer", body.LandID)

		c.String(http.StatusOK, result)
	})

	// Start server on localhost:3001
	router.Run("localhost:3001")
}

// Utility function for transient data
func encodeJSONBytes(data map[string]string) []byte {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic("Failed to encode transient data: " + err.Error())
	}
	return jsonBytes
}
