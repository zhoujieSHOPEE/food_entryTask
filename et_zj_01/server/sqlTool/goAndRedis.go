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


func InitClient(ctx context.Context) (*redis.Client, *redis.Client, *redis.Client) {
	rdb1 := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",  // no password set
		DB:       0,   // use default DB
		PoolSize: 64, // 连接池大小
		MinIdleConns: 64,
	})
	_, err := rdb1.Ping(ctx).Result()
	if err != nil{
		log.Fatalf("connect Redis1 fail : %v", err)
	}
	rdb2 := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "",  // no password set
		DB:       0,   // use default DB
		PoolSize: 64, // 连接池大小
		MinIdleConns: 64,
	})
	_, err = rdb2.Ping(ctx).Result()
	if err != nil{
		log.Fatalf("connect Redis2 fail : %v", err)
	}
	rdb3 := redis.NewClient(&redis.Options{
		Addr:     "localhost:6381",
		Password: "",  // no password set
		DB:       0,   // use default DB
		PoolSize: 64, // 连接池大小
		MinIdleConns: 64,
	})
	_, err = rdb3.Ping(ctx).Result()
	if err != nil{
		log.Fatalf("connect Redis3 fail : %v", err)
	}
	return rdb1,rdb2,rdb3
}

var rdb1, rdb2, rdb3 = InitClient(context.Background())

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

		rdb1.Set(ctx, intToString(v.Id), datas, time.Hour*24*7)
		// 接着存储了所有订单的地理位置信息
		rdb1.GeoAdd(ctx, "outlets", &redis.GeoLocation{
			Name: intToString(v.Id),
			Longitude: v.Longitude,
			Latitude: v.Latitude,
		})

		rdb2.Set(ctx, intToString(v.Id), datas, time.Hour*24*7)
		// 接着存储了所有订单的地理位置信息
		rdb2.GeoAdd(ctx, "outlets", &redis.GeoLocation{
			Name: intToString(v.Id),
			Longitude: v.Longitude,
			Latitude: v.Latitude,
		})

		rdb3.Set(ctx, intToString(v.Id), datas, time.Hour*24*7)
		// 接着存储了所有订单的地理位置信息
		rdb3.GeoAdd(ctx, "outlets", &redis.GeoLocation{
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