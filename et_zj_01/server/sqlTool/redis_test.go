package sqlTool

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"testing"
	"time"
)

//import (
//	"fmt"
//	"github.com/go-redis/redis"
//	"log"
//	"testing"
//	"time"
//)
func TestGetAllDataFromMySQLToRedis(t *testing.T)  {
	now := time.Now()
	getAllDataFromMySQLToRedis()
	fmt.Println("get all data from Mysql to Redis time : ", time.Since(now))
	resRadiu, err := rdb.GeoRadius(ctx,"outlets", 111, -6, &redis.GeoRadiusQuery{
		Radius:      800,
		Unit:        "km",
		WithCoord:   true,
		WithGeoHash: true,
		WithDist:    true,
		Count:       10000,
		Sort:        "ASC",
	}).Result()
	if err != nil {
		log.Fatalf("get outlet location from redis fail : %v", err)
	}

	for _, v := range resRadiu{
		fmt.Print("outletId : ", v.Name, "   Longitude : ", v.Longitude, "   Latitude : ", v.Latitude, "    dist : ", v.Dist, "\n")
	}
	fmt.Println(len(resRadiu))
}

//func TestJob(t *testing.T)  {
//	//将数据进行gob序列化
//	o1 := Outlets{Id: 111}
//	bytes, _ := o1.MarshalBinary()
//	rdb.Set(ctx, "test", bytes, 10*time.Minute)
//
//	fmt.Println(rdb.Get(ctx,"test"))
//	i, _ := rdb.Get(ctx, "test").Bytes()
//	o2_bytes := i
//
//	var o2 Outlets
//
//	o2.UnmarshalBinary(o2_bytes)
//
//	fmt.Println(o2)
//}
//func Test01(t *testing.T)  {
//	initClient()
//}
