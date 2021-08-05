package sqlTool

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var ctx = context.Background()

func GetNearStore(longitude, latitude float64) ([]Outlets, error) {
	//fmt.Println("enter getNearStore")
	//defer fmt.Println("exit getNearStore")
	now := time.Now()
	resRadiu, err := rdb.GeoRadius(ctx,"outlets", longitude, latitude, &redis.GeoRadiusQuery{
		Radius:      5,
		Unit:        "km",
		WithCoord:   true,
		WithGeoHash: true,
		WithDist:    true,
		Count: 1000,
		Sort:        "ASC",
	}).Result()
	if err != nil {
		log.Fatalf("get outlet location from redis fail : %v", err)
	}
	fmt.Println("GeoRadius time : ", time.Since(now))
	now = time.Now()
	outletsListFromRedis := make([]Outlets, 0)
	for _, v := range resRadiu {
		bytes, err := rdb.Get(ctx, v.Name).Bytes()
		if err != nil {
			log.Fatalf("get v from redis fail : %v", err)
		}
		o2 := &Outlets{}
		json.Unmarshal(bytes, o2)
		o2.Dist = v.Dist
		outletsListFromRedis = append(outletsListFromRedis, *o2)
	}
	//fmt.Println(len(resRadiu))
	//for _, v := range resRadiu {
	//	id, _ := strconv.Atoi(v.Name)
	//	o, _ := FindOutletsById(id)
	//	o.Dist = v.Dist
	//	outletsListFromRedis = append(outletsListFromRedis, o)
	//}

	fmt.Println("get from mysql by id time : ", time.Since(now))
	return outletsListFromRedis, err
}
