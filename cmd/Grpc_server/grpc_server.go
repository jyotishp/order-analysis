package main

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/jyotishp/order-analysis/pkg/APIUtil"
	"github.com/tamerh/jsparser"

	//"github.com/gin-gonic/gin"
	"github.com/shubham491/order-analysis/pkg/Models"
	"os"

	//"github.com/gin-gonic/gin"
	//"net/http"

	//"encoding/json"
	"fmt"
	//"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"

	//"github.com/jyotishp/order-analysis/pkg/ErrorHandlers"
	//"github.com/shubham491/order-analysis/pkg/Models"
	//"io/ioutil"
	//"net/http"
	//"os"
	"sort"
	"strconv"
	//"github.com/shubham491/order-analysis/pkg/AuthUtil"
	"github.com/shubham491/order-analysis/pkg/services/orders/orderspb"
)

var Restaurant_count = make(map[string] int64)
var Cuisine_count = make(map[string] int64)
var State_cuisine_count = make(map[string]map[string]int64)
var Orders = make(map[string] int64)

type KV struct {
	Key   string
	Value int64
}

type OrdersServiceServer struct {

}


func KeySort(count map[string] int64, num string) []KV{
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



func (s *OrdersServiceServer) GetAllRestaurant(ctx context.Context, request *orderspb.AllRestaurantRequest) (*orderspb.AllRestaurantResponse, error) {
	res:=&orderspb.AllRestaurantResponse{AllRestaurant:Restaurant_count}
	return res, nil
}
func (s *OrdersServiceServer) GetAllCuisine(ctx context.Context, request *orderspb.AllCuisineRequest) (*orderspb.AllCuisineResponse, error) {
	res:=&orderspb.AllCuisineResponse{AllCuisine:Cuisine_count}
	return res, nil
}

func (s *OrdersServiceServer) GetAllStateCusine(ctx context.Context, request *orderspb.AllStateRequest) (*orderspb.AllStateResponse, error) {
	var res *orderspb.AllStateResponse
	var res1 *orderspb.AllCuisine
	var tempMap=make(map[string]*orderspb.AllCuisine)
	for k,v:= range State_cuisine_count{
		res1= &orderspb.AllCuisine{
			AllCuisine:v,
		}

		tempMap[k]=res1

	}
	res=&orderspb.AllStateResponse{AllState:tempMap}
	return res, nil
}


func (s *OrdersServiceServer) GetTopNumRestaurants(c context.Context, request *orderspb.TopNumRestaurantRequest) (*orderspb.TopNumRestaurantResponse, error) {
	jsonSlice:= KeySort(Restaurant_count, request.Num)
	var kv=make(map[string] string)
	for _,v:= range jsonSlice{
		kv[v.Key]=fmt.Sprint(v.Value)
	}
	res:=&orderspb.TopNumRestaurantResponse{TopNumRestaurant:kv}
	return res, nil
}

func (s *OrdersServiceServer) GetTopNumCuisines(c context.Context, request *orderspb.TopNumCuisineRequest) (*orderspb.TopNumCuisineResponse, error) {
	jsonSlice:= KeySort(Cuisine_count, request.Num)
	var kv=make(map[string] string)
	for _,v:= range jsonSlice{
		kv[v.Key]=fmt.Sprint(v.Value)
	}
	res:=&orderspb.TopNumCuisineResponse{TopNumCuisine:kv}
	return res, nil
}

func (s *OrdersServiceServer) GetTopNumStatesCuisines(c context.Context, request *orderspb.TopNumStatesCuisinesRequest) (*orderspb.TopNumStatesCuisinesResponse, error) {
	jsonSlice:= KeySort(State_cuisine_count[request.State], request.Num)
	var kv=make(map[string] string)
	for _,v:= range jsonSlice{
		kv[v.Key]=fmt.Sprint(v.Value)
	}
	res:=&orderspb.TopNumStatesCuisinesResponse{TopNumState:kv}
	return res, nil
}

func CheckError(err error)  (*orderspb.AddOrderResponse){
	tempMap:=make(map[string] string)

	if err != nil {
		tempMap["error"]=err.Error()
		res:=&orderspb.AddOrderResponse{
			Response:tempMap,
		}
		return res
	}
	return nil
}

func (s *OrdersServiceServer) AddOrder(c context.Context, request *orderspb.AddOrderRequest) (*orderspb.AddOrderResponse, error) {
	var orderData Models.Order
	//var orderData2 Models.Order
	err := json.Unmarshal([]byte(request.Order), &orderData)
	res:=CheckError(err)
	if res!=nil{
		return res,nil
	}
	Id := fmt.Sprint(orderData.Id)
	fmt.Println(Id)
	//fmt.Println(Orders)
	if Orders[Id] >= 1{
		tempMap:=make(map[string] string)
		tempMap["error"]=fmt.Sprintf("OrderId %v already present",Id)
		res:=&orderspb.AddOrderResponse{Response: tempMap}
		return res, nil
	}

	f, err := os.OpenFile("outputs.json", os.O_RDWR, os.ModePerm)
	defer f.Close()
	res=CheckError(err)
	if res!=nil{
		return res,nil
	}

	orderJson, err := json.Marshal(orderData)
	res=CheckError(err)
	if res!=nil{
		return res,nil
	}

	orderString := string(orderJson)
	orderString = "," + orderString

	off := int64(2)
	stat, err := os.Stat("outputs.json")
	fmt.Println("Size : ", stat.Size())
	start := stat.Size() - off

	tmp := []byte(orderString)
	_, err = f.WriteAt(tmp, start)
	res=CheckError(err)
	if res!=nil{
		return res,nil
	}

	str := []byte("]}")
	_, err = f.WriteAt(str, start + int64(len(orderString)))
	res=CheckError(err)
	if res!=nil{
		return res,nil
	}
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
		tempmap:=make(map[string] int64)
		State_cuisine_count[state] = tempmap
		State_cuisine_count[state][cuisine]++
	}

	tempMap:=make(map[string] string)
	tempMap["success"]="Order successfully added"
	res=&orderspb.AddOrderResponse{Response: tempMap}
	return res, nil

}


func main()  {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Sorry failed to load server %v: ", err)
	}

	s := grpc.NewServer()
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

		Restaurant_count[restaurant.(string)]++
		Cuisine_count[cuisine.(string)]++
		Orders[id.(string)]++
		statemap, ok := State_cuisine_count[state.(string)]
		if ok {
			statemap[cuisine.(string)]++
		} else {
			State_cuisine_count[state.(string)] = make(map[string]int64)
			State_cuisine_count[state.(string)][cuisine.(string)]++
		}
	}
	fmt.Println(APIUtil.Orders["2999999"])
	orderspb.RegisterOrdersServiceServer(s, &OrdersServiceServer{})

	fmt.Println("Orders Server starting...")
	if s.Serve(lis); err != nil {
		log.Fatalf("failed to Serve %v", err)
	}
}


