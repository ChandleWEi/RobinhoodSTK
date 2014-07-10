package ARHelper

import (
	"encoding/csv"
	"fmt"
	"os"

	// 	"github.com/moovweb/gokogiri"
)

var (
	// 	CNAStock = "http://www.sse.com.cn/assortment/stock/list/name/"
	shStock = "../../share/SHAStock.csv"
)

// read stock csv file and return string array
func getSHStock() []string {
	fmt.Println("--->trace getSHStock ")
	path, _ := os.Getwd()
	fmt.Println("current path:", path)
	f, err := os.Open(shStock)
	checkErr(err)
	csvReader := csv.NewReader(f)
	stockMatrix, err := csvReader.ReadAll()
	checkErr(err)
	var stockArray []string
	for _, entry := range stockMatrix {
		stockArray = append(stockArray, entry[0])
	}

	return stockArray
	// fmt.Println(stockArray)

}
