package predictor

import (
	"errors"
	"fmt"
	"math"
)

//Predictor Interface for any
type PredictionLogic interface {
	GetPredictorName() string
	Predict(pastArray []float32) (predictedArray []float32, s error)
}

type WaveletTransform struct {
	//wavelet transforms have advantages over Fourier transforms in analyzing acyclic patterns.

	SlidingWindow    int // Parameter D to be used for knowing the windowsize of pastArray
	PredictionWindow int // Parameter W to be used for knowing how many values to predict

}

func (haar *WaveletTransform) GetPredictorName() string {
	return "Haar Wavelet Transform"
}
func (haar *WaveletTransform) Predict(pastArray []float32) (predictedArray []float32, err error) {

	fmt.Print("\nSliding Window size ", haar.SlidingWindow)
	// Check the Sliding windowsize
	//if even continue
	if (haar.SlidingWindow)%2 != 0 {
		return pastArray, errors.New("\nSlidingWindow has to be an even number for Haar - Given :" + string(haar.SlidingWindow) + "\n")
	}
	//Trim additional elements
	if len(pastArray) > haar.SlidingWindow {
		fmt.Print("Length of pastarray is larger than Sliding Window")
		pastArray = append(pastArray[:haar.SlidingWindow])
	}

	transformedArray := Haar(pastArray, int(math.Log2(float64(haar.PredictionWindow))))
	fmt.Println("Transformed array: ", transformedArray)

	invertedArray := Inverse_haar(transformedArray)
	predictedArray = invertedArray
	return
}

// Predictor Funtion to predict which takes input Prediction Logic
func Predictor(p PredictionLogic, pastArray []float32) (predictedArray []float32, s error) {
	fmt.Print(p.GetPredictorName())
	predictedArray, s = p.Predict(pastArray)
	negativeValuesFixer(predictedArray[:])

	return
}

//Predictions are not accurate and for near zero values could predict negative Values. Fixing them to zero
// All it means is the value approaches zero
func negativeValuesFixer(result []float32) {
	fixedCount := 0
	for i, el := range result {
		if el < 0 {
			result[i] = 0
			fixedCount += 1
		}
	}
	fmt.Println("Fixed values in predicted array ", fixedCount)
}
