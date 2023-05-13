package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

type flags struct {
	mean, median, mode, sd bool
}

var fl flags

func init() {
	flag.BoolVar(&fl.mean, "mean", false, "mean flag")
	flag.BoolVar(&fl.median, "median", false, "median flag")
	flag.BoolVar(&fl.mode, "mode", false, "mode flag")
	flag.BoolVar(&fl.sd, "sd", false, "sd flag")
}

func main() {
	flag.Parse()
	nums := readInput()
	if len(nums) == 0 {
		fmt.Println("пустой ввод данных")
		return
	}
	printMetrics(nums)
}

func readInput() []int {
	var result []int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("ошибка: неверный ввод, введите число снова")
			continue
		}
		if num < -100000 || num > 100000 {
			fmt.Println("ошибка: число должен быть от -100000 до 100000 включительно")
			continue
		}
		result = append(result, num)
	}
	sort.Ints(result)
	return result
}

func calcMean(nums []int) float64 {
	var res float64
	for _, num := range nums {
		res += float64(num)
	}
	return res / float64(len(nums))
}

func calcMedian(nums []int) float64 {
	if len(nums)%2 == 0 {
		return (float64(nums[len(nums)/2-1]) + float64(nums[len(nums)/2])) / 2
	}
	return float64(nums[len(nums)/2])
}

func calcMode(nums []int) int {
	var mode, max int
	countmap := make(map[int]int)
	for _, num := range nums {
		countmap[num]++
	}
	for key, count := range countmap {
		if count == max {
			mode = findMin(mode, key)
		}
		if count > max {
			max = count
			mode = key
		}
	}
	return mode
}

func findMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func calcSD(nums []int) float64 {
	var sum, mean, sd float64
	for _, num := range nums {
		sum += float64(num)
	}
	mean = sum / float64(len(nums))
	for _, num := range nums {
		sd += math.Pow(float64(num)-mean, 2)
	}
	return math.Sqrt(sd / float64(len(nums)))
}

func printMetrics(nums []int) {
	if !fl.mean && !fl.median && !fl.mode && !fl.sd {
		fl.mean = true
		fl.median = true
		fl.mode = true
		fl.sd = true
	}
	if fl.mean {
		fmt.Printf("Mean: %.2f\n", calcMean(nums))
	}
	if fl.median {
		fmt.Printf("Median: %.2f\n", calcMedian(nums))
	}
	if fl.mode {
		fmt.Printf("Mode: %d\n", calcMode(nums))
	}
	if fl.sd {
		fmt.Printf("SD: %.2f\n", calcSD(nums))
	}
}
