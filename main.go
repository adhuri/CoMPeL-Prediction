package main

import (
	"fmt"

	predictor "github.com/adhuri/Compel-Prediction/predictor"
)

func main() {

	D := 32 //SlidingWindow
	W := 16 //PredictionWindow

	// Logic to start prediction
	fmt.Print("Prediction started\n")
	haar := predictor.WaveletTransform{SlidingWindow: D, PredictionWindow: W}

	//Sample Data
	a := []float32{6, 5, 4, 4, 4, 3, 4, 4, 3, 4, 3, 5, 3, 4, 4, 5, 3, 3, 5, 4, 3, 3, 5, 7, 4, 5, 5, 4, 4, 5, 5, 3}

	predictedArray, err := predictor.Predictor(&haar, a)
	if err != nil {
		fmt.Print("\nError received from Predictor ", err)
	}
	fmt.Print("\nPredicted Array ", predictedArray, "\n")

}
