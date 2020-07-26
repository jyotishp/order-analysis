package FileUtil

import (
	"encoding/csv"
	"encoding/json"
	"github.com/jyotishp/order-analysis/pkg/Models"
	"github.com/jyotishp/order-analysis/pkg/ErrorHandlers"
	"log"
	"os"
)

type DataHandler struct {
	CsvFilePath  string
	csvFd        *os.File
	JsonFilePath string
	jsonFd       *os.File
	csvReader    *csv.Reader
}





func (r *DataHandler) InitJsonWriter() {
	var err error
	if ErrorHandlers.Exists(r.JsonFilePath) {
		err = os.Remove(r.JsonFilePath)
		ErrorHandlers.HandleErr(err, "Error removing exiting output file")
	}

	r.jsonFd, err = os.OpenFile(r.JsonFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	ErrorHandlers.HandleErr(err, "Error opening JSON output file")

	r.WriteLine("{\"orders\": [")
	ErrorHandlers.HandleErr(err, "Error writing to JSON file")
}

func (r *DataHandler) Init() {
	var err error
	r.csvFd, err = os.Open(r.CsvFilePath)
	ErrorHandlers.HandleErr(err, "Error reading the CSV file")
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
	ErrorHandlers.HandleErr(err, "Error reading from CSV file")
	return data, false
}

func (r *DataHandler) WriteLine(txt string) {
	_, err := r.jsonFd.WriteString(txt)
	ErrorHandlers.HandleErr(err, "Error writing to JSON file")
}

func (r *DataHandler) WriteOrder(order Models.Order) {
	orderJson, _ := json.Marshal(order)
	r.jsonFd.WriteString(",")
	r.jsonFd.WriteString(string(orderJson))
}

func (r *DataHandler) Close() {
	r.WriteLine("]}")
	r.jsonFd.Close()
	r.csvFd.Close()
}

func (r *DataHandler) CreateOrder(data []string) Models.Order {
	custId := ErrorHandlers.ParseInt(data[11])
	custName := data[12]
	restId := ErrorHandlers.ParseInt(data[8])
	restName := data[9]
	state := data[10]
	cuisine := data[6]

	order := Models.Order{
		Id:          ErrorHandlers.ParseInt(data[0]),
		Discount:    ErrorHandlers.ParseFloat(data[1]),
		Amount:      ErrorHandlers.ParseFloat(data[2]),
		PaymentMode: data[3],
		Rating:      ErrorHandlers.ParseInt(data[4]),
		Duration:    ErrorHandlers.ParseInt(data[5]),
		Cuisine:     cuisine,
		Time:        ErrorHandlers.ParseInt(data[7]),
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
