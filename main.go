package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

type DataHandler struct {
	CsvFilePath  string
	csvFd        *os.File
	JsonFilePath string
	jsonFd       *os.File
	csvReader    *csv.Reader
	//customers	map[int]Customer
	//restaurants map[int]Restaurant
	//cuisineCustomer map[string]map[int]bool
	//orderRestaurant map[int]map[int]bool
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

func (r *DataHandler) InitJsonWriter() {
	var err error
	if Exists(r.JsonFilePath) {
		err = os.Remove(r.JsonFilePath)
		HandleErr(err, "Error removing exiting output file")
	}

	r.jsonFd, err = os.OpenFile(r.JsonFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	HandleErr(err, "Error opening JSON output file")

	r.WriteLine("{\"orders\": [")
	HandleErr(err, "Error writing to JSON file")
}

func (r *DataHandler) Init() {
	var err error
	r.csvFd, err = os.Open(r.CsvFilePath)
	HandleErr(err, "Error reading the CSV file")
	r.csvReader = csv.NewReader(r.csvFd)

	r.InitJsonWriter()
	r.ReadLine() // Header line
	firstRow, done := r.ReadLine()
	if done {
		log.Fatal("Got an empty CSV file")
	}
	order := r.CreateOrder(firstRow)
	orderJson, _ := json.Marshal(order)
	r.jsonFd.WriteString(string(orderJson))
}

func (r *DataHandler) ReadLine() ([]string, bool) {
	data, err := r.csvReader.Read()
	if data == nil {
		return nil, true
	}
	HandleErr(err, "Error reading from CSV file")
	return data, false
}

func (r *DataHandler) WriteLine(txt string) {
	_, err := r.jsonFd.WriteString(txt)
	HandleErr(err, "Error writing to JSON file")
}

func (r *DataHandler) WriteOrder(order Order) {
	orderJson, _ := json.Marshal(order)
	r.jsonFd.WriteString(",")
	r.jsonFd.WriteString(string(orderJson))
}

func (r *DataHandler) Close() {
	r.WriteLine("]}")
	r.jsonFd.Close()
	r.csvFd.Close()
}

func (r *DataHandler) CreateOrder(data []string) Order {
	custId := ParseInt(data[11])
	custName := data[12]
	restId := ParseInt(data[8])
	restName := data[9]
	state := data[10]
	cuisine := data[6]

	order := Order{
		Id:          ParseInt(data[0]),
		Discount:    ParseFloat(data[1]),
		Amount:      ParseFloat(data[2]),
		PaymentMode: data[3],
		Rating:      ParseInt(data[4]),
		Duration:    ParseInt(data[5]),
		Cuisine:     cuisine,
		Time:        ParseInt(data[7]),
		RestId:      restId,
		RestName:    restName,
		State:       state,
		CustId:      custId,
		CustName:    custName,
	}

	//r.customers[custId] = Customer{Id: custId, Name: custName, State: state}
	//r.restaurants[restId] = Restaurant{Id: restId, Name: restName, State: state}

	return order
}

func main() {
	csvFilePath := "sample_data_2.csv"
	jsonFilePath := "outputs.json"
	dh := DataHandler{CsvFilePath: csvFilePath, JsonFilePath: jsonFilePath}
	dh.Init()
	defer dh.Close()
	for {
		data, done := dh.ReadLine()
		if done {
			log.Println("Reached end of file")
			break
		}
		order := dh.CreateOrder(data)
		dh.WriteOrder(order)
	}
}
