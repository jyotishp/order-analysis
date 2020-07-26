package main

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/tamerh/jsparser"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

var restaurant_count = make(map[string]int)
var cuisine_count = make(map[string]int)
var state_cuisine_count = make(map[string]map[string]int)

var secrets = gin.H{
	"shubham": gin.H{"email": "shubham.das2@swiggy.in", "phone": "7980365829"},
	"austin":  gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":    gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func getTopNumRestaurants(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		num := c.Param("num")
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

		numint, err := strconv.Atoi(num)
		if err == nil {
			if numint > len(ss) {
				numint = len(ss)
			}
			if numint >= 0 {
				c.JSON(200, ss[:numint])
			} else {
				numint = len(ss) + numint
				if numint < 0 {
					numint = 0
				}
				c.JSON(200, ss[numint:])
			}
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

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

func getTopNumStatesCuisines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {

		num := c.Param("num")
		state := c.Param("state")
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

		numint, err := strconv.Atoi(num)
		if err == nil {
			if numint > len(ss) {
				numint = len(ss)
			}
			if numint >= 0 {
				c.JSON(200, ss[:numint])
			} else {
				numint = len(ss) + numint
				if numint < 0 {
					numint = 0
				}
				c.JSON(200, ss[numint:])
			}
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

func getTopNumCuisines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := secrets[user]; ok {
		num := c.Param("num")
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

		numint, err := strconv.Atoi(num)
		if err == nil {
			if numint > len(ss) {
				numint = len(ss)
			}
			if numint >= 0 {
				c.JSON(200, ss[:numint])
			} else {
				numint = len(ss) + numint
				if numint < 0 {
					numint = 0
				}
				c.JSON(200, ss[numint:])
			}
		}
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
