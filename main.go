package main

/*
import "strings"

import (
	"fmt"
	"log"
	"sort"
	"strconv"
)

func main() {
	sum := 0
	partValues := "1000, 1500, 2000, 2500, 3000, 5000"
	target := 10000

	///
	split := strings.Split(partValues, ", ")
	///
	var NewArr []int

	for _, s := range split {
		log.Println("We in cycle")
		Num, err := strconv.Atoi(s)
		NewArr = append(NewArr, Num)
		if err != nil {
			log.Println(err)
		}
	}

	sort.Ints(NewArr)
	for i, j := 0, len(NewArr)-1; i < j; i, j = i+1, j-1 {
		NewArr[i], NewArr[j] = NewArr[j], NewArr[i]
	}

	fmt.Println(NewArr)

	arr := []int{1500, 1500, 1500, 2500, 1000, 2000, 1500, 1000, 2500, 3000, 5000}

	MapValuesCounters := make(map[int]int)

	for sum <= target {

		for i := 0; i < len(NewArr); i++ {

			if (target-sum)/NewArr[i] <= (SelectCount(arr, NewArr[i])) {
				MapValuesCounters[NewArr[i]] = (target - sum) / NewArr[i]
				sum += (NewArr[i]) * target / (NewArr[i])
			} else if (target-sum)/NewArr[i] > SelectCount(arr, NewArr[i]) {
				sum += NewArr[i] * (SelectCount(arr, NewArr[i]))
				MapValuesCounters[NewArr[i]] = SelectCount(arr, NewArr[i])
			}

			if sum > target {
				i = 0
				sum = 0
			}
			log.Println(MapValuesCounters)
			if sum == target {
				break
			}
		}

	}
}

func SelectCount(arr []int, target int) int {
	counter := 0

	for _, i2 := range arr {
		if i2 == target {
			counter++
		}
	}
	return counter

}
*/
