package predictor

import (
	"errors"
	"math"

	"github.com/Sirupsen/logrus"
)

///AccuracyChecker ... Function to check accuracy +- accuracy  Threshold
func AccuracyChecker(actualArray []float32, predictedArray []float32, size int, accuracyThreshold float32, log *logrus.Logger) (withingThresholdEstimatePercentage float32, underThresholdEstimatePercentage float32, overThresholdEstimatePercentage float32, rmseOverEstimate float32, rmseUnderEstimate float32, e error) {

	//actualArray := []float32{13, 12, 11.9, 13, 13, 13, 13, 11, 13, 12, 11.9, 12, 13, 13, 12, 12, 12, 13, 12.9, 12, 13, 11, 13, 13, 13, 12, 11.9, 13, 13, 12, 13, 12, 13, 13, 13, 13, 13, 12.9, 13, 13, 13, 12, 14, 13, 13, 11.9, 12, 13, 13, 13, 13, 12, 11.9, 12, 12, 13, 12, 12, 13, 12, 9, 10.9, 12, 13}

	var withingThresholdEstimateCount, underThresholdEstimateCount, overThresholdEstimateCount int
	var underEstimateSum, overEstimateSum float64

	if size != len(predictedArray) || size > len(actualArray) {
		return 0, 0, 0, 0, 0, errors.New("Len of predicted Array or actualArray is not same as prediction window size")
	}

	if len(actualArray) < size {
		log.Debugln("Trimming actual array in accuracy checker since D is less ")
		actualArray = actualArray[:size]

	}

	log.Debugln("Accuracy set as +-", accuracyThreshold)
	log.Debug("\nActual Array ", actualArray, "\n")

	log.Debug("\n---------------------- ", "\n")
	log.Debug("\nNon Matching elements ", "\n")
	log.Debug("\n---------------------- ", "\n")
	log.Debug("i\tActual\tPredicted", "\n")

	for i, predictedValue := range predictedArray {

		if predictedValue <= (actualArray[i]+accuracyThreshold) && predictedValue >= (actualArray[i]-accuracyThreshold) {
			// withing Threshold Estimate Count
			withingThresholdEstimateCount++

		} else if predictedValue < (actualArray[i] - accuracyThreshold) {
			//under Threshold Estimate Count
			underThresholdEstimateCount++
			underEstimateSum += math.Pow(float64(actualArray[i]-predictedValue), 2) // RMSE

		} else if predictedValue > (actualArray[i] + accuracyThreshold) {
			//over Threshold Estimate Count
			overThresholdEstimateCount++
			overEstimateSum += math.Pow(float64(predictedValue-actualArray[i]), 2) //RMSE

		}
		log.Debug("- ", i, "\t", actualArray[i], "\t", predictedValue, "")
	}

	if overThresholdEstimateCount == 0 && underThresholdEstimateCount == 0 {
		rmseOverEstimate = 0
		rmseUnderEstimate = 0

	} else if overThresholdEstimateCount == 0 && underThresholdEstimateCount != 0 {
		rmseOverEstimate = 0
		rmseUnderEstimate = float32(math.Sqrt(underEstimateSum / float64(underThresholdEstimateCount)))

	} else if underThresholdEstimateCount == 0 && overThresholdEstimateCount != 0 {
		rmseOverEstimate = float32(math.Sqrt(overEstimateSum / float64(overThresholdEstimateCount)))
		rmseUnderEstimate = 0
	} else {
		rmseOverEstimate = float32(math.Sqrt(overEstimateSum / float64(overThresholdEstimateCount)))
		rmseUnderEstimate = float32(math.Sqrt(underEstimateSum / float64(underThresholdEstimateCount)))
	}

	withingThresholdEstimatePercentage = (float32(withingThresholdEstimateCount) / float32(size)) * 100
	underThresholdEstimatePercentage = (float32(underThresholdEstimateCount) / float32(size)) * 100
	overThresholdEstimatePercentage = (float32(overThresholdEstimateCount) / float32(size)) * 100
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
