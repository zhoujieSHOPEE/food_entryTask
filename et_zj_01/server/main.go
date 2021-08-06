package main

import (
	"context"
	pb "et_zj_01/proto"
	st "et_zj_01/server/sqlTool"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

type server struct {
	//pb.UnimplementedGetBestStoresServiceServer
}

const (
	port = ":50051"
)

var OutLetsMap map[int]st.Outlets

func init()  {
	outletsSlice, err := st.FindAllOutlets()
	if err != nil {
		log.Fatalf("server main init fial : %v", err)
	}
	for _, v := range outletsSlice {
		OutLetsMap[v.Id] = v
	}
}
func (*server) GetBestStoresList (ctx context.Context, req *pb.OutletRequest) (*pb.OutletResponse, error) {

	pos := req.Pos
	method := req.Method

	outletsSlice := make([]st.Outlets, 0)
	var err error

	if method == 1 {
		fmt.Println("method : Distance")
		outletsSlice, err = st.GetNearStore(pos.Longitude, pos.Latitude)
	}else if method == 2 {
		fmt.Println("method : itemsSold")
		outletsSlice, err = st.GetTopSale(int(req.ListNum), pos.Longitude, pos.Latitude)
	}else {
		fmt.Println("method : Distance + itemsSold")

		outletsSliceByItemsSold, _ := st.GetTopSale(int(req.ListNum), pos.Longitude, pos.Latitude)
		outletsSliceByDistance, _ := st.GetNearStore(pos.Longitude, pos.Latitude)
		fmt.Println(len(outletsSliceByItemsSold), len(outletsSliceByDistance))
		outlets := append(outletsSliceByDistance, outletsSliceByItemsSold...)
		outlets = removeRepByMap(outlets)
		outletsSlice = outlets
	}
	if err != nil {
		log.Fatalf("get top sale fail : %v", err)
	}

	outletsSlice = Sort(outletsSlice, 0 ,len(outletsSlice))

	var retMessageList []*pb.RetMessage
	for _,v := range outletsSlice{
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

	res := &pb.OutletResponse{Code: 0,List: retMessageList,ListNum: req.ListNum}
	return res, err
}

func main(){
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
}

func stringToFloat64(s string) float64 {

	f, _ := strconv.ParseFloat(s, 64)
	return f
}

//func Sort(r []*pb.RetMessage, left, right int) []*pb.RetMessage {
//
//	for i := left; i < right; i++ {
//		for j := i+1; j < right; j++ {
//			itemSoldi := stringToFloat64(r[i].ItemsSold)
//			itemSoldj := stringToFloat64(r[j].ItemsSold)
//			disti := stringToFloat64(r[i].Distance)
//			distj := stringToFloat64(r[j].Distance)
//			if itemSoldj/(distj+2) > itemSoldi/(disti+2) {
//				r[i], r[j] = r[j], r[i]
//			}
//		}
//	}
//
//	return r
//}

func Sort(o []st.Outlets, left, right int) []st.Outlets {

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

func removeRepByMap(slc []st.Outlets) []st.Outlets {
	result := make([]st.Outlets, 0)       //存放返回的不重复切片
	tempMap := map[st.Outlets]byte{} // 存放不重复主键
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
