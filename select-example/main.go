package main

import (
	"fmt"
	"sync"
	"time"
)

func routineFunction(ch chan string, wg *sync.WaitGroup) {
	time.Sleep(1 * time.Second)
	ch <- "Message"
	wg.Done()
}

func main() {
	ch := make(chan string)
	var wg sync.WaitGroup
	for index := 0; index < 10; index++ {
		wg.Add(1)
		go routineFunction(ch, &wg)
	}
	select {
	case x := <-ch:
		fmt.Println(x)
	default:
		fmt.Println("Waiting...")
	}
	wg.Wait()
}
