package sqlTool

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestGetRedis(t *testing.T)  {
	id := "100000"
	now := time.Now()
	for i := 0 ; i < 100000 ; i++ {
		rdb.Get(ctx, id)
	}
	fmt.Println("the method 1 : ", time.Since(now))
	now = time.Now()
	idList := make([]string, 0)
	for i := 0 ; i < 100000 ; i++ {
		idList = append(idList, id)
	}
	rdb.MGet(ctx, idList...)
	fmt.Println("the method 1 : ", time.Since(now))

}

func TestJson1(t *testing.T)  {
	o , _ := FindOutletsById(100000)
	bytes , _ := json.Marshal(o)
	fmt.Println(bytes)
	rdb.Set(ctx, "testjson", bytes, 0)


	bytes1, _ := rdb.Get(ctx, "testjson").Bytes()
	fmt.Println(bytes1)

	res, _ := rdb.MGet(ctx, "testjson").Result()

	//s := string(res[0])
	fmt.Printf("%T", res[0])
}