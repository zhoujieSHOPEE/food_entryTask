package sqlTool

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

var ctx = context.Background()

func GetNearStore(longitude, latitude float64) []Outlets {
	fmt.Println("enter getNearStore")
	defer fmt.Println("exit getNearStore")

	resRadiu, err := rdb.GeoRadius(ctx,"outlets", longitude, latitude, &redis.GeoRadiusQuery{
		Radius:      15,
		Unit:        "km",
		WithCoord:   true,
		WithGeoHash: true,
		WithDist:    true,
		Count:       1000,
		Sort:        "ASC",
	}).Result()

	if err != nil {
		log.Fatalf("get outlet location from redis fail : %v", err)
	}


	//这个是直接从redis获取outlets集合的，但是压测很差，我尝试用redis获得Id以后再从数据库获取完整outlets信息。
	// 我已经获得足够的地理位置集合，现在根据地理位置中的id获取完整的outlets
	outletsListFromRedis := make([]Outlets, 0)
	for _, v := range resRadiu {
		bytes, _ := rdb.Get(ctx, v.Name).Bytes()
		//fmt.Println(bytes)
		o2 := &Outlets{}
		json.Unmarshal(bytes, o2)
		o2.Dist = v.Dist
		outletsListFromRedis = append(outletsListFromRedis, *o2)
	}

	/*
	//这是从Redis中获得店铺id以后再从mysql中获得详细信息的方法
	outletsListFromMysql := make([]Outlets, 0)
	for _, v := range resRadiu {
		atoi, _ := strconv.Atoi(v.Name)
		o, err := FindOutletsById(atoi)
		if err != nil {
			log.Println("get outlets from findoutletsbyid fail : %v", err)
		}
		outletsListFromMysql = append(outletsListFromMysql, o)
	}
	 */
	return outletsListFromRedis
}
