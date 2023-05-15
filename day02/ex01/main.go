package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
)

type flags struct {
	w, l, m bool
}

var fl flags

func init() {
	flag.BoolVar(&fl.w, "w", false, "w flag")
	flag.BoolVar(&fl.l, "l", false, "l flag")
	flag.BoolVar(&fl.m, "m", false, "m flag")
}

func main() {
	wg := new(sync.WaitGroup)
	flag.Parse()
	if err := checkFlag(); err != nil {
		panic(err)
	}
	if err := checkValidFile(); err != nil {
		panic(err)
	}
	for _, filename := range flag.Args() {
		wg.Add(1)
		go printCountlines(filename, wg)
	}
	wg.Wait()
}

func checkFlag() error {
	if flag.NFlag() > 1 {
		return fmt.Errorf("ошибка: одновременно можно использовать только один флаг")
	} else if flag.NFlag() == 0 {
		fl.w = true
	}
	return nil
}

func checkValidFile() error {
	for _, filename := range flag.Args() {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return fmt.Errorf("%s - такого файла не существует", filename)
		}
	}
	return nil
}

func printCountlines(filename string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	if fl.w {
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			count++
		}
		fmt.Printf("%d\t%s\n", count, filename)
	} else if fl.l {
		for scanner.Scan() {
			count++
		}
		fmt.Printf("%d\t%s\n", count, filename)
	} else if fl.m {
		scanner.Split(bufio.ScanRunes)
		for scanner.Scan() {
			count++
		}
		fmt.Printf("%d\t%s\n", count, filename)
	}
	if err = scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
}
