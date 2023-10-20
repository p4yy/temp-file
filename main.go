package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"}
	// config.AllowOrigins = []string{"http://localhost:3000"}
	config.AddAllowHeaders("X-Auth-Token")

	indexHTML, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Printf("Failed to read index.html: %v\n", err)
		os.Exit(1)
	}

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	r.POST("/cariPemain", func(c *gin.Context) {
		var requestData struct {
			MatchID string `json:"matchId"`
		}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}

		url := fmt.Sprintf("https://api.football-data.org/v4/persons/%s", requestData.MatchID)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
			return
		}
		req.Header.Set("X-Auth-Token", "37db945bec9c4eadb015167f5e8fed1d") // api key

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			// data json
			var responseData map[string]interface{}
			err := json.NewDecoder(resp.Body).Decode(&responseData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response: " + err.Error()})
				return
			}

			c.JSON(http.StatusOK, responseData)
		} else {
			c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("Request failed with status code %d", resp.StatusCode)})
		}
	})

	r.Run(":8080")
}
