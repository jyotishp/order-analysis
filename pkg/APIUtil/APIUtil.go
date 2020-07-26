package APIUtil

import (
	"github.com/gin-gonic/gin"
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

