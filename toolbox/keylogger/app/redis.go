package app

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

//Prefix redis prefix
const Prefix = "keyboard:"

//LastInputEvent last use event
const LastInputEvent = Prefix + "last-event"

//TotalCount total count ZSET
const TotalCount = Prefix + "total"

//KeyMap cache key code map HASH
const KeyMap = Prefix + "key-map"

var connection *redis.Client

//GetRankKey by time
func GetRankKey(time time.Time) string {
	return GetRankKeyByString(time.Format("2006:01:02"))
}

//GetRankKeyByString
func GetRankKeyByString(time string) string {
	return Prefix + time + ":rank"
}

//GetDetailKey by time
func GetDetailKey(time time.Time) string {
	return GetDetailKeyByString(time.Format("2006:01:02"))
}

func GetDetailKeyByString(time string) string {
	return Prefix + time + ":detail"
}

func GetConnection() *redis.Client {
	return connection
}

func InitConnection(option redis.Options) {
	connection = redis.NewClient(&option)
	result, err := connection.Ping().Result()
	if err != nil {
		fmt.Println(result, err)
		os.Exit(1)
	}
}

func CloseConnection() {
	if connection == nil {
		return
	}
	err := connection.Close()
	if err != nil {
		fmt.Println("close redis connection error: ", err)
	}
}