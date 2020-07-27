package APIUtil

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"

	//"github.com/jyotishp/order-analysis/pkg/ErrorHandlers"
	"github.com/jyotishp/order-analysis/pkg/Models"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"github.com/jyotishp/order-analysis/pkg/AuthUtil"
)

var Restaurant_count = make(map[string]int)
var Cuisine_count = make(map[string]int)
var State_cuisine_count = make(map[string]map[string]int)
var Orders = make(map[string] int)

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



func GetAllRestaurants(c *gin.Context) {

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


func GetTopNumRestaurants(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := AuthUtil.Secrets[user]; ok {
		num := c.Param("num")
		jsonSlice:= KeySort(Restaurant_count, num)
		if jsonSlice == nil{
			c.JSON(200,gin.H{
				"Error":"Provide valid integer value.",
			})
		} else {
			c.JSON(200, jsonSlice)
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}



func GetTopNumStatesCuisines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := AuthUtil.Secrets[user]; ok {

		num := c.Param("num")
		state := c.Param("state")
		jsonSlice:= KeySort(State_cuisine_count[state], num)
		if jsonSlice == nil{
			c.JSON(200,gin.H{
				"Error":"Provide valid integer value.",
			})
		} else {
			c.JSON(200, jsonSlice)
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

func GetTopNumCuisines(c *gin.Context) {

	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := AuthUtil.Secrets[user]; ok {
		num := c.Param("num")
		jsonSlice:= KeySort(Cuisine_count, num)
		if jsonSlice == nil{
			c.JSON(200,gin.H{
				"Error":"Provide valid integer value.",
			})
		} else {
			c.JSON(200, jsonSlice)
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}

func CheckError(err error, c *gin.Context)  {
	if err != nil {
		c.JSON(200,gin.H{
			"error":err.Error(),
		})
	}
}

func AddOrder(c *gin.Context){
	user := c.MustGet(gin.AuthUserKey).(string)
	if _, ok := AuthUtil.Secrets[user]; ok {
	body:=c.Request.Body
	content, _:= ioutil.ReadAll(body)
	var orderData Models.Order
	//var orderData2 Models.Order
	err := json.Unmarshal([]byte(content), &orderData)
	CheckError(err,c)
	err = json.Unmarshal(content, &orderData)
	Id := fmt.Sprint(orderData.Id)
	fmt.Println(Id)
	if Orders[string(Id)] >= 1{
		c.JSON(200, gin.H{
			"Error":"Order ID already there",
		})
		return
	}

	f, err := os.OpenFile("outputs.json", os.O_RDWR, os.ModePerm)
	defer f.Close()
	CheckError(err,c)

	orderJson, err := json.Marshal(orderData)
	CheckError(err,c)

	orderString := string(orderJson)
	orderString = "," + orderString

	off := int64(2)
	stat, err := os.Stat("outputs.json")
	fmt.Println("Size : ", stat.Size())
	start := stat.Size() - off

	tmp := []byte(orderString)
	_, err = f.WriteAt(tmp, start)
	CheckError(err, c)

	str := []byte("]}")
	_, err = f.WriteAt(str, start + int64(len(orderString)))
	CheckError(err, c)

	restaurant := orderData.RestName
	cuisine := orderData.Cuisine
	state := orderData.State

	Restaurant_count[restaurant]++
	Cuisine_count[cuisine]++
	Orders[string(Id)]++
	statemap, ok := State_cuisine_count[state]
	if ok {
		statemap[cuisine]++
	} else {
		State_cuisine_count[state] = make(map[string]int)
		State_cuisine_count[state][cuisine]++
	}

	c.JSON(200,gin.H{
		"success":"order successfully added",
	})

	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	}
}


