package ARHelper

/* 通过mongo db 查询出所有的数据 计算出macd 然后在重新插入回mongo db
 */
import (
	"fmt"
	// "log"
	"io/ioutil"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	// "path/filepath"
	"encoding/csv"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const shortForm = "2006-01-02"
const maxQueryDay = 300
const (
	originType = 0
	macdType   = 1
)

var (
	share   = "../../share/"
	session *mgo.Session
	db      *mgo.Database
)

func BuildAll() {

	var err error

	// init mongo
	//	session, err := mgo.Dial("localhost")
	if session, err = mgo.Dial("localhost"); err != nil {
		panic(err)
	}
	db = session.DB("robinhood")
	//initEnv(db)
	// story
	/*
		  	1.删除所有数据
			2.所有数据导入mongo
			3.生成所有macd导入mongo
	*/
	//read mongo
	// collString := "ss600705"
	buildAllData(db, maxQueryDay)

}

// 通过mongo db 查询出所有的数据 计算出macd 然后在重新插入回mongo db
func buildAllData(db *mgo.Database, day int) {
	collstrings, err := db.CollectionNames()
	fmt.Println("collections count is ", len(collstrings))
	if err != nil {
		fmt.Println("db error is ", err)
		os.Exit(1)
	}
	for _, collString := range collstrings {

		// 有很多不是自己生成的collectionnames
		if collString == "system.indexes" {
			continue
		}

		if matched, _ := regexp.MatchString("ss.*", collString); !matched {
			fmt.Println("Warning It is not good col name is ", collString, " will skiped it ")
			continue
		}

		collection := db.C(collString)
		// collCount, _ := collection.Count()
		// fmt.Println("collection name is ", collection.Name, " count is ", collCount)
		start := time.Now().Add(time.Duration(-24*day) * time.Hour)
		// query := map[string]map[string]time.Time{"date": {"$gte": start}}
		// query := map[string]interface{}{"type": originType, "date": {"$gte": start}}
		query := bson.M{"type": originType, "date": bson.M{"$gte": start}}
		//			}
		queryRS := collection.Find(query)
		var mgoData []map[string]interface{}
		queryRS.All(&mgoData)

		mgoMacd := make(map[string]map[string]float64)
		emaFomula(&mgoMacd, &mgoData, 9)
		emaFomula(&mgoMacd, &mgoData, 12)
		emaFomula(&mgoMacd, &mgoData, 26)
		macdFomula(&mgoMacd)
		wDic2mgo(&mgoMacd, collection)
	}

}

// drop database and init mongo copy all csv data to mongodb
func initEnv(db *mgo.Database) {
	err := db.DropDatabase()
	if err != nil {
		fmt.Println("drop database err is ", err)
	}
	WAllCsv2mgo(false)
}

// 一个map 写入date => {ema12:float64...} ema12 ema9 ema26 macd
func wDic2mgo(mgoMacd_p *map[string]map[string]float64, collection *mgo.Collection) {
	for k, v := range *mgoMacd_p {
		// fmt.Println("collection name is ", collection.Name, " key is ", k, ":", v)
		date, err := time.Parse(shortForm, k)
		// mgoEngry = map[string]interface{}{"date": date, "type": macdType, "ema9": v["ema9"], "ema12": v["ema12"], "ema26": v["ema26"], "macd": v["macd"]}
		// mgoSelect = map[string]interface{}{"date": date, "type": macdType}
		mgoEngry := bson.M{"date": date, "type": macdType, "ema9": v["ema9"], "ema12": v["ema12"], "ema26": v["ema26"], "macd": v["macd"]}
		mgoSelect := bson.M{"date": date, "type": macdType}
		changeInfo, err := collection.Upsert(mgoSelect, mgoEngry)
		if err != nil {
			fmt.Println("csv error is ", err)
			fmt.Println("mongodb changeInfo is ", changeInfo)
		}

		// if verbose {

	}

}

// save all share dir file csv to mongo
func WAllCsv2mgo(verbose bool) {
	var file *os.File
	var err error
	var csvRD *csv.Reader
	var data [][]string
	var dataEntry []string
	var collectionStrArray []string
	var mgoEngry map[string]interface{}
	var mgoSelect map[string]interface{}

	files, _ := ioutil.ReadDir(share)
	// 获取所有csv 文件
	f := files[6]
	// for _, f := range files {
	if verbose {
		fmt.Println(f.Name())
	}
	file, err = os.Open(share + f.Name())
	if err != nil {
		fmt.Println("open file error is ", err)
	}
	defer file.Close()
	csvRD = csv.NewReader(file)

	//Date Open High Low Close Volume Adj Close
	//[2013-01-04 10.30 10.46 10.08 10.21 3488600 9.97]
	collectionStrArray = strings.Split(f.Name(), ".")

	// 每个文件一个文档 600705.ss 存为ss600705
	collectionName := collectionStrArray[1] + collectionStrArray[0]
	fmt.Println("collectionName is ", collectionName)
	if db == nil {
		getDB()
	}
	collection := db.C(collectionName)
	data, err = csvRD.ReadAll()
	if err != nil {
		fmt.Println("csv error is ", err)
	}
	for i := 1; i < len(data); i++ {
		dataEntry = data[i]
		if verbose {
			fmt.Println(dataEntry)
		}
		date, open, high, low, closeM, volume, adj_close := entryParse(dataEntry)
		mgoEngry = map[string]interface{}{"date": date, "type": originType, "open": open, "high": high, "low": low, "close": closeM, "volume": volume, "adj_close": adj_close}
		mgoSelect = map[string]interface{}{"date": date, "type": originType}

		changeInfo, err := collection.Upsert(mgoSelect, mgoEngry)
		if err != nil {
			fmt.Println("csv error is ", err)
		}
		if verbose {
			fmt.Println("mongodb changeInfo is ", changeInfo, " updated ", changeInfo.Updated, " removed ", changeInfo.Removed)
		}
		// 有更新了就不再插入
		if changeInfo.Updated == 1 {
			break
		}
	}

	// }

}

// csv parse type to save to mongo
func entryParse(s []string) (date time.Time, open float64, high float64, low float64, closeM float64, volume float64, adj_close float64) {
	var err error
	// var date time
	// var open, high, low, closeM, volume, adj_close float64

	date, err = time.Parse(shortForm, s[0])
	open, err = strconv.ParseFloat(s[1], 64)
	high, err = strconv.ParseFloat(s[2], 64)
	low, err = strconv.ParseFloat(s[3], 64)
	closeM, err = strconv.ParseFloat(s[4], 64)
	volume, err = strconv.ParseFloat(s[5], 64)
	adj_close, err = strconv.ParseFloat(s[6], 64)
	if err != nil {
		println("entryParse entry error is ", err)
	}

	return date, open, high, low, closeM, volume, adj_close
}

/** =============math==============
*
*
**/
func Average(xs []float64) float64 {
	total := float64(0)
	for _, x := range xs {
		total += x
	}
	return total / float64(len(xs))
}

/** =============macdFomula==============
	// macdFomula(&mgoMacd, &mgoData, 26)
**/
func macdFomula(mgoMacd_p *map[string]map[string]float64) {
	if (*mgoMacd_p) != nil {
		for k, v := range *mgoMacd_p {
			//	fmt.Println("v12 ", v["ema12"], " v26 is ", v["ema26"])
			if ema12, cEma12 := v["ema12"]; cEma12 {
				if ema26, cEma26 := v["ema26"]; cEma26 {
					//	fmt.Println("---------------------------")
					(*mgoMacd_p)[k]["macd"] = ema12 - ema26
				}
			}

		}

	}
}

/** =============emaFomula==============
	// emaFomula(&mgoMacd, &mgoData, 9)
	// emaFomula(&mgoMacd, &mgoData, 12)
	// emaFomula(&mgoMacd, &mgoData, 26)
    // newv 是当天的close 价格oldv 是前天的ema2 第一个是fomulaNum 日期的平均值
	// ema := newv*(2/(fomulaNum+1)) + oldv*(1-(2/(fomulaNum+1)))
**/
func emaFomula(mgoMacd_p *map[string]map[string]float64, mgoData_p *[]map[string]interface{}, fomulaN int) (mgoMacd_r_p *map[string]map[string]float64) {
	// EMA12天平均值
	emaNum := "ema" + strconv.Itoa(fomulaN)
	fomulaNum := float64(fomulaN)
	xs := []float64{}
	last := len(*mgoData_p) - 1
	position := 0
	for i := 0; i < fomulaN; i++ {
		position = last - i
		xs = append(xs, (*mgoData_p)[position]["close"].(float64))

	}
	oldv := Average(xs)
	// fmt.Println("emaFomula#oldv is ", oldv, " fomulaN ", fomulaN)
	position = position - 1
	newv := (*mgoData_p)[position]["close"].(float64)
	ema := newv*(2/(fomulaNum+1)) + oldv*(1-(2/(fomulaNum+1)))
	mgoMacdKey := (*mgoData_p)[position]["date"].(time.Time).Format(shortForm)
	// map[string]float64{emaNum: ema}
	if (*mgoMacd_p)[mgoMacdKey] == nil {
		(*mgoMacd_p)[mgoMacdKey] = make(map[string]float64)
	}
	// fmt.Println("mgoMacdKey is ", mgoMacdKey, "emaNum is ", emaNum, " ema is ", ema)
	(*mgoMacd_p)[mgoMacdKey][emaNum] = ema

	position = position - 1
	for ; position >= 0; position-- {
		// fmt.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")
		// fmt.Println("position is ", position)
		oldv = newv
		newv = (*mgoData_p)[position]["close"].(float64)
		ema = newv*(2/(fomulaNum+1)) + oldv*(1-(2/(fomulaNum+1)))
		mgoMacdKey = (*mgoData_p)[position]["date"].(time.Time).Format(shortForm)
		if _, ok := (*mgoMacd_p)[mgoMacdKey]; !ok {
			// fmt.Println("create position is ", position, " mgoMacdKey is ", mgoMacdKey)
			(*mgoMacd_p)[mgoMacdKey] = make(map[string]float64)
		}
		(*mgoMacd_p)[mgoMacdKey][emaNum] = ema
		// (*mgoMacd_p)[mgoMacdKey] = map[string]float64{emaNum: ema}
		//	fmt.Println("mgoData is ", "mgoMacdKey is ", mgoMacdKey, " position is ", position, emaNum+" is ", ema, " oldv is ", oldv, " newv is ", newv)

	}
	return mgoMacd_p

}
func getDB() {
	var err error
	if session, err = mgo.Dial("localhost"); err != nil {
		panic(err)
	} else {
		db = session.DB("robinhood")
	}
}
