package ErrorHandlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
)

func ParseInt(txt string) int {
	n, err := strconv.Atoi(txt)
	FatalErr(err, "Error converting to integer")
	return n
}

func ParseFloat(txt string) float64 {
	n, err := strconv.ParseFloat(txt, 64)
	FatalErr(err, "Error converting to float")
	return n
}

func FatalErr(err error, txt string) {
	if err != nil {
		log.Fatalf("%v - %v", txt, err.Error())
	}
}

func InfoErr(err error, txt string) {
	if err != nil {
		log.Printf("%v - %v", txt, err.Error())
	}
}

func HandleErr(err error, txt string) {
	if err != nil {
		log.Fatal(txt, err.Error())
	}
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CheckError(err error, c *gin.Context)  {
	if err != nil {
		c.JSON(200,gin.H{
			"error":err.Error(),
		})
	}
}
