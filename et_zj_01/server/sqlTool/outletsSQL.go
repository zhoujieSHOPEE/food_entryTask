package sqlTool

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Outlets struct {
	Id int `gorm:"not null;unique"`
	Name string `gorm:"column:Name"`
	Status int `gorm:"column:Status"`
	IsOfflinePaymentSupported int	`gorm:"column:is_offline_payment_supported"`
	LogoURL string `gorm:"column:Logo_url"`
	ImageURLList string
	CityId int `gorm:"column:City_id"`
	City string `gorm:"column:City"`
	District string
	Longitude float64
	Latitude float64
	Address string
	ItemsSold int
	MerchantId int
	MsOutletId int
	CreateTime int `gorm:"column:Create_time"`
	UpdateTime int `gorm:"column:Update_time"`
	DisplayStatus int `gorm:"column:Display_status"`
	TypeId int `gorm:"column:Type_id"`
	BrandId int `gorm:"column:Brand_id"`
	Location string `gorm:"column:Location"`
	LocationId int `gorm:"column:Location_id"`
	IsBScanCPaymentSupported int
	IsCScanBpaymentSupported int
	CardImage string
	HeadImages string
	Chemas string `gorm:"column:chemas"`
	Sharding string
	Dist float64
}

func (o Outlets) MarshalBinary() ([]byte, error) {
	bytes, err := json.Marshal(o)
	return bytes, err
}

func (o Outlets) UnmarshalBinary(data []byte) error{
	err := json.Unmarshal(data, o)
	return err
}

var db *sql.DB

func init() {
	initDB() // 调用输出化数据库的函数
}

// 定义一个初始化数据库的函数
func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:mysql123456@tcp(127.0.0.1:3306)/et_zj?charset=utf8mb4&parseTime=True"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(64)
	db.SetMaxOpenConns(64)
	db.SetConnMaxLifetime(time.Minute)
	return nil
}

func queryMultiRowDemoByCity(CityId int) ([]Outlets, error){
	sqlStr := "select distinct id, Name, longitude, latitude, Logo_url, address, items_sold from outlets where City_id = ?"
	fmt.Println(CityId)
	rows, err := db.Query(sqlStr, CityId)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil, err
	}
	fmt.Println(CityId)
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据

	//var outletsList *list.List = list.New()
	var outletsSlice []Outlets
	for rows.Next() {
		var o Outlets
		err := rows.Scan(&o.Id, &o.Name, &o.Longitude, &o.Latitude, &o.LogoURL, &o.Address, &o.ItemsSold)
		//outletsList.PushFront(o)
		outletsSlice = append(outletsSlice, o)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil, err
		}
		//fmt.Printf("Name:%s Longitude:%f latitude:%f\n", o.Name, o.Longitude, o.Latitude)
	}
	return outletsSlice, nil
}

func queryMultiRowDemoById(Id int) (Outlets, error){
	sqlStr := "select distinct id, Name, longitude, latitude, Logo_url, address, items_sold from outlets where id = ?"
	//fmt.Println(Id)
	row := db.QueryRow(sqlStr, Id)
	//fmt.Println(Id)
	// 非常重要：关闭rows释放持有的数据库链接
	var o Outlets
	err := row.Scan(&o.Id, &o.Name, &o.Longitude, &o.Latitude, &o.LogoURL, &o.Address, &o.ItemsSold)
	if err != nil {
		log.Printf("scan failed, err:%v\n", err)
	}

	return o, nil
}

func queryMultiRowDemo() ([]Outlets, error){
	sqlStr := "select distinct id, Name, longitude, latitude, Logo_url, address, items_sold from outlets"
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil, err
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据

	//var outletsList *list.List = list.New()
	var outletsSlice []Outlets
	for rows.Next() {
		var o Outlets
		err := rows.Scan(&o.Id, &o.Name, &o.Longitude, &o.Latitude, &o.LogoURL, &o.Address, &o.ItemsSold)
		//outletsList.PushFront(o)
		outletsSlice = append(outletsSlice, o)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil, err
		}
		//fmt.Printf("Name:%s Longitude:%f latitude:%f\n", o.Name, o.Longitude, o.Latitude)
	}
	return outletsSlice, nil
}

func FindOutletsByCityId(CityId int) ([]Outlets, error){

	outletsSlice, err := queryMultiRowDemoByCity(CityId)
	if err != nil {
		fmt.Printf("get outletsList fail,err:%v\n", err)
		return nil, err
	}
	return outletsSlice, nil
}

func FindOutletsById(Id int) (Outlets, error){

	outlets, err := queryMultiRowDemoById(Id)
	if err != nil {
		log.Printf("get outlets by id fail : %v", err)
	}
	return outlets, nil
}

func FindAllOutlets() ([]Outlets, error){

	outletsSlice, err := queryMultiRowDemo()
	if err != nil {
		fmt.Printf("get outletsList fail,err:%v\n", err)
		return nil, err
	}
	return outletsSlice, nil
}