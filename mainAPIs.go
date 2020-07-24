package main

import (
	"bufio"
	"strconv"
	"sort"
	"github.com/tamerh/jsparser"
	"log"
	"os"
	"github.com/gin-gonic/gin"
)

var restaurant_count= make(map[string] int)
var cuisine_count = make(map[string] int)
var state_cuisine_count = make(map[string] map[string] int)


func getTopNumRestaurants(c *gin.Context){

	num:= c.Param("num")
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range restaurant_count {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	numint,err:=strconv.Atoi(num)
	if err == nil {
		if numint>len(ss){
			numint =len(ss)
		}
		if(numint>=0) {
			c.JSON(200, ss[:numint])
		} else {
			numint = len(ss)+numint
			if(numint<0){
				numint=0
			}
			c.JSON(200,ss[numint:])
		}
	}
}

func getAllRestaurants(c *gin.Context){

	c.JSON(200,restaurant_count)
}

func getAllCusines(c *gin.Context){

	c.JSON(200,cuisine_count)
}

func getAllStatesCuisines(c *gin.Context){

	c.JSON(200,state_cuisine_count)
}

func getTopNumStatesCuisines(c *gin.Context){

	num:= c.Param("num")
	state:= c.Param("state")
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range state_cuisine_count[state] {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	numint,err:=strconv.Atoi(num)
	if err == nil {
		if numint>len(ss){
			numint =len(ss)
		}
		if(numint>=0) {
			c.JSON(200, ss[:numint])
		} else {
			numint = len(ss)+numint
			if(numint<0){
				numint=0
			}
			c.JSON(200,ss[numint:])
		}
	}
}

func getTopNumCuisines(c *gin.Context){

	num:= c.Param("num")
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range cuisine_count {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	numint,err:=strconv.Atoi(num)
	if err == nil {
		if numint>len(ss){
			numint =len(ss)
		}
		if(numint>=0) {
			c.JSON(200, ss[:numint])
		} else {
			numint = len(ss)+numint
			if(numint<0){
				numint=0
			}
			c.JSON(200,ss[numint:])
		}
	}
}

func main() {
	router := gin.Default()
	restaurantAPI:=router.Group("/restaurant")
	restaurantAPI.GET("/all",getAllRestaurants)
	restaurantAPI.GET("/top/:num",getTopNumRestaurants)

	cuisineAPI:=router.Group("/cuisine")
	cuisineAPI.GET("/all",getAllCusines)
	cuisineAPI.GET("/top/:num",getTopNumCuisines)

	stateCuisineAPI:=router.Group("/state")
	stateCuisineAPI.GET("/all",getAllStatesCuisines)
	stateCuisineAPI.GET("/top/:state/:num",getTopNumStatesCuisines)

	r, _ := os.Open("outputs.json")
	br := bufio.NewReaderSize(r,65536)
	parser := jsparser.NewJSONParser(br, "orders")


	for json:= range parser.Stream() {
		if json.Err != nil {
			log.Fatal(json.Err)
		}
		//fmt.Println(json.ObjectVals["OrderId"])
		restaurant:= json.ObjectVals["RestName"]
		cuisine:=json.ObjectVals["Cuisine"]
		state:=json.ObjectVals["State"]

		restaurant_count[restaurant.(string)]++
		cuisine_count[cuisine.(string)]++
		statemap,ok:=state_cuisine_count[state.(string)]
		if ok{
			statemap[cuisine.(string)]++
		} else{
			state_cuisine_count[state.(string)]= make(map[string] int)
			state_cuisine_count[state.(string)][cuisine.(string)]++
		}
	}
	router.Run(":5665")

}
