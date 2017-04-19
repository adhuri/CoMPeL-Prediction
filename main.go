package main

import (
	"fmt"
	"time"

	"github.com/adhuri/Compel-Migration/protocol"
	"github.com/adhuri/Compel-Prediction/fetcher"
	predictor "github.com/adhuri/Compel-Prediction/predictor"
)

var debug = false

func main() {

	SlidingWindowSize := 2048               //SlidingWindow
	PredictionWindowSize := 128             //PredictionWindow
	predictionFrequency := time.Second * 10 // in seconds
	fmt.Println("Prediction Frequency is ", predictionFrequency)
	DataFetcher := fetcher.NewDataFetcher()
	var timestamp int64
	var predictedValues []float32
	predictionTimer := time.NewTicker(predictionFrequency).C
	for {
		select {
		case <-predictionTimer:
			{

				fmt.Println("Predicting for time ", time.Now())

				ContainerInfo, err := DataFetcher.GetAgentInformation("127.0.0.1", "9091")
				if err != nil {
					fmt.Println("Error GetAgentInformation did not return client information", err)
					continue
				}
				fmt.Println("Number of Agents for prediction ", len(ContainerInfo.Clients))
				metrics := []string{"cpu", "memory"}

				agentPredictions := []protocol.ClientInfo{}

				for _, agent := range ContainerInfo.Clients {

					containerPredictions := []protocol.ContainerInfo{}

					for _, containerID := range agent.Containers {
						// For every agent Container
						cpuPredictions := []float32{}
						memoryPredictions := []float32{}
						for _, metric := range metrics {
							predictedValues, timestamp = PredictAndStore(DataFetcher, string(agent.ClientIp), string(containerID), metric, SlidingWindowSize, PredictionWindowSize)
							if metric == "cpu" {
								cpuPredictions = predictedValues

							} else if metric == "memory" {
								memoryPredictions = predictedValues
							}
						}
						containerInfo := protocol.NewContainerInfo(string(containerID), cpuPredictions, memoryPredictions)
						containerPredictions = append(containerPredictions, *containerInfo)
					}
					clientInfo := protocol.NewClientInfo(string(agent.ClientIp), containerPredictions)
					agentPredictions = append(agentPredictions, *clientInfo)
				}

				messageToSendToMigration := protocol.NewPredictionData(timestamp, agentPredictions)

				sendDataTOMigration(messageToSendToMigration, log)
			}

		}
	}

}

func PredictAndStore(DataFetcher *fetcher.DataFetcher, agentIP string, containerID string, metric string, SlidingWindowSize int, PredictionWindowSize int) ([]float32, int64) {
	fmt.Println("Predicting ", metric, " for Agent:Container ", agentIP, ":", containerID)

	predictors := []string{"haar", "haargoup", "max"}
	var predictedArray []float32
	var timestamp int64
	for _, predictor := range predictors {

		fmt.Println("For predictor ", predictor)
		//agentIp string, containerId string, metricType string, time int64, numberOfPoints int) returns fetched array and time int64
		timestamp = time.Now().Unix()
		fetchedArray, alignedTimestamp := DataFetcher.GetMetricDataForContainer(agentIP, containerID, metric, timestamp, SlidingWindowSize)
		if debug {
			fmt.Println("Fetched Array for metric", metric, "-", fetchedArray)
		}

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
		if debug {
			fmt.Println("\nPredicted Array by ", predictor, predictedArray)

		}
		//Pass predicted array to store to influx db

		fmt.Println("Storing ", predictor, " predicted array  back to db ")
		//SavePredictedData(agentIP string, containerId string, metric string, predictedValues []float32, startTimeStamp int64) {

		err1 := DataFetcher.SavePredictedData(agentIP, containerID, metric+predictor, predictedArray, alignedTimestamp)
		if err1 != nil {
			fmt.Println("ERROR: Could not store predicted data using SavePredictedData for predictor ", predictor)
		}
	}
	return predictedArray, timestamp
}

func haarPrediction(SlidingWindowSize int, PredictionWindowSize int, fetchedData []float32, logic int) (predictedArray []float32) {

	bin := 30
	// Logic to start prediction

	haar := predictor.WaveletTransform{SlidingWindow: SlidingWindowSize, PredictionWindow: PredictionWindowSize}

	// replace fetchedData with fetched Data from db

	predictedArray, err := predictor.Predictor(&haar, fetchedData, bin, logic)
	if err != nil {
		fmt.Println("Error received from Predictor ", err)
		panic("Exiting due to Predictor not working")
	}

	return
}

func maxPrediction(SlidingWindowSize int, PredictionWindowSize int, fetchedData []float32, logic int) (predictedArray []float32) {

	bin := 0
	// Logic to start prediction

	max := predictor.MaxPredict{SlidingWindow: SlidingWindowSize, PredictionWindow: PredictionWindowSize}

	// replace fetchedData with fetched Data from db

	predictedArray, err := predictor.Predictor(&max, fetchedData, bin, logic)
	if err != nil {
		fmt.Println("Error received from Predictor ", err)
		panic("Exiting due to Predictor not working")
	}

	return
}
