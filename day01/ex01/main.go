package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

func main() {
	old := flag.String("old", "", "json or xml file")
	new := flag.String("new", "", "json or xml file")
	flag.Parse()
	if err := checkFormatFile(*old, *new); err != nil {
		panic(err)
	}
	oldrecipes, err := GetStructRecipes(*old)
	if err != nil {
		panic(err)
	}
	newrecipes, err := GetStructRecipes(*new)
	if err != nil {
		panic(err)
	}
	StructComparison(oldrecipes, newrecipes)
}

func checkFormatFile(f1, f2 string) error {
	if f1 == "" || f2 == "" {
		return fmt.Errorf("ошибка: отсутствуют файлы или не хватает одного файла")
	}
	formatf1, formatf2 := filepath.Ext(f1), filepath.Ext(f2)
	if formatf1 != ".json" && formatf1 != ".xml" || formatf2 != ".json" && formatf2 != ".xml" {
		fmt.Println(formatf1, formatf2)
		return fmt.Errorf("ошибка: файлы должны быть json или xml")
	}
	return nil
}
