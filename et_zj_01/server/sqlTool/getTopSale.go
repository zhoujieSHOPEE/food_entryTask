package sqlTool

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"math"
	"time"
)

func GetTopSale(num int, longitude, latitude float64) ([]Outlets, error) {
	/*
	从全集中选取num个销量最高的产品，首先使用Mysql获取跑通，接着用Redis整个存储
	 */
	idSlice := make([]int, 0)
	now := time.Now()
	idSliceDatas, err := rdb.Get(ctx, "allOutletsId").Bytes()
	fmt.Println("get redis 耗时 : ", time.Since(now))
	if err == redis.Nil {
		fmt.Println("键不存在")
		idSlice, err = FindOrderedOutletsWithLimit(1000)
		if err != nil {
			log.Fatalf("get all outlets fail : %v", err)
		}
		datas, err := json.Marshal(idSlice)
		if err != nil {
			log.Fatalf("json.marshal fail : %v", err)
		}
		rdb.Set(ctx, "allOutletsId", datas, time.Hour*24*2)
	}else if err != nil{
		log.Fatalf("redis get allOutlets fail")
	}else {
		fmt.Println("键存在")
		json.Unmarshal(idSliceDatas, &idSlice)
	}
	var outletsSlice = make([]Outlets, 0)
	for _, v := range idSlice[0:num] {
		o, err:= FindOutletsById(v)
		o.Dist = GeoDistance(o.Longitude, o.Latitude, longitude, latitude, "K")
		if err != nil {
			log.Fatalf("getTopSale FindOutletsById fail : %v", err)
		}
		outletsSlice = append(outletsSlice, o)
	}
	return outletsSlice, err
}

func GeoDistance(lng1 float64, lat1 float64, lng2 float64, lat2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

