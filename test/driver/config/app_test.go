package test

import (
	"myworkers/driver/config"
	"strconv"
	"testing"
)

func TestConfig(t *testing.T) {
	config.Init()
	env := config.GetStringKey("app.env")
	if env == "test" {
		t.Log("config.GetStringKey正常")
	} else {
		t.Error("config.GetStringKey失败:" + env)
	}
	db := config.GetIntKey("redis_main_en.db")
	if db == 0 {
		t.Log("config.GetIntKey正常")
	} else {
		t.Error("config.GetIntKey失败:" + strconv.Itoa(db))
	}
	locales := config.GetStringSlice("app.locales")
	if locales[0] == "en" {
		t.Log("config.GetStringSlice正常")
	} else {
		t.Error("config.GetStringSlice失败:" + locales[0])
	}
	err_key := config.GetStringKey("app.env1")
	if err_key == "" {
		t.Log("config.GetStringSlice使用不存在的key，返回空值")
	} else {
		t.Error("config.GetStringSlice失败:" + err_key)
	}
}
