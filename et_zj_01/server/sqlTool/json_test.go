package sqlTool

//import (
//	"encoding/json"
//	"fmt"
//	"testing"
//	"time"
//)
//func TestJson(t *testing.T) {
//
//	o1 := Outlets{Id: 111, Name: "tom"}
//
//	datas, _ := json.Marshal(o1)
//	fmt.Println("before insert redis : ", datas)
//	rdb.Set(ctx, "t1", datas, 10*time.Minute)
//
//	result, _ := rdb.Get(ctx, "t1").Bytes()
//	fmt.Println("after redis : ", result)
//
//	o2 := &Outlets{}
//	json.Unmarshal(result, o2)
//	fmt.Println(o2)
//}