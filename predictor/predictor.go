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
	//a := [6]float32{2, 3, 5, 7, 11, 13}
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

	predictedArray = pastArray
	return
}

func Predictor(p PredictionLogic, pastArray []float32) (predictedArray []float32, s error) {
	fmt.Print(p.GetPredictorName())
	predictedArray, s = p.Predict(pastArray)
	return
}
