package sqlTool

import (
	"database/sql"
	ent "et_zj_01/server/entity"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)



//func (o Outlets) MarshalBinary() ([]byte, error) {
//	bytes, err := json.Marshal(o)
//	return bytes, err
//}
//
//func (o Outlets) UnmarshalBinary(data []byte) error{
//	err := json.Unmarshal(data, o)
//	return err
//}

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

func queryMultiRowDemoByCity(CityId int) ([]ent.Outlets, error){
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
	var outletsSlice []ent.Outlets
	for rows.Next() {
		var o ent.Outlets
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

func queryMultiRowDemoById(Id int) (ent.Outlets, error){
	sqlStr := "select distinct id, Name, longitude, latitude, Logo_url, address, items_sold from outlets where id = ?"
	//fmt.Println(Id)
	row := db.QueryRow(sqlStr, Id)
	//fmt.Println(Id)
	// 非常重要：关闭rows释放持有的数据库链接
	var o ent.Outlets
	err := row.Scan(&o.Id, &o.Name, &o.Longitude, &o.Latitude, &o.LogoURL, &o.Address, &o.ItemsSold)
	if err != nil {
		log.Printf("scan failed, err:%v\n", err)
	}

	return o, nil
}

func queryMultiRowDemo() ([]ent.Outlets, error){
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
	var outletsSlice []ent.Outlets
	for rows.Next() {
		var o ent.Outlets
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

func queryMultiRowDemoWithLimit(l int) ([]int, error){
	sqlStr := "select id from outlets order by items_sold desc limit ?"
	rows, err := db.Query(sqlStr, l)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil, err
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据

	//var outletsList *list.List = list.New()
	var idSlice []int
	for rows.Next() {
		var i int
		err := rows.Scan(&i)
		idSlice = append(idSlice, i)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil, err
		}
	}
	return idSlice, nil
}

func FindOutletsByCityId(CityId int) ([]ent.Outlets, error){

	outletsSlice, err := queryMultiRowDemoByCity(CityId)
	if err != nil {
		fmt.Printf("get outletsList fail,err:%v\n", err)
		return nil, err
	}
	return outletsSlice, nil
}

func FindOutletsById(Id int) (ent.Outlets, error){

	outlets, err := queryMultiRowDemoById(Id)
	if err != nil {
		log.Printf("get outlets by id fail : %v", err)
	}
	return outlets, nil
}

func FindAllOutlets() ([]ent.Outlets, error){

	outletsSlice, err := queryMultiRowDemo()
	if err != nil {
		fmt.Printf("get outletsList fail,err:%v\n", err)
		return nil, err
	}
	return outletsSlice, nil
}

func FindOrderedOutletsWithLimit(limit int) ([]int, error){
	idSlice, err := queryMultiRowDemoWithLimit(limit)
	if err != nil {
		fmt.Printf("get outletsList fail,err:%v\n", err)
		return nil, err
	}
	return idSlice, nil
}