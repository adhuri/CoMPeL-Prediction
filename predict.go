package main

import (
	"time"

	"github.com/adhuri/Compel-Prediction/fetcher"
	predictor "github.com/adhuri/Compel-Prediction/predictor"
)

func PredictAndStore(DataFetcher *fetcher.DataFetcher, agentIP string, containerID string, metric string, SlidingWindowSize int, PredictionWindowSize int) ([]float32, int64) {
	log.Debugln("Predicting ", metric, " for Agent:Container ", agentIP, ":", containerID)

	predictors := []string{"haar", "haargoup", "max"}
	var predictedArray []float32
	var timestamp int64
	for _, predictor := range predictors {

		log.Debugln("For predictor ", predictor)
		//agentIp string, containerId string, metricType string, time int64, numberOfPoints int) returns fetched array and time int64
		timestamp = time.Now().Unix()
		fetchedArray, alignedTimestamp := DataFetcher.GetMetricDataForContainer(agentIP, containerID, metric, timestamp, SlidingWindowSize)
		log.Debugln("Fetched Array for metric", metric, "-", fetchedArray)

		predictedArray = []float32{}
		if predictor == "haar" {
			// Perform prediction
			predictedArray = haarPrediction(SlidingWindowSize, PredictionWindowSize, fetchedArray, 1)

		} else if predictor == "haargoup" {
			predictedArray = haarPrediction(SlidingWindowSize, PredictionWindowSize, fetchedArray, 2)

		} else if predictor == "max" {
			// Perform prediction
			predictedArray = maxPrediction(SlidingWindowSize, PredictionWindowSize, fetchedArray, 1)

		}

		log.Debugln("\nPredicted Array by ", predictor, predictedArray)

		//Pass predicted array to store to influx db

		log.Debugln("Storing ", predictor, " predicted array  back to db ")
		//SavePredictedData(agentIP string, containerId string, metric string, predictedValues []float32, startTimeStamp int64) {

		err1 := DataFetcher.SavePredictedData(agentIP, containerID, metric+predictor, predictedArray, alignedTimestamp)
		if err1 != nil {
			log.Errorln("ERROR: Could not store predicted data using SavePredictedData for predictor ", predictor)
		}
	}
	return predictedArray, timestamp
}

func haarPrediction(SlidingWindowSize int, PredictionWindowSize int, fetchedData []float32, logic int) (predictedArray []float32) {

	bin := 30
	// Logic to start prediction

	haar := predictor.WaveletTransform{SlidingWindow: SlidingWindowSize, PredictionWindow: PredictionWindowSize}

	// replace fetchedData with fetched Data from db

	predictedArray, err := predictor.Predictor(&haar, fetchedData, bin, logic, log)
	if err != nil {
		log.Errorln("Error received from Predictor ", err)
		panic("Exiting due to Predictor not working")
	}

	return
}

func maxPrediction(SlidingWindowSize int, PredictionWindowSize int, fetchedData []float32, logic int) (predictedArray []float32) {

	bin := 0
	// Logic to start prediction

	max := predictor.MaxPredict{SlidingWindow: SlidingWindowSize, PredictionWindow: PredictionWindowSize}

	// replace fetchedData with fetched Data from db

	predictedArray, err := predictor.Predictor(&max, fetchedData, bin, logic, log)
	if err != nil {
		log.Errorln("Error received from Predictor ", err)
		panic("Exiting due to Predictor not working")
	}

	return
}
