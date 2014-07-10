/*
http://ichart.finance.yahoo.com/table.csv?s=600705
http://table.finance.yahoo.com/table.csv?s=600705.SS
*/
package ARHelper

import (
	"bufio"
	"fmt"
	"os/exec"
	//	"github.com/moovweb/gokogiri"
	// "io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	// "sync"
	"net/url"

	"path/filepath"
	"time"
)

const (
	UA = "Golang Downloader form home.com"
)

var yahooApi = "http://ichart.finance.yahoo.com/table.csv?s="

// func GetCsv(wg *sync.WaitGroup, stock string) (err error) {
func GetCsv(stock string) (err error) {
	stock = stock + ".ss"
	// 固定存储位置
	var file *os.File
	var fileExist bool
	var needSync bool
	needSync = true
	//
	execPath, _ := exec.LookPath(os.Args[0])

	if strings.Contains(execPath, "go-build") {
		//this is under  is dev env go run
		execSharePath = "../../share/"
	} else {
		//this is under production env
		execSharePath = filepath.Dir(execPath) + "./share/"
	}
	filename := execSharePath + stock
	fmt.Println("filename is ", filename)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		path, _ := os.Getwd()
		fmt.Println("no such file or directory: will create it ", path, filename)

		if os.MkdirAll(filepath.Dir(filename), 0755) != nil {
			panic("Unable to create directory for tagfile!")
		}

		file, err = os.Create(filename)
		fileExist = false
	} else {

		file, err = os.OpenFile(filename, os.O_RDWR, 0666)
		fileExist = true
	}
	checkErr(err)
	defer file.Close()

	//TODO 自动更新数据
	t := time.Now()
	// fmt.Printf("hour is %d, week is %d\n", t.Hour(), t.Weekday())

	if fileExist {
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		scanner.Scan()
		if len(scanner.Text()) > 0 {
			// 获取那一行的日期
			lastDate := strings.Split(scanner.Text(), ",")[0]
			fmt.Println("lastDate is ", lastDate)
			//parse 日期
			const shortForm = "2006-01-02"
			lastTime, err := time.Parse(shortForm, lastDate)
			checkErr(err)
			//如果日期相等

			if t.Day() == lastTime.Day() {
				needSync = false

			} else {
				//星期六星期天
				if t.Weekday() == 0 {
					if (t.Day() - 2) == lastTime.Day() {
						needSync = false
					}
				} else if t.Weekday() == 6 {
					if (t.Day() - 1) == lastTime.Day() {
						needSync = false
					}
				}
			}
		}
	}

	checkErr(err)
	fmt.Println(t)
	if needSync {
		var urlStr = yahooApi + stock
		// urlStr := "http://127.0.0.1:4567/hello"
		fmt.Println("get url is ", urlStr)

		var req http.Request
		req.Method = "GET"

		header := http.Header{}
		header.Set("User-Agent", UA)
		req.Header = header
		// req.UserAgent = UA
		req.Close = true
		req.URL, err = url.Parse(urlStr)
		if err != nil {
			fmt.Println("error url parse is ", err)
			panic(err)
		}

		resp, err := http.DefaultClient.Do(&req)

		if err != nil {
			fmt.Println("Default client do error url parse is ", err)
			panic(err)
		}

		defer resp.Body.Close()
		//status 为200的时候存入文件
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			file.Truncate(0)
			written, err := file.Write(body)
			// written, err := io.Copy(file, resp.Body)
			if err != nil {
				panic(err)
			}
			println("written: ", written)
		}

	}
	fmt.Printf(stock+" is finished at The call took %v to run.\n", time.Now().Sub(t))
	// wg.Done()
	return err
}
