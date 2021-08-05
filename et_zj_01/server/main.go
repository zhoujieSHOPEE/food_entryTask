package main

import (
	"context"
	pb "et_zj_01/proto"
	st "et_zj_01/server/sqlTool"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math"
	"net"
	"strconv"
	"time"
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
	pos := req.Pos
	now := time.Now()
	outletsSliceByDistance, err := st.GetNearStore(pos.Longitude, pos.Latitude)
	fmt.Println("GetNearStore time : ", time.Since(now))
	//outletsSliceByTopSale, err := st.GetTopSale(int(req.ListNum), pos.Longitude, pos.Latitude)
	if err != nil {
		log.Fatalf("get top sale fail : %v", err)
	}

	fmt.Println(outletsSliceByDistance[0])
	var retMessageList []*pb.RetMessage
	for _,v := range outletsSliceByDistance{
		itemsSold := strconv.Itoa(v.ItemsSold)
		//dist := GeoDistance(pos.Longitude, pos.Latitude, v.Longitude,v.Latitude)
		retMessage := &pb.RetMessage{
			Name: v.Name,
			Address: v.Address,
			LogoURL: v.LogoURL,
			ItemsSold: itemsSold,
			Distance: fmt.Sprintf("%fkm", v.Dist),
		}


		retMessageList = append(retMessageList, retMessage)
	}
	//s := RetMessageWrapper{retMessageList}
	////sort.Sort(s)
	//fmt.Println(sort.IsSorted(s))

	retMessageList = Sort(retMessageList, 0, len(retMessageList))

	res := &pb.OutletResponse{Code: 0,List: retMessageList,ListNum: req.ListNum}
	return res, err
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

	itemsSold1, itemsSold2 = Normalization(itemsSold1, itemsSold2)
	//dis1, dis2 = Normalization(dis1, dis2)

	score1 := itemsSold1 + DistScore(dis1)
	score2 := itemsSold2 + DistScore(dis2)

	fmt.Println(score1, score2)

	return score1 > score2
}

func Normalization(i, j float64) (float64, float64) {
	/*
	对两个数归一化
	 */
	if i == 0 {
		return 0,1
	}
	if j == 0{
		return 1,0
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


func stringTOFloat64(s string) float64 {

	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func DistScore(dist float64) float64 {
	if dist > 0 && dist < 2 {
		return 1
	}else if dist >= 2 && dist < 5 {
		return 0.8
	}else if dist >=5 && dist < 10 {
		return 0.6
	}else {
		return 0.2
	}
}

func Sort(r []*pb.RetMessage, left, right int) []*pb.RetMessage {

	for i := left; i < right; i++ {
		for j := i+1; j < right; j++ {
			itemSoldi := stringTOFloat64(r[i].ItemsSold)
			itemSoldj := stringTOFloat64(r[j].ItemsSold)
			disti := stringTOFloat64(r[i].Distance)
			distj := stringTOFloat64(r[j].Distance)
			if itemSoldj/(distj+2) > itemSoldi/(disti+2) {
				r[i], r[j] = r[j], r[i]
			}
		}
	}

	return r
}