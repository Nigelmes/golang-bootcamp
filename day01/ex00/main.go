package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

func main() {
	filename := flag.String("f", "", "json or xml file")
	flag.Parse()
	formatfile, err := checkFormatFile(*filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if formatfile == "json" {
		ConvertingAndPrint(jsonfile(*filename))
	} else if formatfile == "xml" {
		ConvertingAndPrint(xmlfile(*filename))
	}
}

func checkFormatFile(filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("ошибка: отсутствует файл")
	}
	if filepath.Ext(filename) == ".json" {
		return "json", nil
	}
	if filepath.Ext(filename) == ".xml" {
		return "xml", nil
	}
	return "", fmt.Errorf("ошибка: файл должен быть json или xml")
}
