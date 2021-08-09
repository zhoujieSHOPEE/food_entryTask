package sqlTool

import (
	"context"
	ent "et_zj_01/server/entity"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"time"
)

var ctx = context.Background()

func GetNearStore(longitude, latitude float64, outletsMap map[int]ent.Outlets, count int) ([]ent.Outlets, error) {
	now := time.Now()
	fmt.Println("count : ", count)
	var rdb *redis.Client
	if count % 3 == 0 {
		rdb = rdb1
	}else if count % 3 == 1 {
		rdb = rdb2
	}else {
		rdb = rdb3
	}
	fmt.Println("=========")
	resRadiu, err := rdb.GeoRadius(ctx,"outlets", longitude, latitude, &redis.GeoRadiusQuery{
		Radius:      5,
		Unit:        "km",
		WithCoord:   true,
		WithGeoHash: true,
		WithDist:    true,
		Count: 500,
		Sort:        "ASC",
	}).Result()
	if err != nil {
		if len(resRadiu) == 0 {
			return nil, nil
		}
		log.Fatalf("get outlet location from redis fail : %v", err)
	}

	fmt.Println("GetNearStore GeoRadius : ", time.Since(now))
	outletsListFromRedis := make([]ent.Outlets, 0)
	for _, v := range resRadiu {
		atoi, _ := strconv.Atoi(v.Name)
		o := outletsMap[atoi]
		o.Dist = v.Dist
		outletsListFromRedis = append(outletsListFromRedis, o)
	}
	//fmt.Println(len(resRadiu))
	//for _, v := range resRadiu {
	//	id, _ := strconv.Atoi(v.Name)
	//	o, _ := FindOutletsById(id)
	//	o.Dist = v.Dist
	//	outletsListFromRedis = append(outletsListFromRedis, o)
	//}

	return outletsListFromRedis, err
}
