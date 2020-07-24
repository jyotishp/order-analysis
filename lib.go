package main

import (
	"log"
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
