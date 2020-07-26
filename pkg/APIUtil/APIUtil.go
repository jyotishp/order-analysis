package APIUtil

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
)

var Secrets = gin.H{
	"shubham": gin.H{"email": "shubham.das2@swiggy.in", "phone": "7980365829"},
	"austin":  gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":    gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

type KV struct {
	Key   string
	Value int
}

func KeySort(count map[string] int, num string) []KV{
	var ss []KV
	for k, v := range count {
		ss = append(ss, KV{k, v})
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
			return ss[:numint]
		} else {
			numint = len(ss) + numint
			if numint < 0 {
				numint = 0
			}
			return ss[numint:]
		}
	}
	return nil
}

var Restaurant_count = make(map[string]int)
var Cuisine_count = make(map[string]int)
var State_cuisine_count = make(map[string]map[string]int)

func GetAllRestaurants(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := Secrets[user]; ok {
		c.JSON(200, Restaurant_count)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}

}

func GetAllCusines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := Secrets[user]; ok {
		c.JSON(200, Cuisine_count)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

func GetAllStatesCuisines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := Secrets[user]; ok {
		c.JSON(200, State_cuisine_count)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}


func GetTopNumRestaurants(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := Secrets[user]; ok {
		num := c.Param("num")
		jsonSlice:= KeySort(Restaurant_count, num)
		c.JSON(200,jsonSlice)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}



func GetTopNumStatesCuisines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := Secrets[user]; ok {

		num := c.Param("num")
		state := c.Param("state")
		jsonSlice:= KeySort(State_cuisine_count[state], num)
		c.JSON(200,jsonSlice)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

func GetTopNumCuisines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := Secrets[user]; ok {
		num := c.Param("num")
		jsonSlice:= KeySort(Cuisine_count, num)
		c.JSON(200,jsonSlice)
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}


