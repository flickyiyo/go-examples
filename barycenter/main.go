package main

import (
	"errors"
	"fmt"
	"io"
	"os"
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Incorret number of args")
	}

	file, err := os.Open(os.Args[1])
	handle(err)
	defer closeFile(file)

	var massPoints []MassPoint
	startLoading := time.Now()

	for {
		var newMassPoint MassPoint
		_, err = fmt.Fscanf(file, "%d:%d:%d:%d", &newMassPoint.x, &newMassPoint.y, &newMassPoint.z, &newMassPoint.mass)
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}

		massPoints = append(massPoints, newMassPoint)
	}

	fmt.Printf("Loaded %d values form file in %s", len(massPoints), time.Since(startLoading))

	if len(massPoints) <= 1 {
		handle(errors.New("Insufficient"))
	}

	startCalculation := time.Now()

	for len(massPoints) > 1 {
		var newMassPoints []MassPoint
		for i := 0; i < len(massPoints)-1; i += 2 {
			newMassPoints = append(newMassPoints, avgMassPointsWeighted(massPoints[i], massPoints[i+1]))
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
