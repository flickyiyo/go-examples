package main

import (
	"fmt"
	"time"
)

func main() {
	go pingpong()
	time.Sleep(3 * time.Second)
}

func ping(ball chan<- int, action chan<- string) {
	ball <- 1
	action <- "Player ping"

}

func pong(ball chan int, action chan string) {
	ball <- 2
	action <- "Player pong"
}

func referee(action <-chan string) {
	for {
		fmt.Println(<-action)
	}
}

func pingpong() {
	ball := make(chan int)
	action := make(chan string)
	go referee(action)
	go ping(ball, action)
	go pong(ball, action)
	for {
		value := <-ball
		switch value {
		case 1:
			go pong(ball, action)
		case 2:
			go ping(ball, action)
		}
	}
}
