package predictor

import (
	"errors"
	"fmt"
	"math"
)

//Predictor Interface for any predictor
type PredictionLogic interface {
	GetPredictorName() string
	Predict(pastArray []float32) (predictedArray []float32, s error)
}

type WaveletTransform struct {
	//wavelet transforms have advantages over Fourier transforms in analyzing acyclic patterns.

	SlidingWindow    int // Parameter D to be used for knowing the windowsize of pastArray
	PredictionWindow int // Parameter W to be used for knowing how many values to predict

}

// Predictor Funcion to predict which takes input Prediction Logic
func Predictor(p PredictionLogic, pastArray []float32) (predictedArray []float32, err error) {
	fmt.Println("Predictor Name : ", p.GetPredictorName())
	predictedArray, err = p.Predict(pastArray)
	if err != nil {
		return
	}
	negativeValuesFixer(predictedArray[:])

	return
}

func (haar *WaveletTransform) GetPredictorName() string {
	return "Haar Wavelet Transform"
}

//Used to Predict WaveletTransform
func (haar *WaveletTransform) Predict(pastArray []float32) (predictedArray []float32, err error) {
	// Check the Sliding windowsize
	//if even continue

	if !isPowerOfTwo(haar.SlidingWindow) {
		fmt.Println("Sliding window size configured ", haar.SlidingWindow)
		return pastArray, errors.New("  Sliding number has to be power of 2 for Haar Wavelet")
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

// Utility function to check if number is power of two
func isPowerOfTwo(num int) bool {
	for num >= 2 {
		if num%2 != 0 {
			return false
		}
		num = num / 2
	}
	return num == 1
}
