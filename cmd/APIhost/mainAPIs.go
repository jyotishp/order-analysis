package main

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/tamerh/jsparser"
	"log"
	"net/http"
	"os"
)

var restaurant_count = make(map[string]int)
var cuisine_count = make(map[string]int)
var state_cuisine_count = make(map[string]map[string]int)

func getAllRestaurants(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		c.JSON(200, restaurant_count)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}

}

func getAllCusines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		c.JSON(200, cuisine_count)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

func getAllStatesCuisines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		c.JSON(200, state_cuisine_count)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}


func getTopNumRestaurants(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		num := c.Param("num")
		jsonSlice:= keySort(restaurant_count, num)
		c.JSON(200,jsonSlice)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}



func getTopNumStatesCuisines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {

		num := c.Param("num")
		state := c.Param("state")
		jsonSlice:= keySort(state_cuisine_count[state], num)
		c.JSON(200,jsonSlice)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

func getTopNumCuisines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		num := c.Param("num")
		jsonSlice:= keySort(cuisine_count, num)
		c.JSON(200,jsonSlice)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

func main() {
	router := gin.Default()

	accounts := gin.Accounts{
		"shubham": "das",
		"austin":  "1234",
		"lena":    "hello2",
		"manu":    "4321",
	}

	restaurantAPI := router.Group("/restaurant", gin.BasicAuth(accounts))
	cuisineAPI := router.Group("/cuisine", gin.BasicAuth(accounts))
	stateCuisineAPI := router.Group("/state", gin.BasicAuth(accounts))

	//restaurantAPI:=router.Group("/restaurant")
	restaurantAPI.GET("/all", getAllRestaurants)
	restaurantAPI.GET("/top/:num", getTopNumRestaurants)

	//cuisineAPI:=router.Group("/cuisine")
	cuisineAPI.GET("/all", getAllCusines)
	cuisineAPI.GET("/top/:num", getTopNumCuisines)

	//stateCuisineAPI:=router.Group("/state")
	stateCuisineAPI.GET("/all", getAllStatesCuisines)
	stateCuisineAPI.GET("/top/:state/:num", getTopNumStatesCuisines)

	r, _ := os.Open("outputs.json")
	br := bufio.NewReaderSize(r, 65536)
	parser := jsparser.NewJSONParser(br, "orders")

	for json := range parser.Stream() {
		if json.Err != nil {
			log.Fatal(json.Err)
		}
		//fmt.Println(json.ObjectVals["OrderId"])
		restaurant := json.ObjectVals["RestName"]
		cuisine := json.ObjectVals["Cuisine"]
		state := json.ObjectVals["State"]

		restaurant_count[restaurant.(string)]++
		cuisine_count[cuisine.(string)]++
		statemap, ok := state_cuisine_count[state.(string)]
		if ok {
			statemap[cuisine.(string)]++
		} else {
			state_cuisine_count[state.(string)] = make(map[string]int)
			state_cuisine_count[state.(string)][cuisine.(string)]++
		}
	}
	router.Run(":5665")

}
