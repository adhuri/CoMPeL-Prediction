package predictor

import "github.com/Sirupsen/logrus"

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
	log.Infoln("--> Predictor Name : ", p.GetPredictorName())
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

//Prediction is not accurate and for near zero values could predict negative Values. Fixing them to zero
// All it means is the value approaches zero
func negativeValuesFixer(result []float32, log *logrus.Logger) {
	fixedCount := 0
	for i, el := range result {
		if el < 0 {
			result[i] = 0
			fixedCount += 1
		}
	}

	log.Debugln("Fixed negative values in predicted array ", fixedCount)

}

// To create a variation for under prediction
func valueRaiser(result []float32, valueRaisedPercentage float32, log *logrus.Logger) {
	//_, max := findMinMax(result)
	//fromMaxValueEnhancer := 100 - max
	for i, _ := range result {
		fromValueRaiser := (valueRaisedPercentage / 100) * (result[i])

		result[i] = result[i] + fromValueRaiser

		// if fromMaxValueEnhancer > fromValueRaiser {
		// 	result[i] = result[i] + fromMaxValueEnhancer
		// } else {
		// 	result[i] = result[i] + fromValueRaiser
		// }

	}
	log.Infoln("Value raised in predicted by ", valueRaisedPercentage, "%")
}
