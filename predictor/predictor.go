package predictor

import (
	"errors"
	"fmt"
)

var debug = true

//Predictor Interface for any predictor
type PredictionLogic interface {
	GetPredictorName() string
	Predict(pastArray []float32, bin int, logic int) (predictedArray []float32, s error)
}

type WaveletTransform struct {
	//wavelet transforms have advantages over Fourier transforms in analyzing acyclic patterns.

	SlidingWindow    int // Parameter D to be used for knowing the windowsize of pastArray
	PredictionWindow int // Parameter W to be used for knowing how many values to predict

}

type MaxPredict struct {
	// Choose all predicted values = max of the sliding window
	SlidingWindow    int // Parameter D to be used for knowing the windowsize of pastArray
	PredictionWindow int // Parameter W to be used for knowing how many values to predict

}

// Predictor Funcion to predict which takes input Prediction Logic
func Predictor(p PredictionLogic, pastArray []float32, bin int, logic int) (predictedArray []float32, err error) {
	fmt.Println("Predictor Name : ", p.GetPredictorName())
	predictedArray, err = p.Predict(pastArray, bin, logic)
	if err != nil {
		return
	}

	// Value Raiser set to 0
	valueRaiser(predictedArray[:], 0)

	// CPU , memory cannot be negative

	negativeValuesFixer(predictedArray[:])

	return
}

func (haar *WaveletTransform) GetPredictorName() string {
	return "Haar Wavelet Transform"
}

//Used to Predict WaveletTransform
func (haar *WaveletTransform) Predict(pastArray []float32, bin int, logic int) (predictedArray []float32, err error) {
	// Check the Sliding windowsize
	//if even continue

	if !isPowerOfTwo(haar.SlidingWindow) {
		fmt.Println("Sliding window size configured ", haar.SlidingWindow)
		return predictedArray, errors.New("  Sliding number has to be power of 2 for Haar Wavelet")
	}
	// Ignore the
	if len(pastArray) < haar.SlidingWindow {
		fmt.Println("No Prediction - Length of past array is smaller than Sliding Window ", len(pastArray))
		return
	}
	//Trim additional elements - Redundant code but dont believe other module - Safe side check
	if len(pastArray) > haar.SlidingWindow {
		fmt.Println("Length of pastarray is larger than Sliding Window - Trimming ")
		pastArray = append(pastArray[len(pastArray)-haar.SlidingWindow:]) // To trim from totallength- sliding window to end of array
	}

	predictedCoefficients := Haar(pastArray, haar.PredictionWindow, bin, logic)

	if debug {
		fmt.Println("Predicted coefficients array: ", predictedCoefficients)
	}

	invertedArray := InverseHaar(predictedCoefficients)
	predictedArray = invertedArray
	return
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

// Max values

func (mp *MaxPredict) GetPredictorName() string {
	return "Haar Wavelet Transform"
}

//Used to Predict WaveletTransform
func (mp *MaxPredict) Predict(pastArray []float32, bin int, logic int) (predictedArray []float32, err error) {
	_, max := findMinMax(pastArray)
	for i := 0; i < mp.PredictionWindow; i++ {
		predictedArray = append(predictedArray, max)
	}

	return
}
