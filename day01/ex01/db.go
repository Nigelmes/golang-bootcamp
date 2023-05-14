package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type DBReader interface {
	Read() (Recipes, error)
	Convert(dbreader Recipes) (string, error)
}

type Recipes struct {
	Cake []struct {
		Name        string `xml:"name" json:"name"`
		Time        string `xml:"stovetime" json:"time"`
		Ingredients []struct {
			IngredientName  string `xml:"itemname" json:"ingredient_name"`
			IngredientCount string `xml:"itemcount" json:"ingredient_count"`
			IngredientUnit  string `xml:"itemunit" json:"ingredient_unit,omitempty"`
		} `xml:"ingredients>item" json:"ingredients"`
	} `xml:"cake" json:"cake"`
}

type jsonfile string

func (j jsonfile) Read() (Recipes, error) {
	var recipes Recipes
	text, err := os.ReadFile(string(j))
	if err != nil {
		return Recipes{}, fmt.Errorf("Ошибка: %w", err.Error())
	}
	if err = json.Unmarshal(text, &recipes); err != nil {
		return Recipes{}, fmt.Errorf("Ошибка: %w", err.Error())
	}
	return recipes, nil
}

func (j jsonfile) Convert(recipes Recipes) (string, error) {
	res, err := xml.MarshalIndent(recipes, "", "    ")
	if err != nil {
		return "", err
	}
	return string(res), nil
}

type xmlfile string

func (x xmlfile) Read() (Recipes, error) {
	var recipes Recipes
	text, err := os.ReadFile(string(x))
	if err != nil {
		return Recipes{}, fmt.Errorf("Ошибка: %w", err.Error())
	}
	if err = xml.Unmarshal(text, &recipes); err != nil {
		return Recipes{}, fmt.Errorf("Ошибка: %w", err.Error())
	}
	return recipes, nil
}

func (x xmlfile) Convert(recipes Recipes) (string, error) {
	res, err := json.MarshalIndent(recipes, "", "    ")
	if err != nil {
		return "", err
	}
	return string(res), err
}

func ConvertingAndPrint(dbreader DBReader) {
	recipes, err := dbreader.Read()
	if err != nil {
		fmt.Printf("ошибка чтения: %w\n", err.Error())
		return
	}
	result, err := dbreader.Convert(recipes)
	if err != nil {
		fmt.Printf("ошибка конвертации: %w\n", err.Error())
		return
	}
	fmt.Println(result)
}

func GetStructRecipes(file string) (Recipes, error) {
	if filepath.Ext(file) == ".json" {
		return jsonfile(file).Read()
	}
	return xmlfile(file).Read()
}
