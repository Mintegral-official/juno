package main

import (
	"fmt"
	"github.com/Mintegral-official/juno/conf"
	"github.com/Mintegral-official/juno/model"
)

//
///**
// * @author: tangye
// * @Date: 2019/11/4 19:31
// * @Description:
// */
//
/////*
////#include <stdio.h>
////int t() {
////    return rand() % (1000000000 - 0 + 1) + 0;
////}
////*/
////import "C"
//import "C"
//import (
//	"fmt"
//	"github.com/Mintegral-official/juno/helpers"
//	"github.com/Mintegral-official/juno/index"
//	"time"
//	"unsafe"
//)
//
//var s = index.NewSkipListIterator(index.DEFAULT_MAX_LEVEL, helpers.IntCompare)
//// var s = New(helpers.IntCompare)
//var arr1 [200001]int
//var s1 = make([]int, 500001)
//
//func init() {
//	for i := 0; i < 200001; i++ {
//		arr1[i] = int(C.t())
//	}
//
//	for i := 0; i < 200001; i++ {
//		s.Add(arr1[i], nil)
//	}
//
//	var sl index.SkipList
//	var el index.Element
//	fmt.Printf("Structure sizes: SkipList is %v, Element is %v bytes\n", unsafe.Sizeof(sl), unsafe.Sizeof(el))
//}
//
////func t2() {
////	for i := 0; i <= 20000; i++ {
////		v, _ := s.Get(arr1[i])
////		fmt.Println("****2**** %d ", v)
////	}
////}
////func t3() {
////	for i := 20001; i <= 40000; i++ {
////		v, _ := s.Get(arr1[i])
////		fmt.Println("****3**** %d ", v)
////	}
////}
////func t4() {
////	for i := 40001; i <= 60000; i++ {
////		v, _ := s.Get(arr1[i])
////		fmt.Println("****4**** %d ", v)
////	}
////}
////func t5() {
////	for i := 60001; i <= 80000; i++ {
////		v, _ := s.Get(arr1[i])
////		fmt.Println("****5**** %d ", v)
////	}
////}
////func t6() {
////	for i := 80001; i <= 100000; i++ {
////		v, _ := s.Get(arr1[i])
////		fmt.Println("****6**** %d ", v)
////	}
////}
////func t7() {
////	for i := 100001; i <= 120000; i++ {
////		v, _ := s.Get(arr1[i])
////		fmt.Println("****7**** %d ", v)
////	}
////}
////func t8() {
////	for i := 120001; i <= 140000; i++ {
////		v, _ := s.Get(arr1[i])
////		fmt.Println("****8**** %d ", v)
////	}
////}
////func t9() {
////	for i := 140001; i <= 160000; i++ {
////		v, _ := s.Get(arr1[i])
////		fmt.Println("****9**** %d ", v)
////	}
////}
////func t10() {
////	for i := 160001; i <= 180000; i++ {
////		v, _ := s.Get(arr1[i])
////		fmt.Println("****10**** %d ", v)
////	}
////}
////func t11() {
////	for i := 180001; i <= 200000; i++ {
////		v, _ := s.Get(arr1[i])
////		fmt.Println("****11**** %d ", v)
////	}
////}
////func t12() {
////	for i := 0; i < 200001; i++ {
////		s.Add(arr1[i], nil)
////	}
////}
////func t1() {
////	go func() {
////        for i := 0; i < 200000; i++ {
////        	time.Sleep(250 * time.Microsecond)
////        	//s.Get(arr1[i], nil)
////        	fmt.Printf("111    %d  ", s.Len())
////		}
////	}()
////    time.Sleep(10000)
////
////    go func() {
////		for i := 0; i < 200000; i++ {
////			time.Sleep(3000 * time.Microsecond)
////			v, _ := s.Get(arr1[i])
////			fmt.Printf("222    %d  ", v)
////		}
////	}()
////
////}
//
//func main() {
//	//go t12()
//	//go t2()
//	//go t3()
//	//go t4()
//	//go t5()
//	//go t6()
//	//go t7()
//	//go t8()
//	//go t9()
//	//go t10()
//	//go t11()
//	time.Sleep(15 * time.Second)
//	fmt.Println("\n*****************\n")
//	//fmt.Println(s.Len())
//	c := 0
//	for s.HasNext() {
//		s.Next()
//		c++
//	}
//	fmt.Println("\n*********\n")
//	fmt.Println(s.Len())
//	fmt.Println(c)
//
//}

func main() {

	cfg := &conf.MongoCfg{
		URI:            "mongodb://localhost:27017",
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
	}

	mon, err := model.NewMongo(cfg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mon)
	r, e := mon.Find()
	if err == nil {
		fmt.Println(r)
	}
	fmt.Println(e)

	for i := 0; i < len(r); i++ {
		fmt.Println(r[i].CampaignId)
	}

}