package main

import (
	"context"
	pb "et_zj_01/proto"
	st "et_zj_01/server/sqlTool"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sort"
	"strconv"
)

type server struct {
	//pb.UnimplementedGetBestStoresServiceServer
}

const (
	port = ":50051"
)

func (*server) GetBestStoresList (ctx context.Context, req *pb.OutletRequest) (*pb.OutletResponse, error) {
	/*
		仅根据店铺距离用户的距离和店铺销量对店铺打分推荐
	*/
	defer fmt.Println("service finish")
	fmt.Println("server get the service, start the service...")
	pos := req.Pos

	outletsSlice := st.GetNearStore(pos.Longitude, pos.Latitude)
	//outletsSlice := st.GetTopSale(1000)
	//fmt.Println("getNearStore耗时:", time.Since(now))
	//outletsSlice, _ := st.FindOutletsByCityId(cityId)
	//outletsSlice, _ := st.FindAllOutlets()
	//距离的计算多余，之前其实已经有了，下一步改进
	var retMessageList []*pb.RetMessage
	for _,v := range outletsSlice{
		itemsSold := strconv.Itoa(v.ItemsSold)
		retMessage := &pb.RetMessage{
			Name: v.Name,
			Address: v.Address,
			LogoURL: v.LogoURL,
			ItemsSold: itemsSold,
			Distance: fmt.Sprintf("%fkm", v.Dist),
		}
		retMessageList = append(retMessageList, retMessage)
	}

	sort.Sort(RetMessageWrapper{retMessageList})
	//fmt.Println("排序耗时：", time.Since(now))
	res := &pb.OutletResponse{Code: 0,List: retMessageList,ListNum: req.ListNum}
	return res, nil
}

func main(){
	//now := time.Now()
	lis, err := net.Listen("tcp", port)
	if err != nil{
		log.Fatalf("failed to listen : %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterGetBestStoresServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	//fmt.Println("服务器端准备耗时:", time.Since(now))
}

//func GeoDistance(lng1 float64, lat1 float64, lng2 float64, lat2 float64, unit ...string) float64 {
//	const PI float64 = 3.141592653589793
//	radlat1 := float64(PI * lat1 / 180)
//	radlat2 := float64(PI * lat2 / 180)
//	theta := float64(lng1 - lng2)
//	radtheta := float64(PI * theta / 180)
//	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
//	if dist > 1 {
//		dist = 1
//	}
//	dist = math.Acos(dist)
//	dist = dist * 180 / PI
//	dist = dist * 60 * 1.1515
//	if len(unit) > 0 {
//		if unit[0] == "K" {
//			dist = dist * 1.609344
//		} else if unit[0] == "N" {
//			dist = dist * 0.8684
//		}
//	}
//	return dist
//}

type RetMessageWrapper struct {
	retMessages []*pb.RetMessage
}

func (sm RetMessageWrapper) Len() int {
	return len(sm.retMessages)
}

func (sm RetMessageWrapper) Swap(i, j int) {
	sm.retMessages[i], sm.retMessages[j] = sm.retMessages[j], sm.retMessages[i]
}

func (sm RetMessageWrapper) Less(i, j int) bool{
	dis1, _ := strconv.ParseFloat(sm.retMessages[i].Distance, 64)
	dis2, _ := strconv.ParseFloat(sm.retMessages[j].Distance, 64)

	itemsSold1, _ := strconv.ParseFloat(sm.retMessages[i].ItemsSold, 64)
	itemsSold2, _ := strconv.ParseFloat(sm.retMessages[j].ItemsSold, 64)
	dis1, dis2 = Normalization(dis1, dis2)
	itemsSold1, itemsSold2 = Normalization(itemsSold1, itemsSold2)
	score1 := 0.4*itemsSold1 + 0.6*dis2
	score2 := 0.4*itemsSold2 + 0.6*dis1
	return score1 > score2
}

func Normalization(i, j float64) (float64, float64) {
	/*
	对两个数归一化
	 */
	if i == 0 || j == 0{
		return i,j
	}
	if i > j{
		j = j/i
		i = 1
	}else{
		i = i/j
		j = 1
	}
	return i,j
}