package ARHelper

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	// "github.com/moovweb/gokogiri"
	// "net/http"
	"encoding/csv"
	"encoding/json"
	// "io/ioutil"
)

var (
	// CNAStock = "http://www.sse.com.cn/assortment/stock/list/name/"
	stockJsonFile = "../../share/AStock.json"
)

// func getSHStock() {
func DownloadShStockCsv() {
	fmt.Println("============>DownloadShStockCsv start")
	file, err := os.Open(stockJsonFile)
	defer file.Close()
	if err != nil {
		fmt.Println("open file error is ", err)
	}
	scanner := bufio.NewScanner(file)
	// scanner.Scan()
	// fmt.Println(scanner.Text()) // Println will add back the final '\n'

	var dat map[string]interface{}
	var result []interface{}
	for scanner.Scan() {
		// 	// result = append(result, dat["result"])
		if err := json.Unmarshal(scanner.Bytes(), &dat); err != nil {
			panic(err)
		}

		switch reflect.TypeOf(dat["result"]).Kind() {
		case reflect.Slice:
			// fmt.Println(dat["result"])
			v := reflect.ValueOf(dat["result"])

			i := v.Interface()
			s := i.([]interface{})
			for _, value := range s {
				result = append(result, value)
			}

		}
	}
	cfile, err := os.Create("../../share/SHAStock.csv")
	defer cfile.Close()
	var csvArray [][]string
	// fmt.Println(result)
	for _, e := range result {
		entry := e.(map[string]interface{})

		enArray := []string{entry["PRODUCTID"].(string), entry["PRODUCTNAME"].(string), entry["FULLNAME"].(string)}
		csvArray = append(csvArray, enArray)
	}
	csvFile := csv.NewWriter(cfile)
	err = csvFile.WriteAll(csvArray)
	if err != nil {
		fmt.Println("create sh stock csv file error is ", err)
	} else {
		fmt.Println("============>DownloadShStockCsv successed")
	}

}
