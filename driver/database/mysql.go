package database

import (
	"database/sql"
	"log"
	"myworkers/driver/config"

	_ "github.com/go-sql-driver/mysql"
)

// 数据库连接 key为语言(en,es,ko,ru,id)
var MainDB map[string]*sql.DB
var MainDatabaseName map[string]string

func MysqlInit() {
	locales := config.GetStringSlice("app.locales")

	MainDB = make(map[string]*sql.DB)
	MainDatabaseName = make(map[string]string)
	var err error

	for _, locale := range locales {
		log.Println("mysql init locale:" + locale)

		mainConfName := "mysql_main_" + locale
		mainHost := config.GetStringKey(mainConfName + ".host")
		mainPort := config.GetStringKey(mainConfName + ".port")
		mainUsername := config.GetStringKey(mainConfName + ".username")
		mainPassword := config.GetStringKey(mainConfName + ".password")
		MainDatabaseName[locale] = config.GetStringKey(mainConfName + ".database")
		MainDB[locale], err = sql.Open("mysql", mainUsername+":"+mainPassword+"@tcp("+mainHost+":"+mainPort+")/?parseTime=true&charset=utf8mb4")
		if err != nil {
			log.Fatalln(err)
		}

		MainDB[locale].SetMaxIdleConns(3)
		MainDB[locale].SetMaxOpenConns(10)

		if err := MainDB[locale].Ping(); err != nil {
			log.Fatalln(err)
		}
	}
}

func GetMainDb(locale string) *sql.DB {
	return MainDB[locale]
}
