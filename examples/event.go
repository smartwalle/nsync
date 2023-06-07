package main

import (
	"fmt"
	"github.com/smartwalle/nsync"
	"time"
)

func main() {
	var event = nsync.NewEvent()

	var done = nsync.NewEvent()

	go func() {
		select {
		case <-event.Done():
			fmt.Println("1 done")
			done.Fire()
		}
	}()

	go func() {
		select {
		case <-event.Done():
			fmt.Println("2 done")
			done.Fire()
		}
	}()

	time.Sleep(time.Second * 2)

	event.Fire()

	<-done.Done()
}
