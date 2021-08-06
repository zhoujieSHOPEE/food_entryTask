package main

import (
	"context"
	pb "et_zj_01/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
)

const (
	address     = "localhost:50051"
	port = ":81"
)


var conn *grpc.ClientConn
func main()  {
	router := gin.Default()
	var err error
	conn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024*4)))
	if err != nil {
		log.Fatalf("grpc Dial fail : %v", err)
	}
	router.GET("/getBestStoresList", getBestStoresList)
	router.Run(port)

}

func getBestStoresList(c *gin.Context) {

	longitude := stringToFloat64(c.Query("longitude"))
	latitude := stringToFloat64(c.Query("latitude"))
	listNum := stringToInt32(c.Query("list_num"))
	method := stringToInt32(c.Query("method"))
	client := pb.NewGetBestStoresServiceClient(conn)

	pos := pb.Position{
		Longitude: longitude,
		Latitude: latitude,
	}

	r, err := client.GetBestStoresList(context.Background(), &pb.OutletRequest{Pos:&pos, ListNum:listNum, Method: method})
	if err != nil {
		log.Fatalf("could not getStoresList: %v", err)
	}
	for i,v := range r.List{
		if i == int(15){
			break
		}
		c.IndentedJSON(http.StatusOK, gin.H{
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
