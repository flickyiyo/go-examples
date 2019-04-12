package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type MassPoint struct {
	x, y, z, mass int
}

func addMassPoints(a, b MassPoint) MassPoint {
	return MassPoint{
		a.x + b.x,
		a.y + b.y,
		a.z + b.z,
		a.mass + b.mass,
	}
}

func avgMassPoints(a, b MassPoint) MassPoint {
	sum := addMassPoints(a, b)
	return MassPoint{
		sum.x / 2,
		sum.y / 2,
		sum.z / 2,
		sum.mass,
	}
}

func toWeightedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		a.x * a.mass,
		a.y * a.mass,
		a.z * a.mass,
		a.mass,
	}
}

func fromWeightedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		a.x / a.mass,
		a.y / a.mass,
		a.z / a.mass,
		a.mass,
	}
}

func avgMassPointsWeighted(a, b MassPoint) MassPoint {
	aWeighted := toWeightedSubspace(a)
	bWeighted := toWeightedSubspace(b)
	return fromWeightedSubspace(avgMassPoints(aWeighted, bWeighted))
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func closeFile(fi *os.File) {
	err := fi.Close()
	handle(err)
}

func stringToPointAsync(s string, c chan<- MassPoint, wg *sync.WaitGroup) {
	defer wg.Done()
	var newMassPoint MassPoint
	_, err := fmt.Sscanf(s, "%d:%d:%d:%d", &newMassPoint.x, &newMassPoint.y, &newMassPoint.z, &newMassPoint.mass)
	if err != nil {
		return
	}
	c <- newMassPoint
}

func avgMassPointsWeightedAsync(a, b MassPoint, c chan<- MassPoint) {
	c <- avgMassPointsWeighted(a, b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Incorret number of args")
	}

	file, err := os.Open(os.Args[1])
	handle(err)
	defer closeFile(file)

	var massPoints []MassPoint
	startLoading := time.Now()

	r := bufio.NewReader(file)
	massPointsChan := make(chan MassPoint, 128)
	var wg sync.WaitGroup
	for {
		str, err := r.ReadString('\n')
		if len(str) == 0 || err != nil {
			break
		}
		wg.Add(1)
		go stringToPointAsync(str, massPointsChan, &wg)
	}

	syncChan := make(chan bool)
	go func() { wg.Wait(); syncChan <- false }()

	run := true

	for run || len(massPointsChan) > 0 {
		select {
		case value := <-massPointsChan:
			massPoints = append(massPoints, value)
		case _ = <-syncChan:
			run = false
		}
	}

	fmt.Printf("Loaded %d values form file in %s", len(massPoints), time.Since(startLoading))

	if len(massPoints) <= 1 {
		handle(errors.New("Insufficient"))
	}

	c := make(chan MassPoint, len(massPoints)/2)

	startCalculation := time.Now()

	for len(massPoints) > 1 {
		var newMassPoints []MassPoint
		goroutines := 0
		for i := 0; i < len(massPoints)-1; i += 2 {
			go avgMassPointsWeightedAsync(massPoints[i], massPoints[i+1], c)
			goroutines++
		}

		for i := 0; i < goroutines; i++ {
			newMassPoints = append(newMassPoints, <-c)
		}
		if len(massPoints)%2 != 0 {
			newMassPoints = append(newMassPoints, massPoints[len(massPoints)-1])
		}

		massPoints = newMassPoints
	}

	systemAverge := massPoints[0]
	fmt.Printf("Calculation took %s\n", time.Since(startCalculation))
	fmt.Println(systemAverge)
}
