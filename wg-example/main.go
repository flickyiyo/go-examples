package main

import (
	"fmt"
	"sync"
	"time"
)

func myFunc(wg *sync.WaitGroup, goRoutineNumber int) {
	time.Sleep(1 * time.Second)
	fmt.Println("Finishing myFunc", goRoutineNumber)
	wg.Done()
}

func main() {
	fmt.Println("Init")
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go myFunc(&wg, i)
	}
	wg.Wait()
	fmt.Println("Finish")
}
