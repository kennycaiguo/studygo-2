package db

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"os"
	"property"
	"strconv"
)

var (
	host     string
	port     string
	user     string
	password string
	dbname   string
)

func init() {
	cfg := property.Cfg
	host, _ = cfg.GetValue("pgsql", "host")
	port, _ = cfg.GetValue("pgsql", "port")
	user, _ = cfg.GetValue("pgsql", "user")
	password, _ = cfg.GetValue("pgsql", "password")
	dbname, _ = cfg.GetValue("pgsql", "dbname")
}

var db *sql.DB

// PGs数据库信息

func opensql() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	checkErr(err)
}

func QueryUrls() (urls []map[string]string) {
	opensql()
	//查询数据
	rows, err := db.Query("SELECT user_login_id,loginname FROM t_sys_userinfo")
	defer db.Close()
	checkErr(err)

	urls = []map[string]string{}
	for rows.Next() {
		var id string
		var url string
		err = rows.Scan(&id, &url)
		data := make(map[string]string)
		data["id"] = id
		data["url"] = url
		urls = append(urls, data)
		checkErr(err)
	}
	return
}

func SaveUrlInfo(urlinfo []map[string]string) {
	opensql()
	defer db.Close()
	stmt, err := db.Prepare("update userinfo set username=$1 where uid=$2")
	checkErr(err)
	for _, urls := range urlinfo {
		stmt.Exec(urls["url"], urls["path"], urls["id"])
	}
}

func QueryUrlstest() (urls []map[string]string) {
	urls = []map[string]string{}

	allurls := ReadUrl()
	for i := 0; i < len(allurls); i++ {
		data := make(map[string]string)
		data["id"] = strconv.Itoa(i)
		data["url"] = allurls[i]
		urls = append(urls, data)
	}

	return
}

func ReadUrl() (urls []string) {
	urls = make([]string, 0)
	fi, err := os.Open("/root/spideraddr.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		urls = append(urls, string(a))
	}
	return
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
