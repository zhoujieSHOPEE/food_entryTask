package sqlTool

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"time"
)


func InitClient(ctx context.Context) (rdb *redis.Client) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",  // no password set
		DB:       0,   // use default DB
		PoolSize: 64, // 连接池大小
		MinIdleConns: 64,
		//DialTimeout: 1*time.Second,
		//ReadTimeout: 800*time.Millisecond,
		//WriteTimeout: 800*time.Millisecond,
		//PoolTimeout: 900*time.Millisecond,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil{
		log.Fatalf("connect Redis fail : %v", err)
	}
	return rdb
}

var rdb = InitClient(context.Background())

func getAllDataFromMySQLToRedis() {

	outletsList, err := FindAllOutlets()
	ctx := context.Background()
	if err != nil {
		log.Fatalf("get all outlets from mysql fail : %v", err)
	}
	for _ , v := range outletsList {
		datas, err := json.Marshal(v)
		if err != nil {
			log.Fatalf("struct to bytes fail : %v", err)
		}
		rdb.Set(ctx, intToString(v.Id), datas, time.Hour*12)
		// 接着存储了所有订单的地理位置信息
		rdb.GeoAdd(ctx, "outlets", &redis.GeoLocation{
			Name: intToString(v.Id),
			Longitude: v.Longitude,
			Latitude: v.Latitude,
		})
	}
}

func intToString(i int) string {
	s := strconv.Itoa(i)
	return s
}