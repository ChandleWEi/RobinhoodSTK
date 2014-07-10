package ARHelper

import (
	// "encoding/csv"
	"fmt"
	//	"os"
	//	"os/exec"
	//	"path/filepath"
	"runtime"

	//	"strings"
	"sync"
	// "time"
)

var (
	execPath      string
	execSharePath string
	//yahooApi = "http://table.finance.yahoo.com/table.csv?s="
)

func DownloadAllCSV() {
	cpuNum := runtime.NumCPU()
	status := runtime.GOMAXPROCS(cpuNum)
	if status < 1 {
		fmt.Println("CPU Nothing to change")
	} else {
		fmt.Println("CPU Num is ", cpuNum)
	}

	CheckGoEnv()
	// _, filename, _, _ := runtieme.Caller(1)
	// execPath, _ = exec.LookPath(os.Args[0])

	// if strings.Contains(execPath, "go-build") {
	// 	//this is under  is dev env go run
	// 	execSharePath = "../share/"
	// } else {
	// 	//this is under production env
	// 	execSharePath = filepath.Dir(execPath) + "./share/"
	// }
	// fmt.Println(">>>>>>>>" + execSharePath)
	// os.MkdirAll(execSharePath, 0755)

	// fmt.Println(os.Args)

	// fmt.Println("filepath is " + execPath)
	// fmt.Println("start to get http csv")
	stockArray := getSHStock()
	//  创建等待数组
	// wg := sync.WaitGroup{}
	// wg.Add(len(stockArray))
	//	length := len(stockArray)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(stockArray []string) {
		defer wg.Done()
		for _, s := range stockArray {
			// stock := s + ".ss"
			// //	fmt.Print(stock)
			// // go GetCsv(&wg, stock)
			// // GetCsv(&wg, stock)
			// url := yahooApi + stock
			// fmt.Println("url is " + url)
			// dstPath := "../share/" + stock
			// fmt.Println("dstPath is  " + dstPath)
			// cmd := exec.Command("curl", "-o", dstPath, url)
			// fmt.Println("cmd is ", cmd.Args)
			// err := cmd.Run()
			// checkErr(err)

			// fmt.Print(" id is ", i, " \n")
			// time.Sleep(2000 * time.Millisecond)
			//curl -o taglist.zip http://www.vim.org/scripts/download_script.php?src_id=7701
			GetCsv(s)

		}
	}(stockArray)
	wg.Wait()

	fmt.Println("DONE.")
}
