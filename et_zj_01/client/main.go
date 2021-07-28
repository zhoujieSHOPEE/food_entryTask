package main

import (
	"context"
	"fmt"
	"time"

	pb "et_zj_01/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

const (
	address     = "localhost:50051"
	port = ":81"
)

// 127.0.0.1/get_list?longitude=&latitude=&list_num=&city_id=
var conn *grpc.ClientConn
func main()  {
	router := gin.Default()
	conn, _ = grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024*4)))
	router.Use(func(c *gin.Context) {
		now := time.Now()
		c.Next()
		fmt.Println(time.Since(now).Seconds())
	})

	router.GET("/getBestStoresList", getBestStoresList)
	router.Run(port)

}

func getBestStoresList(c *gin.Context) {

	longitude := stringToFloat64(c.Query("longitude"))
	latitude := stringToFloat64(c.Query("latitude"))
	listNum := stringToInt32(c.Query("list_num"))
	cityId := stringToInt32(c.Query("city_id"))

	client := pb.NewGetBestStoresServiceClient(conn)

	pos := pb.Position{
		Longitude: longitude,
		Latitude: latitude,
	}

	r, err := client.GetBestStoresList(context.Background(), &pb.OutletRequest{Pos:&pos, ListNum:listNum, CityId:cityId})
	if err != nil {
		log.Fatalf("could not getStoresList: %v", err)
	}
	log.Printf("getStoresList: %d", r.GetCode())
	for i,v := range r.List{
		if i == int(listNum){
			break
		}
		//fmt.Println(v)
		c.IndentedJSON(200, gin.H{
			"name": v.Name,
			"distance": v.Distance,
			"logo_url": v.LogoURL,
			"address": v.Address,
			"itemsSold": v.ItemsSold,
		})
	}

}

func stringToFloat64(str string) float64 {
	float,_ := strconv.ParseFloat(str,64)
	return float
}

func stringToInt32(str string) int32 {
	i,_ := strconv.ParseInt(str, 10, 32)
	return int32(i)
}
