package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	args := os.Args[1:]
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			cmd := exec.Command(args[0], append(args[1:], line)...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println("ошибка выполнения команды:", err)
			}
		}
	}
	if scanner.Err() != nil {
		fmt.Println("ошибка чтения ввода:", scanner.Err())
	}
}
