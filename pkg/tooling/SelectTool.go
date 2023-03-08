package tooling

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
)

type ProdStore interface {
	GetCountForSelectTool(ctx context.Context, nominal int) (int, error)
	GetForClient(ctx context.Context, SubcatID, Count int) ([]map[string]interface{}, error)
}
type SubcatStore interface {
	GetValueByID(ctx context.Context, SubcategoryID int) (int, error)
	GetData(ctx context.Context, SubcategoryID, CategoryID int) (string, int, error)
	GetIDByValue(ctx context.Context, Value int, CategoryID int, IsComposite bool) (int, error)
}
type Tool struct {
	ProdStore   ProdStore
	SubcatStore SubcatStore
}

func NewTool(prodStore ProdStore, subcatStore SubcatStore) *Tool {
	return &Tool{ProdStore: prodStore, SubcatStore: subcatStore}
}

type SubcatCounts struct {
	Count    int
	SubcatID int
}

type ElementCount struct {
	Element int
	Count   int
}

func (t *Tool) SelectTool(ctx context.Context, SubcategoryID int, CategoryID int) ([]map[string]interface{}, error) {
	//Getting String with Available Slice
	logrus.Println("Select Tool ЮХУ МЫ ЗДЕСЬ")
	AvValuesString, TargetSum, err := t.SubcatStore.GetData(ctx, SubcategoryID, CategoryID)
	if err != nil {
		logrus.Printf("Step 1: %s", err.Error())
		return nil, err
	}
	logrus.Printf("Step 1: %s, %d", AvValuesString, TargetSum)

	// String to Slice with Available Nominals
	NominalSlice := stringToIntSlice(AvValuesString)
	logrus.Printf("Step 2 (Slice Nominals): %v", NominalSlice)

	// Getting Full Array by selecting counts from DB
	FullNotSortedArray, err := t.GetFullArray(ctx, NominalSlice, CategoryID)
	if err != nil {
		logrus.Printf("Step 3 (Full Array): %v", err)
		return nil, err
	}
	//_____________ СЛомалось здесь

	logrus.Printf("Step 3 (Full Array): %v", FullNotSortedArray)

	NeedingArray, err := minimumSumArray(FullNotSortedArray, TargetSum)
	if err != nil {
		logrus.Printf("Step 4: %s", err.Error())
		return nil, err
	}
	logrus.Println("Needing Array: ", NeedingArray)

	CountsElements := countElements(NeedingArray)
	logrus.Printf("Step 5 (CountElements): %+v", CountsElements)

	prods, err := t.GetCompositeKeys(ctx, CountsElements, CategoryID)
	if err != nil {
		logrus.Printf("Step 6 (Get Products): %s", err.Error())
		return nil, err
	}
	logrus.Printf("Step 6 (Get Products): %v", prods)
	return prods, nil
}
func (t *Tool) SelectToolCheck(ctx context.Context, SubcategoryID int, CategoryID int) (bool, error) {
	//Getting String with Available Slice
	logrus.Println("Select Tool ЮХУ МЫ ЗДЕСЬ")
	AvValuesString, TargetSum, err := t.SubcatStore.GetData(ctx, SubcategoryID, CategoryID)
	if err != nil {
		logrus.Printf("Step 1: %s", err.Error())
		return false, err
	}
	logrus.Printf("Step 1: %s, %d", AvValuesString, TargetSum)

	// String to Slice with Available Nominals
	NominalSlice := stringToIntSlice(AvValuesString)
	logrus.Printf("Step 2 (Slice Nominals): %v", NominalSlice)

	// Getting Full Array by selecting counts from DB
	FullNotSortedArray, err := t.GetFullArray(ctx, NominalSlice, CategoryID)
	if err != nil {
		logrus.Printf("Step 3 (Full Array): %v", err)
		return false, err
	}

	logrus.Printf("Step 3 (Full Array): %v", FullNotSortedArray)

	_, err = minimumSumArray(FullNotSortedArray, TargetSum)
	if err != nil {
		logrus.Printf("Step 4: %s", err.Error())
		return false, err
	}

	return true, nil
}

func (t *Tool) GetCompositeKeys(ctx context.Context, CountElements []ElementCount, CategoryID int) ([]map[string]interface{}, error) {
	var SubCounts []SubcatCounts

	logrus.Printf("Get Composite Keys GOT Array: %v", CountElements)
	//Getting SubCat IDS
	for _, element := range CountElements {
		var SubCount SubcatCounts
		SubcatID, err := t.SubcatStore.GetIDByValue(ctx, element.Element, CategoryID, false)
		if err != nil {
			logrus.Printf("GetCompositeKeys: %v", err)
			return nil, err
		}
		SubCount.Count = element.Count
		SubCount.SubcatID = SubcatID
		SubCounts = append(SubCounts, SubCount)
	}

	logrus.Printf("SubCounts: %v", SubCounts)
	var products []map[string]interface{}

	for _, value := range SubCounts {
		prods, err := t.ProdStore.GetForClient(ctx, value.SubcatID, value.Count)
		if err != nil {
			logrus.Printf("Get For Client From SelectTool: %v", err)
			return nil, err
		}

		if len(prods) == 0 {
			logrus.Printf("Len of Prods: %v", value)
			return nil, err
		}
		logrus.Println("Prods: ", prods)
		for _, prod := range prods {
			ValueOfProd, err := t.SubcatStore.GetValueByID(ctx, prod["subcategory_id"].(int))
			if err != nil {
				return nil, err
			}
			prod["content_key"] = fmt.Sprintf("%d: %s", ValueOfProd, prod["content_key"].(string))
			products = append(products, prod)
		}
	}
	return products, nil
}

func (t *Tool) GetFullArray(ctx context.Context, Nominals []int, CategoryID int) ([]int, error) {
	var FullArr []int
	logrus.Println("GetFullArray (Nominals Got): ", Nominals)
	for _, nominal := range Nominals {
		logrus.Println("Get Full Array Searching: ", nominal)
		SubCatID, err := t.SubcatStore.GetIDByValue(ctx, nominal, CategoryID, false)
		if err != nil {
			logrus.Printf("Get ID By Value Subcat: %v", err)
			return nil, err
		}
		count, err := t.ProdStore.GetCountForSelectTool(ctx, SubCatID)
		if err != nil {
			logrus.Printf("Get Count Products: %v", err)
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
			logrus.Printf("String To Slice: %v", err)
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
		return nil, errors.New("Мы не можем подобрать")
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
