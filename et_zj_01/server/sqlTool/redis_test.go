package sqlTool

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestGetAllDataFromMySQLToRedis(t *testing.T)  {
	now := time.Now()
	getAllDataFromMySQLToRedis()
	fmt.Println("get all data from Mysql to Redis time : ", time.Since(now))

}


func TestSlice(t *testing.T)  {
	s := []string{"tom", "jack", "lili"}
	fmt.Println(s)

	datas, _ := json.Marshal(s)
	rdb.Set(ctx, "test", datas, time.Hour*1)

	datas, _ = rdb.Get(ctx, "test").Bytes()
	var o []string
	json.Unmarshal(datas, &o)
	fmt.Println(o)

}

func TestAllOutlets(t *testing.T)  {
	now := time.Now()
	for i := 100000; i < 110000; i++ {
		id, _ := FindOutletsById(i)
		fmt.Println(id)
	}
	fmt.Println("获取一次的时间: ", time.Since(now)/10000)
}