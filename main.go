package main
import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
)
type Order struct {
	OrderId int
	Discount float64
	Amount float64
	PaymentMode  string
	Rating int
	Duration int
	Cuisine string
	OrderTime int
	RestId int
	State  string
	CustId string

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

	r.WriteLine("{")
	r.WriteLine("\"orders\": [")
	HandleErr(err, "Error writing to JSON file")
}

func (r *DataHandler) Init() {
	var err error
	r.csvFd, err = os.Open(r.CsvFilePath)
	HandleErr(err, "Error reading the CSV file")
	r.csvReader = csv.NewReader(r.csvFd)

	r.InitJsonWriter()
	firstRow, done := r.ReadLine()
	firstRow, done = r.ReadLine()
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
	r.WriteLine("]\n}")
	r.jsonFd.Close()
	r.csvFd.Close()
}

func (r *DataHandler) CreateOrder(data []string) Order {
	order := Order{
		OrderId:    ParseInt(data[0]),
		Discount:   ParseFloat(data[1]),
		Amount:     ParseFloat(data[2]),
		PaymentMode:data[3],
		Rating:     ParseInt(data[4]),
		Duration:   ParseInt(data[5]),
		Cuisine:    data[6],
		OrderTime:  ParseInt(data[7]),
		RestId:     ParseInt(data[8]),
		State:      data[9],
		CustId:     data[10],

	}
	return order
}

func main() {
	csvFilePath := "sample_data.csv"
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
