package main

import (
	"github.com/shubham491/order-analysis/pkg/FileUtil"
	"log"
)

func main() {
	csvFilePath := "sample_data_2.csv"
	jsonFilePath := "outputs.json"
	dh := FileUtil.DataHandler{CsvFilePath: csvFilePath, JsonFilePath: jsonFilePath}
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
