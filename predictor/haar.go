package predictor

import (
	"fmt"
	"math"
)

var debug = false

// Haar ...  pastArray , scale is , logic numbers
func Haar(pastArray []float32, scale int, logic int) [][]float32 {

	/// Map for the scale and corresponding min,max. This will be used for reconstruction
	scaleMap := make(map[int][]float32)

	var ApproximateCoefficient []float32 = pastArray[:]
	var res [][]float32
	scaleNum := 1

	for len(ApproximateCoefficient) > 1 && scaleNum <= scale {
		var DetailedCoefficent []float32
		ApproximateCoefficient, DetailedCoefficent = haar_level(ApproximateCoefficient)
		if debug {
			fmt.Print("\n Haar Level ", scaleNum, "\t | \t", ApproximateCoefficient, "\t | \t", DetailedCoefficent, "\n")
		}

		// Transform ACoefficient and DCoefficient matrix to timeseries  A and D matrix to

		//A, D = convertHaarCoeeficientToTimeSeries(f, ACoefficient, DCoefficient)
		numStates := 12
		stateDeciderD := make([]float32, numStates+1)

		/// Find min, max
		minD, maxD := findMinMax(DetailedCoefficent)
		diffD := (maxD - minD) / float32(numStates)

		/// Add to map
		scaleMap[scaleNum] = []float32{minD, maxD, diffD}

		/// Calculate endpoints of the intervals which will determine states
		/// e.g state decider array [a,b,c,d]: State 1 is [a,b), state 2 is [b,c)

		stateDeciderD[0] = minD
		stateDeciderD[numStates] = maxD
		for i := 1; i < numStates; i++ {
			stateDeciderD[i] = stateDeciderD[i-1] + diffD
		}

		if debug {
			fmt.Println("State decider is: ", stateDeciderD)
		}

		/// Original array D transformed to state array containing the state numbers
		stateArrayD := make([]int, len(DetailedCoefficent))

		for index, elem := range DetailedCoefficent {
			stateArrayD[index] = findState(elem, stateDeciderD)
		}

		if debug {
			fmt.Println("State array of D is: ", stateArrayD)
		}

		/// Predicted States
		predictedArray := predictMarkov(stateArrayD, numStates, int(math.Pow(float64(2), float64(scale-scaleNum))))
		if debug {
			fmt.Println("Predicted array of D is ", predictedArray)
		}

		/// Convert States back to Avg between the gaps
		predictedActualD := convertStatesToValues(predictedArray, scaleMap[scaleNum])
		if debug {
			fmt.Println("Predicted Actual array of D is ", predictedActualD)
		}

		res = append([][]float32{predictedActualD}, res...) // prepend

		//res = append([][]float32{D}, res...) // prepend

		if scaleNum == scale {
			if debug {
				fmt.Println("ScaleNum == scale ? ", scale, scaleNum)
			}
			stateDeciderA := make([]float32, numStates+1)
			stateArrayA := make([]int, len(ApproximateCoefficient))
			minA, maxA := findMinMax(ApproximateCoefficient)
			if debug {
				fmt.Println("Min max A ", minA, ",", maxA)
			}
			diffA := (maxA - minA) / float32(numStates)
			// index = 0

			stateDeciderA[0] = minA
			stateDeciderA[numStates] = maxA
			for i := 1; i < numStates; i++ {
				stateDeciderA[i] = stateDeciderA[i-1] + diffA
			}
			for index, elem := range ApproximateCoefficient {
				stateArrayA[index] = findState(elem, stateDeciderA)
			}
			if debug {
				fmt.Println("State decider of A is: ", stateDeciderA)
				fmt.Println("State array of A in the last scale is: ", stateArrayA)
			}

			/// Predicted States
			predictedArrayA := predictMarkov(stateArrayA, numStates, 0)
			//float64(scale-scaleNum))))
			if debug {
				fmt.Println("Predicted array of A is ", predictedArrayA)
			}

			// Min Max A
			scaleMapA := []float32{minA, maxA, diffA}

			/// Convert States back to Avg between the gaps
			predictedActualA := convertStatesToValues(predictedArrayA, scaleMapA)
			if debug {
				fmt.Println("Predicted Actual array of A is ", predictedActualA)
			}

			res = append([][]float32{predictedActualA}, res...) // prepend

		}

		scaleNum++
	}
	// res = append([][]float32{A}, res...) // prepend

	return res
}

//To find haar Approximate and detailed coefficients for each level
func haar_level(f []float32) (a []float32, d []float32) {
	base := float32(1.0 / math.Pow(2, 0.5))
	var n int = len(f) / 2
	for i := 0; i < n; i++ {
		a = append(a, base*(f[2*i]+f[2*i+1]))
		d = append(d, base*(f[2*i]-f[2*i+1]))
	}
	return
}

