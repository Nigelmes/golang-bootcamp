package main

import (
	"bufio"
	"fmt"
	"os"
)

func processFile(filename string) (map[string]struct{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		result[line] = struct{}{}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func printChanges(oldmap, newmap map[string]struct{}) {
	for key := range newmap {
		if _, ok := oldmap[key]; !ok {
			fmt.Printf("ADDED %s\n", key)
		}
	}
	for key := range oldmap {
		if _, ok := newmap[key]; !ok {
			fmt.Printf("REMOVED %s\n", key)
		}
	}
}
