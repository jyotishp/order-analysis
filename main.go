package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type Order struct {
	Id           int
	CustomerId   int
	RestaurantId int
	Amount       float64
	Status       string
	DEId         int
	Cart         string
	PaymentMode  string
}

func ParseInt(txt string) int {
	n, _ := strconv.Atoi(txt)
	return n
}

func ParseFloat(txt string) float64 {
	n, _ := strconv.ParseFloat(txt, 64)
	return n
}

type DataHandler struct {
	CsvFilePath  string
	csvFd        *os.File
	JsonFilePath string
	jsonFd       *os.File
	csvReader    *csv.Reader
}

func HandleErr(err error, txt string) {
	if err != nil {
		log.Fatal(txt, err.Error())
	}
}

func (r *DataHandler) InitJsonWriter() {
	var err error

	err = os.Remove(r.JsonFilePath)
	HandleErr(err, "Error removing exiting output file")

	r.jsonFd, err = os.OpenFile(r.JsonFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	HandleErr(err, "Error opening JSON output file")

	r.WriteLine("[")
	HandleErr(err, "Error writing to JSON file")
}

func (r *DataHandler) Init() {
	var err error
	r.csvFd, err = os.Open(r.CsvFilePath)
	HandleErr(err, "Error reading the CSV file")
	r.csvReader = csv.NewReader(r.csvFd)

	r.InitJsonWriter()
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
	r.WriteLine("]")
	r.jsonFd.Close()
	r.csvFd.Close()
}

func (r *DataHandler) CreateOrder(data []string) Order {
	order := Order{
		Id:           ParseInt(data[0]),
		CustomerId:   ParseInt(data[1]),
		RestaurantId: ParseInt(data[2]),
		Amount:       ParseFloat(data[3]),
		Status:       data[4],
		DEId:         ParseInt(data[5]),
		Cart:         data[6],
		PaymentMode:  data[7],
	}
	return order
}

func main() {
	csvFilePath := "MOCK_DATA.csv"
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
