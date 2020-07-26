package main

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/tamerh/jsparser"
	"log"
	"os"
	"github.com/jyotishp/order-analysis/pkg/APIUtil"
	"github.com/jyotishp/order-analysis/pkg/AuthUtil"
)

func addAPIPaths(router *gin.Engine){
	restaurantAPI := router.Group("/restaurant", gin.BasicAuth(AuthUtil.Accounts))
	cuisineAPI := router.Group("/cuisine", gin.BasicAuth(AuthUtil.Accounts))
	stateCuisineAPI := router.Group("/state", gin.BasicAuth(AuthUtil.Accounts))

	//restaurantAPI:=router.Group("/restaurant")
	restaurantAPI.GET("/all", APIUtil.GetAllRestaurants)
	restaurantAPI.GET("/top/:num", APIUtil.GetTopNumRestaurants)

	//cuisineAPI:=router.Group("/cuisine")
	cuisineAPI.GET("/all", APIUtil.GetAllCusines)
	cuisineAPI.GET("/top/:num", APIUtil.GetTopNumCuisines)

	//stateCuisineAPI:=router.Group("/state")
	stateCuisineAPI.GET("/all", APIUtil.GetAllStatesCuisines)
	stateCuisineAPI.GET("/top/:state/:num", APIUtil.GetTopNumStatesCuisines)
}

func main() {
	router := gin.Default()

	addAPIPaths(router)

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

		APIUtil.Restaurant_count[restaurant.(string)]++
		APIUtil.Cuisine_count[cuisine.(string)]++
		statemap, ok := APIUtil.State_cuisine_count[state.(string)]
		if ok {
			statemap[cuisine.(string)]++
		} else {
			APIUtil.State_cuisine_count[state.(string)] = make(map[string]int)
			APIUtil.State_cuisine_count[state.(string)][cuisine.(string)]++
		}
	}
	router.Run(":5665")

}
