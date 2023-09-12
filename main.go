package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/gin-gonic/gin"
	// "github.com/go-oauth2/oauth2/v4/manage"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func main() {

	router := gin.Default()

	router.POST("/oauth2/token", Login)

	port := "8068"
	// log.Info().Msg("Starting server on :" + port)
	router.Run(":" + port)

}

func Login(c *gin.Context) {

	twitterToken := "1488730958270726144-E7EGFotCFgNTYpm8vrVelEnVUsQHHO"
	twitterTokenSecret := "tP4mzrhOEp0ChPoZP38RVwAAWkHn5ZfDoqtypc322MwME"

	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET_KEY"))
	token := oauth1.NewToken(twitterToken, twitterTokenSecret)
	
	httpClient := config.Client(oauth1.NoContext, token)
	type TwitterResponseStatus struct {
		Place struct {
			Country     string `json:"country"`
			CountryCode string `json:"country_code"`
		} `json:"place"`
	}

	type TwitterResponse struct {
		ID     int64                 `json:"id"`
		Email  string                `json:"email"`
		Name   string                `json:"name"`
		Status TwitterResponseStatus `json:"status"`
	}

	// example Twitter API request
	path := "https://api.twitter.com/1.1/account/verify_credentials.json?include_email=true&skip_status=false"
	resp, _ := httpClient.Get(path)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("Raw Response Body:\n%v\n", string(body))

	var user map[string]interface{}
	if err := json.Unmarshal(body, &user); err != nil {
		log.Fatalf("Failed to parse JSON response: %v", err)
		return
	}

	c.JSON(200, user)
}
