package main

import (
	"fmt"
	"github.com/smartwalle/nsync/singleflight"
	"time"
)

func main() {
	var g1 = singleflight.New()

	go func() {
		var v, _ = g1.Do("k1", func(key string) (interface{}, error) {
			fmt.Println("begin1")
			time.Sleep(time.Second * 2)
			fmt.Println("end1")
			return "v1", nil
		})
		fmt.Println(v)
	}()

	var g2 = *g1

	go func() {
		var v, _ = g2.Do("k1", func(key string) (interface{}, error) {
			fmt.Println("begin2")
			time.Sleep(time.Second * 2)
			fmt.Println("end2")
			return "v2", nil
		})
		fmt.Println(v)
	}()

	select {
	case <-time.After(time.Second * 10):
	}
}
