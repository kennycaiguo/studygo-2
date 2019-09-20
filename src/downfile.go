package main

import (
	"db"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"property"
	"regexp"
	"sync"
	"time"
)

var (
	cfg              = property.Cfg
	fileStorePath, _ = cfg.GetValue("file", "storepath")
)

/**
生成当天目录 格式 yyyy-MM-dd
*/
func initpath() {
	currentTime := time.Now()
	fileStorePath = fileStorePath + "/" + currentTime.Format("2006-01-02")

	_, err := os.Stat(fileStorePath)
	if err != nil {
		err = os.Mkdir(fileStorePath, os.ModePerm)
	}
}

/**
下载文件
*/

func downfile(dataMap map[string]string, filename string) {
	url := dataMap["url"]
	start := time.Now()
	client := &http.Client{}

	response, err := client.Get(url)

	if err != nil {
		fmt.Sprintf("get fail %s", url)
		fmt.Println(err)
		wg.Done()
		return
	}
	defer response.Body.Close()
	//now := time.Now()

	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		//fmt.Println(url+"返回失败")
		wg.Done()
		return
	}

	err = ioutil.WriteFile(fileStorePath+"/"+filename, data, 0666)
	if err != nil {
		fmt.Println(url + "写入失败")
		wg.Done()
		return
	}
	dataMap["path"] = fileStorePath + "/" + filename

	fmt.Printf("%s下载总时间:%v\n", url, time.Now().Sub(start))
	wg.Done()

}

var wg = new(sync.WaitGroup)

func main() {
	initpath()
	urls := db.QueryUrlstest()
	wg.Add(len(urls))
	for _, urlmap := range urls {
		url := urlmap["url"]
		filename := getfilename(url)
		go downfile(urlmap, filename)
	}

	wg.Wait()
	fmt.Println("process over")

}

/**
下载后的文件以url为文件名存储 并替换掉https 和一些特殊符号
*/
func getfilename(url string) (filename string) {
	re3, _ := regexp.Compile("http[s]?://")
	rep := re3.ReplaceAllString(url, "")
	re3, _ = regexp.Compile("[/|?|\\||\\*]")
	filename = re3.ReplaceAllString(rep, "_")
	return
}
