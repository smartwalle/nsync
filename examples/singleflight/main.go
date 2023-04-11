package main

import (
	"fmt"
	"github.com/smartwalle/nsync/singleflight"
	"time"
)

func main() {
	var g1 = singleflight.New()

	go func() {
		fmt.Println("Goroutine---1")
		var v, _ = g1.Do("k1", func(key string) (interface{}, error) {
			fmt.Println("Goroutine---1: begin")
			time.Sleep(time.Second * 2)
			fmt.Println("Goroutine---1: end")
			return "v1", nil
		})
		fmt.Println("Goroutine---1 结果:", v)
	}()

	var g2 = g1

	go func() {
		fmt.Println("Goroutine---2")
		var v, _ = g2.Do("k1", func(key string) (interface{}, error) {
			fmt.Println("Goroutine---2: begin")
			time.Sleep(time.Second * 2)
			fmt.Println("Goroutine---2: end")
			return "v2", nil
		})
		fmt.Println("Goroutine---2 结果:", v)
	}()

	select {
	case <-time.After(time.Second * 10):
	}
}
