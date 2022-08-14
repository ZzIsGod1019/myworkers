package database

import (
	"crypto/tls"
	"log"
	"myworkers/config"

	"github.com/go-redis/redis"
)

// MainRedis 主实例 key为语言(en,es,ko,ru,id)
var MainRedis map[string]*redis.Client

// MainRedisPrefix 主实例 redis key 的前缀
var MainRedisPrefix map[string]string

func RedisInit() {
	locales := config.GetStringSlice("app.locales")

	MainRedis = make(map[string]*redis.Client)
	MainRedisPrefix = make(map[string]string)

	for _, locale := range locales {
		log.Println("redis init locale:" + locale)

		// --- main实例 ---
		mainConfName := "redis_main_" + locale
		mainHost := config.GetStringKey(mainConfName + ".host")
		mainPort := config.GetStringKey(mainConfName + ".port")
		mainDb := config.GetIntKey(mainConfName + ".db")
		mainPoolSize := config.GetIntKey(mainConfName + ".pool_size")
		mainPassword := config.GetStringKey(mainConfName + ".password")
		mainScheme := config.GetStringKey(mainConfName + ".scheme")
		MainRedisPrefix[locale] = config.GetStringKey(mainConfName + ".prefix")

		mainConf := &redis.Options{
			Addr:     mainHost + ":" + mainPort,
			Password: mainPassword, // no password set
			DB:       mainDb,       // use default DB
			PoolSize: mainPoolSize,
		}

		if mainScheme == "tls" {
			mainConf.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
				//Certificates: []tls.Certificate{cert}
			}
		}
		MainRedis[locale] = redis.NewClient(mainConf)
	}
}
