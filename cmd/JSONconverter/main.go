package main

import (
	"encoding/csv"
	"encoding/json"
	"github.com/jyotishp/order-analysis/pkg/WriteUtil"
	"log"
	"os"
)





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
