package property

import (
	"flag"
	"fmt"
	"github.com/Unknwon/goconfig"
	"os"
)

var Cfg *goconfig.ConfigFile

func init() {
	filestorepath := flag.String("inipath", "C:/Users/20160712/Desktop/spider.ini", "ini path")
	flag.Parse()
	config, err := goconfig.LoadConfigFile(*filestorepath) //加载配置文件
	if err != nil {
		fmt.Println("get config file error")
		os.Exit(-1)
	}
	Cfg = config

}
