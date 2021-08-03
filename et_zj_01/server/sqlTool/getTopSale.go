package sqlTool

import (
	"log"
	"sort"
)

func GetTopSale(num int) []Outlets {
	/*
	从全集中选取num个销量最高的产品，首先使用Mysql获取跑通，接着用Redis整个存储
	 */
	outletsSlice, err := FindAllOutlets()
	if err != nil {
		log.Fatalf("get all outlets fail : %v", err)
	}

	sort.Sort(outletsWrapper{outletsSlice: outletsSlice})
	return outletsSlice[0:num]
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