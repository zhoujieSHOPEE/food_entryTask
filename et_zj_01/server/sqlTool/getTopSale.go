package sqlTool

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"strconv"
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
		fmt.Println(err)
		idSlice, err = FindOrderedOutletsWithLimit(100000)
		if err != nil {
			log.Fatalf("get all outlets fail : %v", err)
		}
		datas, err := json.Marshal(idSlice)
		if err != nil {
			log.Fatalf("json.marshal fail : %v", err)
		}
		rdb.Set(ctx, "allOutletsId", datas, time.Hour*24*7)
	}else if err != nil{
		log.Fatalf("redis get allOutlets fail")
	}else {
		now := time.Now()
		fmt.Println("键存在")
		var i []int
		json.Unmarshal(idSliceDatas, &i)
		idSlice = i
		fmt.Println("反序列化 耗时 : ", time.Since(now))
	}

	now = time.Now()
	var outletsSlice = make([]Outlets, 0)
	for _, v := range idSlice[0:num] {
		bytes, _ := rdb.Get(ctx, strconv.Itoa(v)).Bytes()
		o2 := &Outlets{}
		json.Unmarshal(bytes, o2)
		outletsSlice = append(outletsSlice, *o2)
	}
	fmt.Println("根据id反序列化的时间 : ", time.Since(now))

	return outletsSlice[0:num], err
}



type outletsWrapper struct {
	outletsSlice []Outlets
}

func (ow outletsWrapper) Len() int {
	return len(ow.outletsSlice)
}

func (ow outletsWrapper) Swap(i, j int) {
	ow.outletsSlice[i], ow.outletsSlice[j] = ow.outletsSlice[j], ow.outletsSlice[i]
}

func (ow outletsWrapper) Less(i, j int) bool{

	return ow.outletsSlice[i].ItemsSold > ow.outletsSlice[j].ItemsSold
}

