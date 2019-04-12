package main

import (
	"merge-concurrent/constants"
)

func main() {
	arreglo := constants.GetArreglo()
	for index := 0; index < 100; index++ {
		arreglo = append(arreglo, arreglo...)
	}
	channel := make(chan []int)
	defer close(channel)
	go merge1(append(arreglo, arreglo...), channel)
	<-channel

}

func merge1(m []int, channel chan []int) {
	if len(m) <= 1 {
		channel <- m
		return
	}
	mid := len(m) / 2
	channelLeft := make(chan []int)
	channelRight := make(chan []int)
	go merge1(m[:mid], channelLeft)
	go merge1(m[mid:], channelRight)
	left := <-channelLeft
	right := <-channelRight
	close(channelLeft)
	close(channelRight)
	channel <- merge2(left, right)
	return
}

func merge2(left []int, right []int) []int {
	list := []int{}
	for len(left) > 0 && len(right) > 0 {
		if left[0] < right[0] {
			list = append(list, left[0])
			left = left[1:]
		} else {
			list = append(list, right[0])
			right = right[1:]
		}
	}
	if len(left) > 0 {
		list = append(list, left...)
	}
	if len(right) > 0 {
		list = append(list, right...)
	}
	return list
}
