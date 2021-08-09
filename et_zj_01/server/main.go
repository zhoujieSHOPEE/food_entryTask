package main

import (
	"context"
	pb "et_zj_01/proto"
	ent "et_zj_01/server/entity"
	ft "et_zj_01/server/funcTool"
	st "et_zj_01/server/sqlTool"
	"fmt"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
	"log"
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

var OutLetsMap = make(map[int]ent.Outlets)
var OutletsIdSlice = make([]int,0)
var count = 0
func initServer()  {
	OutletsSlice, err := st.FindAllOutlets()
	if err != nil {
		log.Fatalf("server main init fial : %v", err)
	}
	for _, v := range OutletsSlice {
		OutLetsMap[v.Id] = v
	}
	OutletsIdSlice, err = st.FindOrderedOutletsWithLimit(100000)
	if err != nil {
		log.Fatalf("server main FindOrderedOutletsWithLimit fail : %v", err)
	}
}
func (*server) GetBestStoresList (ctx context.Context, req *pb.OutletRequest) (*pb.OutletResponse, error) {
	pos := req.Pos
	method := req.Method
	count += 1
	if count % 3 == 0 {

	}
	outletsSlice := make([]ent.Outlets, 0)
	var err error

	if method == 1 {
		fmt.Println("method : Distance")
		now := time.Now()
		outletsSlice, err = st.GetNearStore(pos.Longitude, pos.Latitude, OutLetsMap, count)
		fmt.Println("GetNearStore : ", time.Since(now))

		if outletsSlice == nil {
			outletsSlice, err = st.GetTopSale(int(req.ListNum), pos.Longitude, pos.Latitude, OutLetsMap, OutletsIdSlice)
		}
	}else if method == 2 {
		fmt.Println("method : itemsSold")
		outletsSlice, err = st.GetTopSale(int(req.ListNum), pos.Longitude, pos.Latitude, OutLetsMap, OutletsIdSlice)
	}else {
		fmt.Println("method : Distance + itemsSold")

		outletsSliceByItemsSold, _ := st.GetTopSale(int(req.ListNum), pos.Longitude, pos.Latitude, OutLetsMap, OutletsIdSlice)
		outletsSliceByDistance, _ := st.GetNearStore(pos.Longitude, pos.Latitude, OutLetsMap, count)
		//fmt.Println(len(outletsSliceByItemsSold), len(outletsSliceByDistance))
		outlets := append(outletsSliceByDistance, outletsSliceByItemsSold...)
		outlets = ft.RemoveRepByMap(outlets)
		outletsSlice = outlets
	}
	if err != nil {
		log.Fatalf("get top sale fail : %v", err)
	}

	outletsSlice = ft.SortByScore(outletsSlice, 0 ,len(outletsSlice))

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
	initServer()
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
	c := cron.New(cron.WithSeconds())
	spec := "00 02 03 * * ?"
	c.AddFunc(spec, func() {
		initServer()
	})
	c.Start()
}







