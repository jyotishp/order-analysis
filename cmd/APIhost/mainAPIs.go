package main

import (
	"github.com/gin-gonic/gin"
	//"github.com/shubham491/order-analysis/pkg/Grpc_server"
	"github.com/shubham491/order-analysis/pkg/AuthUtil"
	"github.com/shubham491/order-analysis/pkg/Grpc_client"
)

func addAPIPaths(router *gin.Engine){
	restaurantAPI := router.Group("/restaurant", gin.BasicAuth(AuthUtil.Accounts))
	cuisineAPI := router.Group("/cuisine", gin.BasicAuth(AuthUtil.Accounts))
	stateCuisineAPI := router.Group("/state", gin.BasicAuth(AuthUtil.Accounts))
	orderAPI := router.Group("/order", gin.BasicAuth(AuthUtil.Accounts))

	//restaurantAPI:=router.Group("/restaurant")
	restaurantAPI.GET("/all", Grpc_client.GetAllRestaurants)
	restaurantAPI.GET("/top/:num", Grpc_client.GetTopNumRestaurants)

	//cuisineAPI:=router.Group("/cuisine")
	cuisineAPI.GET("/all", Grpc_client.GetAllCusines)
	cuisineAPI.GET("/top/:num", Grpc_client.GetTopNumCuisines)

	//stateCuisineAPI:=router.Group("/state")
	stateCuisineAPI.GET("/all", Grpc_client.GetAllStatesCuisines)
	stateCuisineAPI.GET("/top/:state/:num", Grpc_client.GetTopNumStateCuisines)

	orderAPI.POST("/add",Grpc_client.AddOrder)
}

func main() {
	router := gin.Default()

	addAPIPaths(router)


	router.Run(":5665")

}
