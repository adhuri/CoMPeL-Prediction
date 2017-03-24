package predictor

import (
	"fmt"
	"math/rand"
	"time"
)

func predictMarkov(arr []int, numStates int, predictionWindow int) (res []int) {

	// states := 10

	//matrixSize := states+1

	//x := []int{1,2,3,1,2,1,2,1,2,1,3,3,3,3,3,1}
	// x := []int{6, 5, 4, 4, 4, 3, 4, 4, 3, 4, 3, 5, 3, 4, 4, 5, 3, 3, 5, 4, 3, 3, 5, 7, 4, 5, 5, 4}

	n := len(arr)

	if debug {
		fmt.Print(arr, "\n")
		fmt.Print(n, "\n")
	}

	var p [][]float32
	for i := 0; i <= numStates; i++ {
		p = append(p, make([]float32, numStates+1))
	}

	for t := 0; t < n-1; t++ {
		f := arr[t]
		s := arr[t+1]
		p[f][s] = p[f][s] + 1
	}

	for i := 1; i <= numStates; i++ {
		sum := sumOfAllElements(p[i])
		for j := 1; j <= numStates; j++ {
			if sum == 0 {
				p[i][j] = 0
			} else {
				p[i][j] = p[i][j] / sum
			}
		}

	}
	var results []int // Declaring

	//Initializing
	if predictionWindow == 0 {
		results = make([]int, 1)
		lastElement := arr[len(arr)-1]
		results[0] = predictNext(lastElement, p, 0)

	} else {
		results = make([]int, predictionWindow)
	}

	lastElement := arr[len(arr)-1]
	for i := 0; i < predictionWindow; i++ {
		results[i] = predictNext(lastElement, p, predictionWindow)

		lastElement = results[i]
	}
	// fmt.Print(p, "\n")
	return results

}

func predictNext(lastElement int, transitionMatrix [][]float32, D int) int {

	var maxIndices = make([]int, D)
	max := transitionMatrix[lastElement][1] // since n+1 by n+1 matrix

	// 1st pass for max
	for _, elem := range transitionMatrix[lastElement] {
		if elem > max {
			max = elem
		}
	}
	/// Special case all zeros
	if D == 0 || max == 0 {
		//fmt.Print("\n-~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~``---D is possibily 0 ----- \n")
		return random(1, len(transitionMatrix[1]))
	}

	// 2nd pass for max

	for i, elem := range transitionMatrix[lastElement] {
		if elem == max {
			maxIndices = append(maxIndices, i)

		}
	}

	return maxIndices[random(0, len(maxIndices))]

}

// func findMaxProbabilityState(int rowNumber, transitionMatrix [][]float32) int {
//
// }

func random(min, max int) int {

	if !debug {
		//Use rand.Seed() when in production for vvariable output
		rand.Seed(time.Now().Unix())
	}

	return rand.Intn(max-min) + min
}

func sumOfAllElements(array []float32) (sum float32) {
	for _, i := range array {
		sum += i
	}
	//fmt.Print("Sum",sum)
	return
}
