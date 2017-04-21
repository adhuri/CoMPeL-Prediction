package predictor

import (
	"math"

	"github.com/Sirupsen/logrus"
)

type HaarPredictionLogic interface {
	/*

	   P 1 : Logic with equidistance fixed bins
	   P 2 : Logic  with variable number of bins based on- Closed proximity and No empty bin
	*/

	GetLogicName() string
	ConvertValuesToStates()
	ConvertStatesToValues()
}

// P1 struct logic 1 - fixed bin size
type P1 struct {
	name         string
	numberOfBins int // Number of bins
}

// P1 struct logic 1 - fixed bin size - Go Up Strategy
type P1GoUp struct {
	name         string
	numberOfBins int // Number of bins
}

// P2 ... logic 2 no fixed bins
type P2 struct {
	name            string
	cpl             float32 // Close proximity length
	maxNumberOfBins int     //Number of max bins
}

// Haar ...  pastArray , bin is max numberStates to generate , logic numbers
func Haar(pastArray []float32, predictionWindow int, bin int, logic int, log *logrus.Logger) [][]float32 {
	log.Debugln("Max Bins configured for states", bin)

	var goUp = false
	switch logic {
	case 1:
		log.Debugln("P1 Logic chosen ")
		goUp = false
	case 2:
		log.Debugln("P1 Go Up Logic Chosen")
		goUp = true
	default:
		log.Debugln("No valid logic chosen - Default P1 Logic chosen ")
		goUp = false
	}

	scale := int(math.Log2(float64(predictionWindow)))

	// Logic decider

	//  Scale for prediction

	/// Map for the scale and corresponding min,max. This will be used for reconstruction
	scaleMap := make(map[int][]float32)

	ApproximateCoefficient := pastArray[:]

	// To store the results
	var res [][]float32
	scaleNum := 1

	// Generate either till ApproximateCoefficient has length 1 or if scaleNum is equal to the required number of scales
	// based on predictionWindow Size
	for len(ApproximateCoefficient) > 1 && scaleNum <= scale {
		var DetailedCoefficent []float32
		ApproximateCoefficient, DetailedCoefficent = haarLevel(ApproximateCoefficient)
		log.Debugln("\n Haar Level ", scaleNum, "\t -|- \t", ApproximateCoefficient, "\t -|- \t", DetailedCoefficent, "\n")

		stateDeciderD := make([]float32, bin+1)

		/// Find min, max
		minD, maxD := findMinMax(DetailedCoefficent)
		diffD := (maxD - minD) / float32(bin)

		/// Add to map
		scaleMap[scaleNum] = []float32{minD, maxD, diffD}

		/// Calculate endpoints of the intervals which will determine states
		/// e.g state decider array [a,b,c,d]: State 1 is [a,b), state 2 is [b,c)

		stateDeciderD[0] = minD
		stateDeciderD[bin] = maxD
		for i := 1; i < bin; i++ {
			stateDeciderD[i] = stateDeciderD[i-1] + diffD
		}

		log.Debugln("State decider is: ", stateDeciderD)

		/// Original array D transformed to state array containing the state numbers
		stateArrayD := make([]int, len(DetailedCoefficent))

		for index, elem := range DetailedCoefficent {
			stateArrayD[index] = findState(elem, stateDeciderD)
		}

		log.Debugln("State array of D is: ", stateArrayD)

		/// Predicted States
		predictedArray := predictMarkov(stateArrayD, bin, int(math.Pow(float64(2), float64(scale-scaleNum))), goUp, log)
		log.Debugln("Predicted array of D is ", predictedArray)

		/// Convert States back to Avg between the gaps
		predictedActualD := convertStatesToValues(predictedArray, scaleMap[scaleNum], log)
		log.Debugln("Predicted Actual array of D is ", predictedActualD)

		res = append([][]float32{predictedActualD}, res...) // prepend

		//res = append([][]float32{D}, res...) // prepend

		if scaleNum == scale {
			log.Debugln("ScaleNum == scale ? ", scale, scaleNum)
			stateDeciderA := make([]float32, bin+1)
			stateArrayA := make([]int, len(ApproximateCoefficient))
			minA, maxA := findMinMax(ApproximateCoefficient)
			log.Debugln("Min max A ", minA, ",", maxA)
			diffA := (maxA - minA) / float32(bin)
			// index = 0

			stateDeciderA[0] = minA
			stateDeciderA[bin] = maxA
			for i := 1; i < bin; i++ {
				stateDeciderA[i] = stateDeciderA[i-1] + diffA
			}
			for index, elem := range ApproximateCoefficient {
				stateArrayA[index] = findState(elem, stateDeciderA)
			}
			log.Debugln("State decider of A is: ", stateDeciderA)
			log.Debugln("State array of A in the last scale is: ", stateArrayA)

			/// Predicted States
			predictedArrayA := predictMarkov(stateArrayA, bin, 0, goUp, log)
			//float64(scale-scaleNum))))
			log.Debugln("Predicted array of A is ", predictedArrayA)

			// Min Max A
			scaleMapA := []float32{minA, maxA, diffA}

			/// Convert States back to Avg between the gaps
			predictedActualA := convertStatesToValues(predictedArrayA, scaleMapA, log)
			log.Debugln("Predicted Actual array of A is ", predictedActualA)

			res = append([][]float32{predictedActualA}, res...) // prepend

		}

		scaleNum++
	}
	// res = append([][]float32{A}, res...) // prepend

	return res
}

// Converts States to values by converting int array to float32 array taking average between gaps
func convertStatesToValues(predictedArray []int, minMaxDiff []float32, log *logrus.Logger) []float32 {
	convertedArray := make([]float32, len(predictedArray))
	if len(minMaxDiff) != 3 {
		log.Errorln("Warn: Looks some issue in convertStatesToValues")
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

// Haar coefficients
//To find haar Approximate and detailed coefficients for each level
func haarLevel(f []float32) (a []float32, d []float32) {
	base := float32(1.0 / math.Pow(2, 0.5))
	var n int = len(f) / 2
	for i := 0; i < n; i++ {
		a = append(a, base*(f[2*i]+f[2*i+1]))
		d = append(d, base*(f[2*i]-f[2*i+1]))
	}
	return
}

//Inverse Haar coefficients

func inverseHaarLevel(a []float32, d []float32) (res []float32) {
	base := float32(1.0 / math.Pow(2, 0.5))
	for i := 0; i < len(a); i++ {
		res = append(res, base*(a[i]+d[i]))
		res = append(res, base*(a[i]-d[i]))
	}
	return
}

func InverseHaar(h [][]float32) (an []float32) {
	an = h[0]
	for i := 1; i < len(h); i++ {
		an = inverseHaarLevel(an, h[i])
	}
	return
}
