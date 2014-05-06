package main

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"

	"github.com/coopernurse/gorp"
)

// データソース文字列を変換
func convert_datasource(ds string) (result string) {
	url, _ := url.Parse(ds)
	result = fmt.Sprintf("%s@tcp(%s:3306)%s", url.User.String(), url.Host, url.Path)
	return
}

func initDb() *gorp.DbMap {
	var datasource string
	// for heroku with cleardb
	if os.Getenv("CLEARDB_DATABASE_URL") != "" {
		datasource = convert_datasource(os.Getenv("CLEARDB_DATABASE_URL"))
	} else {
		//		datasource = "root:pass@/database_name?charset=utf8"
		datasource = "root@/martini_guestbook?charset=utf8"
	}
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		panic(err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(Greeting{}, "greetings").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		panic(err)
	}
	return dbmap
}
