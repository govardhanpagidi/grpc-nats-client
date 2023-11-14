package main

import (
	"math/rand"
	"time"
)

func currentTime() string {
	return time.Now().Format("2006-01-02T15:04:05.000")
}

func randomFloat() float64 {
	rand.Seed(time.Now().UnixNano())

	// Define the range for the random float (e.g., between 1.0 and 10.0)
	minVal := 1.0
	maxVal := 100000.0

	// Generate a random float within the specified range
	return roundToTwoDecimalPlaces(minVal + rand.Float64()*(maxVal-minVal))
}

func roundToTwoDecimalPlaces(num float64) float64 {
	return float64(int(num*100)) / 100
}
