package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Data struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Location struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"location"`
}

func makeData(record []string) (Data, error) {
	if len(record) != 6 {
		return Data{}, fmt.Errorf("invalid person slice: %v", record)
	}
	id := record[0]
	name := record[1]
	address := record[2]
	phone := record[3]
	lon, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		return Data{}, fmt.Errorf("invalid Longitude: %s", record[4])
	}
	lat, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return Data{}, fmt.Errorf("invalid Latitude: %s", record[5])
	}
	return Data{
		ID: id, Name: name, Address: address, Phone: phone,
		Location: struct {
			Lon float64 `json:"lon"`
			Lat float64 `json:"lat"`
		}{Lon: lon, Lat: lat},
	}, nil
}

func parseCsvFile() ([]Data, error) {
	var res []Data
	file, err := os.Open(csvFileName)
	if err != nil {
		return []Data{}, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	read := csv.NewReader(file)
	read.Comma = '\t'
	_, _ = read.Read() // скпипаем первую строку csv файла
	for {
		record, err := read.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []Data{}, err
		}
		data, err := makeData(record)
		if err != nil {
			return []Data{}, err
		}
		res = append(res, data)
	}
	return res, nil
}
