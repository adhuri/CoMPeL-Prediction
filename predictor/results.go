package predictor

import (
	"errors"
	"fmt"
)

///Function to check accuracy +- accuracy  Threshold
func AccuracyChecker(actualArray []float32, predictedArray []float32, size int, accuracyThreshold float32) (withingThresholdEstimatePercentage float32, underThresholdEstimatePercentage float32, overThresholdEstimatePercentage float32, averageOverEstimate float32, averageUnderEstimate float32, e error) {

	//actualArray := []float32{13, 12, 11.9, 13, 13, 13, 13, 11, 13, 12, 11.9, 12, 13, 13, 12, 12, 12, 13, 12.9, 12, 13, 11, 13, 13, 13, 12, 11.9, 13, 13, 12, 13, 12, 13, 13, 13, 13, 13, 12.9, 13, 13, 13, 12, 14, 13, 13, 11.9, 12, 13, 13, 13, 13, 12, 11.9, 12, 12, 13, 12, 12, 13, 12, 9, 10.9, 12, 13}

	var withingThresholdEstimateCount, underThresholdEstimateCount, overThresholdEstimateCount int
	var underEstimateSum, overEstimateSum float32

	if size != len(predictedArray) || size > len(actualArray) {
		return 0, 0, 0, 0, 0, errors.New("Len of predicted Array or actualArray is not same as prediction window size")
	}

	if len(actualArray) < size {
		fmt.Println("Trimming actual array in accuracy checker since D is less ")
		actualArray = actualArray[:size]

	}
	if debug {
		fmt.Println("Accuracy set as +-", accuracyThreshold)
		fmt.Print("\nActual Array ", actualArray, "\n")

		fmt.Print("\n---------------------- ", "\n")
		fmt.Print("\nNon Matching elements ", "\n")
		fmt.Print("\n---------------------- ", "\n")
		fmt.Print("\n i\tActual\tPredicted", "\n")
	}

	for i, predictedValue := range predictedArray {

		if predictedValue <= (actualArray[i]+accuracyThreshold) && predictedValue >= (actualArray[i]-accuracyThreshold) {
			// withing Threshold Estimate Count
			withingThresholdEstimateCount++

		} else if predictedValue < (actualArray[i] - accuracyThreshold) {
			//under Threshold Estimate Count
			underThresholdEstimateCount++
			underEstimateSum = actualArray[i] - predictedValue
			if debug {
				fmt.Print("- ", i, "\t", actualArray[i], "\t", predictedValue, "\n")
			}

		} else if predictedValue > (actualArray[i] + accuracyThreshold) {
			//over Threshold Estimate Count
			overThresholdEstimateCount++
			overEstimateSum = predictedValue - actualArray[i]
			if debug {
				fmt.Print("+ ", i, "\t", actualArray[i], "\t", predictedValue, "\n")
			}
		}
	}
	averageOverEstimate = overEstimateSum / float32(overThresholdEstimateCount)
	averageUnderEstimate = underEstimateSum / float32(underThresholdEstimateCount)
	withingThresholdEstimatePercentage = (float32(withingThresholdEstimateCount) / float32(size)) * 100
	underThresholdEstimatePercentage = (float32(underThresholdEstimateCount) / float32(size)) * 100
	overThresholdEstimatePercentage = (float32(overThresholdEstimateCount) / float32(size)) * 100
	return

}

//Prediction is not accurate and for near zero values could predict negative Values. Fixing them to zero
// All it means is the value approaches zero
func negativeValuesFixer(result []float32) {
	fixedCount := 0
	for i, el := range result {
		if el < 0 {
			result[i] = 0
			fixedCount += 1
		}
	}
	if debug {
		fmt.Println("Fixed negative values in predicted array ", fixedCount)
	}
}

// To create a variation for under prediction
func valueRaiser(result []float32, valueRaisedPercentage float32) {
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
	if debug {
		fmt.Println("Value raised in predicted by ", valueRaisedPercentage, "%")
	}
}
