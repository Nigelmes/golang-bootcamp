package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

func main() {
	old := flag.String("old", "", "txt file")
	new := flag.String("new", "", "txt file")
	flag.Parse()
	if err := checkFormatFile(*old, *new); err != nil {
		panic(err)
	}
	mapold, err := processFile(*old)
	if err != nil {
		panic(err)
	}
	mapnew, err := processFile(*new)
	if err != nil {
		panic(err)
	}
	printChanges(mapold, mapnew)
}

func checkFormatFile(f1, f2 string) error {
	if f1 == "" || f2 == "" {
		return fmt.Errorf("ошибка: отсутствуют файлы или не хватает одного файла")
	}
	if filepath.Ext(f1) != ".txt" || filepath.Ext(f2) != ".txt" {
		return fmt.Errorf("ошибка: файлы должны быть формата txt")
	}
	return nil
}
