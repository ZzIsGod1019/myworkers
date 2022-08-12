package config

import (
	"log"

	"github.com/spf13/viper"
)

var conf *viper.Viper = nil

var JwtSecrets map[string]string

var Timezones map[string]string

var SplitDbCount int
var SplitTableCount int

// Init 配置加载
func Init() {
	// 测试时，临时写死配置文件
	exPath := "/mnt/d/workspace/myworkers"
	// ex, _ := os.Executable()
	// exPath := filepath.Dir(ex)

	conf = viper.New()
	conf.AddConfigPath(exPath + "/")
	conf.SetConfigName("app")
	conf.SetConfigType("toml")
	log.Println("Loading config file:" + exPath + "/app.toml")

	if err := conf.ReadInConfig(); err != nil {
		log.Println("error when reading config", err)
		panic(err)
	}

	locales := GetStringSlice("app.locales")

	JwtSecrets = make(map[string]string)
	for _, locale := range locales {
		JwtSecrets[locale] = GetStringKey("app.jwt_secret_" + locale)
	}

	SplitDbCount = GetIntKey("app.mysql_split_db_count")
	SplitTableCount = GetIntKey("app.mysql_split_table_count")

	Timezones = make(map[string]string)
	Timezones["en"] = "America/New_York"
	Timezones["es"] = "America/Mexico_City"
	Timezones["ru"] = "Europe/Moscow"
	Timezones["vi"] = "Asia/Ho_Chi_Minh"
	Timezones["ko"] = "Asia/Seoul"
	Timezones["id"] = "Asia/Jakarta"
}

// GetStringKey 获取某个配置项 外部使用: config.GetKey("app.env")
func GetStringKey(configKey string) string {
	return conf.GetString(configKey)
}

func GetIntKey(configKey string) int {
	return conf.GetInt(configKey)
}

func GetStringSlice(configKey string) []string {
	return conf.GetStringSlice(configKey)
}
