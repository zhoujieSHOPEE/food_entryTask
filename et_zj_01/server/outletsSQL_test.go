package main

import (
	"fmt"
	"testing"
)

func TestFindOutletsByCityId(t *testing.T) {
	outletsSlice, err := FindOutletsByCityId(1)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("get from list---------------------")
	for _,v := range outletsSlice{
		fmt.Println(v)
	}
}