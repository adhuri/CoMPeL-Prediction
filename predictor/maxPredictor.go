package predictor

// Max values

func (mp *MaxPredict) GetPredictorName() string {
	return "Max Predictor "
}

//Used to Predict WaveletTransform
func (mp *MaxPredict) Predict(pastArray []float32, bin int, logic int) (predictedArray []float32, err error) {
	_, max := findMinMax(pastArray)
	for i := 0; i < mp.PredictionWindow; i++ {
		predictedArray = append(predictedArray, max)
	}

	return
}
