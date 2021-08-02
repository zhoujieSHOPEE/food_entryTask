package sqlTool

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
)

var ctx = context.Background()

func GetNearStore(longitude, latitude float64) []Outlets {
	//fmt.Println("into GetNearStore")
	//now := time.Now()
	resRadiu, err := rdb.GeoRadius(ctx,"outlets", longitude, latitude, &redis.GeoRadiusQuery{
		Radius:      5,
		Unit:        "km",
		WithCoord:   true,
		WithGeoHash: true,
		WithDist:    true,
		Count:       1000,
		Sort:        "ASC",
	}).Result()
	//fmt.Println("pos:1")
	//fmt.Println("georaius耗时：", time.Since(now))
	if err != nil {
		log.Fatalf("get outlet location from redis fail : %v", err)
	}


	//这个是直接从redis获取outlets集合的，但是压测很差，我尝试用redis获得Id以后再从数据库获取完整outlets信息。
	// 我已经获得足够的地理位置集合，现在根据地理位置中的id获取完整的outlets
	outletsListFromRedis := make([]Outlets, 0)
	for _, v := range resRadiu {
		bytes, _ := rdb.Get(ctx, v.Name).Bytes()
		o2 := &Outlets{}
		json.Unmarshal(bytes, o2)
		outletsListFromRedis = append(outletsListFromRedis, *o2)
	}
	//fmt.Println("pos:2")


	//outletsListFromMysql := make([]Outlets, 0)
	//for _, v := range resRadiu {
	//	atoi, _ := strconv.Atoi(v.Name)
	//	o, err := FindOutletsById(atoi)
	//	if err != nil {
	//		log.Println("get outlets from findoutletsbyid fail : %v", err)
	//	}
	//	outletsListFromMysql = append(outletsListFromMysql, o)
	//}
	return outletsListFromRedis
}
