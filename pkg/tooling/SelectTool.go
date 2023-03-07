package tooling

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
)

type ProdStore interface {
	GetCount(ctx context.Context, nominal int) (int, error)
	GetForClient(ctx context.Context, SubcatID, Count int) ([]map[string]interface{}, error)
}
type SubcatStore interface {
	GetData(ctx context.Context, SubcategoryID int) (string, int, error)
	GetIDByValue(ctx context.Context, Value int) (int, error)
}
type Tool struct {
	ProdStore   ProdStore
	SubcatStore SubcatStore
}

type SubcatCounts struct {
	Count    int
	SubcatID int
}

type ElementCount struct {
	Element int
	Count   int
}

func (t *Tool) SelectTool(ctx context.Context, SubcategoryID int) ([]map[string]interface{}, error) {
	//Getting String with Available Slice
	AvValuesString, TargetSum, err := t.SubcatStore.GetData(ctx, SubcategoryID)
	if err != nil {
		logrus.Printf("Step 1: %s", err.Error())
		return nil, err
	}
	logrus.Printf("Step 1: %s, %d", AvValuesString, TargetSum)

	// String to Slice with Available Nominals
	NominalSlice := stringToIntSlice(AvValuesString)
	logrus.Printf("Step 2 (Slice Nominals): %v", NominalSlice)

	// Getting Full Array by selecting counts from DB
	FullNotSortedArray, err := t.GetFullArray(ctx, NominalSlice)
	if err != nil {
		logrus.Printf("Step 3 (Full Array): %v", err)
		return nil, err
	}
	logrus.Printf("Step 3 (Full Array): %v", FullNotSortedArray)

	NeedingArray, err := minimumSumArray(FullNotSortedArray, TargetSum)
	if err != nil {
		logrus.Printf("Step 4: %s", err.Error())
		return nil, err
	}

	CountsElements := countElements(NeedingArray)
	logrus.Printf("Step 5 (CountElements): %v", CountsElements)

	prods, err := t.GetCompositeKeys(ctx, CountsElements)
	if err != nil {
		logrus.Printf("Step 6 (Get Products): %s", err.Error())
		return nil, err
	}
	logrus.Printf("Step 6 (Get Products): %v", prods)
	return prods, nil
}

func (t *Tool) GetCompositeKeys(ctx context.Context, Array []ElementCount) ([]map[string]interface{}, error) {
	var SubCounts []SubcatCounts

	//Getting SubCat IDS
	for _, element := range Array {
		var SubCount SubcatCounts
		SubcatID, err := t.SubcatStore.GetIDByValue(ctx, element.Element)
		if err != nil {
			return nil, err
		}
		SubCount.Count = element.Count
		SubCount.SubcatID = SubcatID
	}

	var products []map[string]interface{}
	for _, value := range SubCounts {
		prods, err := t.ProdStore.GetForClient(ctx, value.SubcatID, value.Count)
		if err != nil {
			logrus.Println(err)
			return nil, err
		}
		for _, prod := range prods {
			products = append(products, prod)
		}
	}
	return products, nil
}

func (t *Tool) GetFullArray(ctx context.Context, Nominals []int) ([]int, error) {
	var FullArr []int
	for _, nominal := range Nominals {

		SubCatID, err := t.SubcatStore.GetIDByValue(ctx, nominal)

		if err != nil {
			return nil, err
		}

		count, err := t.ProdStore.GetCount(ctx, SubCatID)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			FullArr = append(FullArr, nominal)
		}
	}
	return FullArr, nil
}
func stringToIntSlice(str string) []int {
	var arr []int
	for _, s := range strings.Split(str, ", ") {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}
		arr = append(arr, i)
	}
	return arr
}
func minimumSumArray(arr []int, sum int) ([]int, error) { //Непосредственно алгос
	result := make([]int, 0)
	currentSum := 0
	sortedArr := make([]int, len(arr))
	copy(sortedArr, arr)
	sort.Ints(sortedArr)

	for i := len(sortedArr) - 1; i >= 0; i-- {
		if currentSum+sortedArr[i] <= sum {
			result = append(result, sortedArr[i])
			currentSum += sortedArr[i]
		}
		if currentSum == sum {
			break
		}
	}

	if currentSum != sum {
		return nil, errors.New("NewError")
	}
	return result, nil
}

func countElements(arr []int) []ElementCount { //получаем массив возвращаем структуру с элементами и их колличетсвами
	counts := make(map[int]int)

	for _, element := range arr {
		counts[element]++
	}

	result := make([]ElementCount, 0, len(counts))

	for element, count := range counts {
		result = append(result, ElementCount{element, count})
	}

	return result
}