// Haar ...
func Haar_old(f []float32, scale int) [][]float32 {

	/// Map for the scale and corresponding min,max. This will be used for reconstruction
	scaleMap := make(map[int][]float32)

	n := int(math.Log2(float64(len(f))))
	fmt.Println("n ", n)
	m := int(math.Pow(2, float64(n)))
	var A []float32 = f[:m]
	var res [][]float32
	scaleNum := 1

	for len(A) > 1 && scaleNum <= scale {
		var D []float32
		A, D = haar_level(A)
		fmt.Print("\n Haar Level - ", scaleNum, "---", A, "--", D, "\n")

		// Transform ACoefficient and DCoefficient matrix to timeseries  A and D matrix to

		//A, D = convertHaarCoeeficientToTimeSeries(f, ACoefficient, DCoefficient)
		numStates := 10
		stateDeciderD := make([]float32, numStates+1)

		/// Find min, max
		minD, maxD := findMinMax(D)
		diffD := (maxD - minD) / float32(numStates)

		/// Add to map
		scaleMap[scaleNum] = []float32{minD, maxD, diffD}

		/// Calculate endpoints of the intervals which will determine states
		/// e.g state decider array [a,b,c,d]: State 1 is [a,b), state 2 is [b,c)

		stateDeciderD[0] = minD
		stateDeciderD[numStates] = maxD
		for i := 1; i < numStates; i++ {
			stateDeciderD[i] = stateDeciderD[i-1] + diffD
		}

		fmt.Println("State decider is: ", stateDeciderD)

		/// Original array D transformed to state array containing the state numbers
		stateArrayD := make([]int, len(D))

		for index, elem := range D {
			stateArrayD[index] = findState(elem, stateDeciderD)
		}

		fmt.Println("State array of D is: ", stateArrayD)

		/// Predicted States
		predictedArray := predictMarkov(stateArrayD, numStates, int(math.Pow(float64(2), float64(scale-scaleNum))))
		fmt.Println("Predicted array of D is ", predictedArray)

		/// Convert States back to Avg between the gaps
		predictedActualD := convertStatesToValues(predictedArray, scaleMap[scaleNum])
		fmt.Println("Predicted Actual array of D is ", predictedActualD)

		res = append([][]float32{predictedActualD}, res...) // prepend

		//res = append([][]float32{D}, res...) // prepend

		if scaleNum == scale {
			fmt.Println("ScaleNum == scale ? ", scale, scaleNum)
			stateDeciderA := make([]float32, numStates+1)
			stateArrayA := make([]int, len(A))
			minA, maxA := findMinMax(A)
			fmt.Println("Min max A ", minA, ",", maxA)
			diffA := (maxA - minA) / float32(numStates)
			// index = 0

			stateDeciderA[0] = minA
			stateDeciderA[numStates] = maxA
			for i := 1; i < numStates; i++ {
				stateDeciderA[i] = stateDeciderA[i-1] + diffA
			}
			for index, elem := range A {
				stateArrayA[index] = findState(elem, stateDeciderA)
			}
			fmt.Println("State decider of A is: ", stateDeciderA)
			fmt.Println("State array of A in the last scale is: ", stateArrayA)

			/// Predicted States
			predictedArrayA := predictMarkov(stateArrayA, numStates, 0)
			//float64(scale-scaleNum))))
			fmt.Println("Predicted array of A is ", predictedArrayA)

			// Min Max A
			scaleMapA := []float32{minA, maxA, diffA}

			/// Convert States back to Avg between the gaps
			predictedActualA := convertStatesToValues(predictedArrayA, scaleMapA)
			fmt.Println("Predicted Actual array of A is ", predictedActualA)

			res = append([][]float32{predictedActualA}, res...) // prepend

		}

		scaleNum++
	}
	// res = append([][]float32{A}, res...) // prepend

	return res
}

// stateArrayGeneratorEquidistantWidth()

/*
// Converts coefficients to timeseries
func convertHaarCoeeficientToTimeSeries(pastArray, ACoefficient, DCoefficient []float32) (A, D []float32) {

	numOfCoefficients := len(D)
	numOfTimeSeriesPoints := len(pastArray)

	return
}
*/
// Converts States to values by converting int array to float32 array taking average between gaps
func convertStatesToValues(predictedArray []int, minMaxDiff []float32) []float32 {
	convertedArray := make([]float32, len(predictedArray))
	if len(minMaxDiff) != 3 {
		fmt.Println("Warn: Looks some issue in convertStatesToValues")
	}
	low := minMaxDiff[0]
	//max := minMaxDiff[1]
	diff := minMaxDiff[2]

	for i, el := range predictedArray {
		localLow := low + (float32(el-1) * diff)
		convertedArray[i] = float32(localLow + (diff / 2))
	}
	return convertedArray
}

/// Find the state by doing binary search on state decider
/// e.g array is [c,d,a,b,e]. If the element lies in [a,b), the state is 3. The last interval contains e inclusive
func findState(elem float32, stateDecider []float32) int {
	low := 0
	high := len(stateDecider) - 1

	for low < high {
		if elem == stateDecider[high] {
			return high
		}

		if elem == stateDecider[low] {
			return low + 1
		}

		median := (low + high) / 2
		if stateDecider[median] == elem {
			return median + 1
		}
		if stateDecider[low] < elem && low < high-1 && elem < stateDecider[low+1] {
			return low + 1
		}
		if high >= 1 && stateDecider[high-1] < elem && elem < stateDecider[high] {
			return high
		}
		if stateDecider[median] < elem && median < high-1 && elem < stateDecider[median+1] {
			return median + 1
		}
		if median >= 1 && stateDecider[median-1] < elem && elem < stateDecider[median] {
			return median
		}

		if stateDecider[median] < elem {
			low = median + 1
		} else {
			high = median - 1
		}

	}

	return low + 1
}

func findMinMax(arr []float32) (float32, float32) {
	min := float32(10000)
	max := float32(0)
	for i := 0; i < len(arr); i++ {
		if arr[i] < min {
			min = arr[i]
		}
		if arr[i] > max {
			max = arr[i]
		}
	}
	return min, max
}

func inverse_haar_level(a []float32, d []float32) (res []float32) {
	base := float32(1.0 / math.Pow(2, 0.5))
	for i := 0; i < len(a); i++ {
		res = append(res, base*(a[i]+d[i]))
		res = append(res, base*(a[i]-d[i]))
	}
	return
}

func Inverse_haar(h [][]float32) (an []float32) {
	an = h[0]
	for i := 1; i < len(h); i++ {
		an = inverse_haar_level(an, h[i])
	}
	return
}
