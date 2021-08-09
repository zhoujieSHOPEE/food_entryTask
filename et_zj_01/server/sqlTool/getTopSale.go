package sqlTool

import (
	ent "et_zj_01/server/entity"
	ft "et_zj_01/server/funcTool"
)

func GetTopSale(num int, longitude, latitude float64, outletsMap map[int]ent.Outlets, idSlice []int) ([]ent.Outlets, error) {

	var outletsSlice = make([]ent.Outlets, 0)
	for _, v := range idSlice[0:500] {
		o := outletsMap[v]
		o.Dist = ft.GeoDistance(o.Longitude, o.Latitude, longitude, latitude, "K")
		outletsSlice = append(outletsSlice, o)
	}
	return outletsSlice, nil
}


