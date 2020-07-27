package Grpc_client

import (
	"github.com/gin-gonic/gin"
	"github.com/jyotishp/order-analysis/pkg/AuthUtil"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func GetAllRestaurants(c *gin.Context) {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Sorry client cannot talk to server: %v: ", err)
	}
	defer conn.Close();
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := AuthUtil.Secrets[user]; ok {
		c.JSON(200, Restaurant_count)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}

}

func GetAllCusines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := AuthUtil.Secrets[user]; ok {
		c.JSON(200, Cuisine_count)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

func GetAllStatesCuisines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := AuthUtil.Secrets[user]; ok {
		c.JSON(200, State_cuisine_count)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}