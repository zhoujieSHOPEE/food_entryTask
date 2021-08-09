package funcTool

import (
	ent "et_zj_01/server/entity"
	"fmt"
	"math"
	"strconv"
)

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

func stringToFloat64(s string) float64 {

	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func RemoveRepByMap(slc []ent.Outlets) []ent.Outlets {
	result := make([]ent.Outlets, 0)       //存放返回的不重复切片
	tempMap := map[ent.Outlets]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0 //当e存在于tempMap中时，再次添加是添加不进去的，，因为key不允许重复
		//如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e) //当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

func SortByScore(o []ent.Outlets, left, right int) []ent.Outlets {
	if o == nil {
		fmt.Println("结果为空")
		return nil
	}
	for i := left; i < right; i++ {
		for j := i+1; j < right; j++ {
			itemSoldI := float64(o[i].ItemsSold)
			itemSoldJ := float64(o[j].ItemsSold)
			distI := o[i].Dist
			distJ := o[j].Dist
			if itemSoldJ/(distJ+2) > itemSoldI/(distI+2) {
				o[i], o[j] = o[j], o[i]
			}
		}
	}

	return o
}
