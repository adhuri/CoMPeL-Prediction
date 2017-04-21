package predictor

import (
	"errors"

	"github.com/Sirupsen/logrus"
)

//Predictor Interface for any predictor
type PredictionLogic interface {
	GetPredictorName() string
	Predict(pastArray []float32, bin int, logic int, log *logrus.Logger) (predictedArray []float32, s error)
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
func Predictor(p PredictionLogic, pastArray []float32, bin int, logic int, log *logrus.Logger) (predictedArray []float32, err error) {
	log.Infoln("Predictor Name : ", p.GetPredictorName())
	predictedArray, err = p.Predict(pastArray, bin, logic, log)
	if err != nil {
		return
	}

	// Value Raiser set to 0
	valueRaiser(predictedArray[:], 0, log)

	// CPU , memory cannot be negative

	negativeValuesFixer(predictedArray[:], log)

	return
}

func (haar *WaveletTransform) GetPredictorName() string {
	return "Haar Wavelet Transform"
}

//Used to Predict WaveletTransform
func (haar *WaveletTransform) Predict(pastArray []float32, bin int, logic int, log *logrus.Logger) (predictedArray []float32, err error) {
	// Check the Sliding windowsize
	//if even continue

	if !isPowerOfTwo(haar.SlidingWindow) {
		log.Debugln("Sliding window size configured ", haar.SlidingWindow)
		return predictedArray, errors.New("  Sliding number has to be power of 2 for Haar Wavelet")
	}
	// Ignore the
	if len(pastArray) < haar.SlidingWindow {
		log.Debugln("No Prediction - Length of past array is smaller than Sliding Window ", len(pastArray))
		return
	}
	//Trim additional elements - Redundant code but dont believe other module - Safe side check
	if len(pastArray) > haar.SlidingWindow {
		log.Debugln("Length of pastarray is larger than Sliding Window - Trimming ")
		pastArray = append(pastArray[len(pastArray)-haar.SlidingWindow:]) // To trim from totallength- sliding window to end of array
	}

	predictedCoefficients := Haar(pastArray, haar.PredictionWindow, bin, logic, log)

	log.Debugln("Predicted coefficients array: ", predictedCoefficients)

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
