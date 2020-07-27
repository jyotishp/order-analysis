package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tamerh/jsparser"
	"log"
	"os"
	"github.com/jyotishp/order-analysis/pkg/APIUtil"
	"github.com/jyotishp/order-analysis/pkg/AuthUtil"
	"github.com/jyotishp/order-analysis/pkg/Grpc_client"
)

func addAPIPaths(router *gin.Engine){
	restaurantAPI := router.Group("/restaurant", gin.BasicAuth(AuthUtil.Accounts))
	cuisineAPI := router.Group("/cuisine", gin.BasicAuth(AuthUtil.Accounts))
	stateCuisineAPI := router.Group("/state", gin.BasicAuth(AuthUtil.Accounts))
	orderAPI := router.Group("/order", gin.BasicAuth(AuthUtil.Accounts))

	//restaurantAPI:=router.Group("/restaurant")
	restaurantAPI.GET("/all", Grpc_client.GetAllRestaurants)
	restaurantAPI.GET("/top/:num", APIUtil.GetTopNumRestaurants)

	//cuisineAPI:=router.Group("/cuisine")
	cuisineAPI.GET("/all", Grpc_client.GetAllCusines)
	cuisineAPI.GET("/top/:num", APIUtil.GetTopNumCuisines)

	//stateCuisineAPI:=router.Group("/state")
	stateCuisineAPI.GET("/all", Grpc_client.GetAllStatesCuisines)
	stateCuisineAPI.GET("/top/:state/:num", APIUtil.GetTopNumStatesCuisines)

	orderAPI.POST("/add",APIUtil.AddOrder)
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
		id := json.ObjectVals["Id"]

		APIUtil.Restaurant_count[restaurant.(string)]++
		APIUtil.Cuisine_count[cuisine.(string)]++
		APIUtil.Orders[id.(string)]++
		statemap, ok := APIUtil.State_cuisine_count[state.(string)]
		if ok {
			statemap[cuisine.(string)]++
		} else {
			APIUtil.State_cuisine_count[state.(string)] = make(map[string]int)
			APIUtil.State_cuisine_count[state.(string)][cuisine.(string)]++
		}
	}
	fmt.Println(APIUtil.Orders["2999999"])
	router.Run(":5665")

}
