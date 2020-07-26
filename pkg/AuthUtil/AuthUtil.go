package AuthUtil

import "github.com/gin-gonic/gin"

var Secrets = gin.H{
	"shubham": gin.H{"email": "shubham.das2@swiggy.in", "phone": "7980365829"},
	"austin":  gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":    gin.H{"email": "lena@guapa.com", "phone": "523443"},
}
var Accounts = gin.Accounts{
"shubham": "das",
"austin":  "1234",
"lena":    "hello2",
"manu":    "4321",
}
