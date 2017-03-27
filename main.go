package main

import (
	"fmt"
	"time"

	fetcher "github.com/adhuri/Compel-Prediction/fetcher"
	predictor "github.com/adhuri/Compel-Prediction/predictor"
)

var debug = true

func main() {

	SlidingWindowSize := 2048               //SlidingWindow
	PredictionWindowSize := 128             //PredictionWindow
	predictionFrequency := time.Second * 10 // in seconds
	fmt.Println("Prediction Frequency is ", predictionFrequency)

	// Data fetcher
	DataFetcher := fetcher.NewDataFetcher()

	predictionTimer := time.NewTicker(predictionFrequency).C
	for {
		select {
		case <-predictionTimer:
			{

				fmt.Println("Predicting for time ", predictionTimer)
				//Store Container_name to IP mapping from monitoring server
				err := storeAgentDetails()
				if err != nil {
					fmt.Println("Error: Unable to store Agent details ", err)
				}

				// Fetch data from db
				//fetchedData := []float32{26.7, 20, 14, 25, 10.9, 22, 22, 11, 7, 9, 12.9, 9, 13, 19, 19, 24, 32.7, 35, 27, 23, 15, 22, 37.6, 36, 41, 42, 42.6, 42, 40.6, 43, 39, 39.6, 40, 39, 42, 38.6, 38, 40.6, 45, 43.6, 43, 44, 37.6, 42, 41.6, 41, 39, 16, 19.8, 15, 23, 23, 16.8, 22, 22, 16, 16, 16, 14, 16, 26.7, 22, 25, 12, 14, 28, 13.9, 13, 13, 13, 14, 12, 14, 13, 10, 13, 11.9, 9, 10, 9, 11, 9, 10, 11, 12, 19.8, 22, 23, 22, 13, 13, 26.7, 14, 16, 14, 11, 19, 19.8, 13, 13, 12, 13, 12, 12, 12.9, 13, 13, 13, 24, 13, 14, 19.8, 23, 15, 26, 17, 27, 13, 14.9, 13, 13, 14, 11, 8, 6, 7.9, 8, 6, 5, 8, 12, 11, 11.9, 10, 11, 9, 10, 9, 11, 23, 10.9, 9, 12, 11, 15, 16, 29, 11.9, 11, 13, 18, 13, 12, 17, 30.7, 13, 11, 12, 28, 16, 12, 10.9, 12, 11, 14, 13, 13, 13, 14, 13, 13.9, 14, 12, 12, 12, 15, 14, 23, 15.8, 14, 27, 13, 20, 20, 10, 11, 10.9, 11, 11, 10, 11, 11, 11, 11.9, 11, 12, 17, 13, 14, 13, 13, 12, 11.9, 12, 12, 13, 11, 11, 12, 12, 10.9, 11, 12, 12, 13, 10, 12, 12.9, 10, 9, 11, 13, 13, 10, 11, 12.9, 13, 13, 11, 9, 12, 12, 14, 10.9, 12, 13, 13, 12, 13, 13, 13, 13.9, 12, 11, 12, 11, 11, 12, 9, 14, 12, 12, 14, 13, 11, 12, 10.9, 11, 10, 12, 10, 11, 9, 12, 9, 12, 11, 14, 12, 11, 13, 11, 11.9, 13, 12, 12, 9, 11, 10, 10, 9, 12.9, 10, 10, 10, 13, 12, 10, 13, 11.9, 12, 12, 11, 9, 10, 12, 12, 12, 12.9, 11, 12, 12, 12, 13, 13, 13, 12.9, 13, 13, 12, 12, 12, 11, 13, 12, 13.9, 12, 14, 12, 12, 12, 12, 12, 11.9, 11, 11, 10, 11, 12, 9, 10, 9, 10, 9, 11, 10, 11, 11, 13, 12.9, 13, 12, 13, 13, 13, 13, 12, 12, 10.9, 12, 13, 12, 12, 10, 12, 12, 13, 12, 10.9, 11, 13, 11, 13, 13, 10, 9, 9, 11, 12, 13, 13, 11, 12, 12, 12.9, 13, 12, 14, 12, 13, 13, 12, 12.9, 12, 13, 13, 11, 13, 13, 12, 12, 11.9, 13, 11, 12, 12, 12, 10, 10, 10.9, 10, 10, 12, 11, 10, 10, 11, 9, 11, 12, 9, 9, 13, 12, 11, 11.9, 14, 13, 13, 13, 10, 13, 13, 12, 13, 12.9, 13, 13, 12, 12, 12, 13, 11.9, 13, 11, 13, 12, 11, 13, 13, 11.9, 11, 12, 12, 13, 12, 13, 13, 13, 11.9, 10, 13, 10, 11, 11, 10, 11, 12, 12, 11.9, 12, 13, 14, 13, 12, 12, 12, 12, 9, 12, 10, 12, 10, 10, 9, 10, 11, 9, 9, 9, 9, 13, 12, 14, 12, 9, 10.9, 12, 13, 12, 12, 12, 12, 11, 11.9, 11, 11, 12, 13, 10, 9, 9, 9, 14, 13, 11, 10, 12, 10, 9, 10.9, 10, 9, 11, 11, 11, 9, 10, 11, 10.9, 9, 11, 10, 11, 11, 12, 11, 9, 11, 12, 13, 9, 11, 9, 10, 12, 12.9, 12, 12, 12, 12, 13, 12, 12, 12, 12.9, 13, 12, 12, 12, 12, 10, 10, 12, 11.9, 13, 10, 14, 13, 12, 12, 10, 10.9, 13, 11, 13, 9, 9, 12, 12, 11.9, 12, 11, 12, 12, 11, 9, 11, 11.9, 12, 11, 11, 12, 12, 13, 13, 12, 11.9, 12, 12, 12, 10, 11, 11, 12, 12, 12, 13.9, 11, 13, 11, 12, 11, 11, 12, 5, 5, 5, 5, 7, 11, 13, 12, 12.9, 12, 12, 12, 13, 13, 13, 12, 13, 10.9, 12, 12, 11, 11, 12, 13, 13, 11.9, 13, 13, 12, 12, 11, 11, 14, 12, 11.9, 12, 12, 12, 11, 12, 12, 12, 12, 12.9, 12, 11, 13, 11, 13, 12, 11.9, 12, 12, 11, 13, 12, 12, 13, 12.9, 13, 12, 13, 13, 13, 13, 11, 10, 11.9, 9, 12, 13, 12, 12, 12, 11.9, 12, 12, 12, 12, 11, 11, 13, 12.9, 13, 12, 11, 11, 12, 12, 9, 12, 11.9, 11, 12, 13, 10, 13, 11, 13, 13, 11.9, 12, 12, 12, 12, 11, 13, 12, 12.9, 13, 13, 11, 12, 10, 11, 11, 12, 11.9, 11, 11, 9, 13, 11, 11, 11, 11, 11.9, 12, 10, 9, 11, 13, 13, 13, 13, 10.9, 13, 11, 11, 10, 10, 10, 11, 10.9, 12, 13, 11, 13, 11, 10, 10, 11, 11.9, 11, 10, 12, 10, 14, 11, 13, 13, 12.9, 13, 10, 13, 10, 13, 13, 12, 10, 13.9, 10, 11, 11, 12, 10, 11, 11, 13, 10.9, 13, 12, 13, 13, 12, 12, 12, 12.9, 13, 13, 11, 12, 12, 13, 11.9, 12, 12, 13, 11, 10, 12, 13, 12, 11.9, 12, 13, 11, 11, 9, 11, 12, 12, 9, 11, 13, 11, 13, 12, 13, 12, 11.9, 13, 11, 12, 12, 12, 12, 12, 12.9, 13, 12, 11, 13, 13, 12, 12, 12, 12.9, 13, 9, 9, 12, 14, 12, 13, 11.9, 13, 12, 12, 13, 12, 12, 13, 13, 12.9, 10, 14, 12, 11, 12, 12, 11, 11.9, 12, 12, 12, 12, 13, 13, 12.9, 13, 12, 13, 12, 13, 12, 12, 13, 10.9, 11, 11, 11, 11, 10, 12, 12, 13, 12.9, 11, 9, 12, 12, 12, 12, 10, 12, 11.9, 12, 12, 12, 14, 13, 11, 12, 11, 10.9, 12, 12, 14, 12, 12, 13, 13, 11.9, 13, 12, 13, 13, 13, 11, 12.9, 12, 12, 13, 13, 12, 13, 11, 11.9, 11, 11, 13, 11, 13, 13, 12, 10, 12.9, 12, 11, 12, 12, 11, 12, 11, 11.9, 11, 11, 10, 12, 11, 11, 11, 11, 12, 11.9, 13, 13, 13, 9, 9, 13, 13, 11.9, 12, 12, 10, 10, 11, 11, 12, 13, 10.9, 12, 14, 11, 12, 11, 12, 13, 12.9, 13, 13, 13, 12, 11, 12, 12, 11.9, 14, 12, 13, 13, 13, 13, 10, 9, 12, 12, 12, 12, 9, 14, 12, 12.9, 11, 13, 12, 12, 11, 12, 13, 11.9, 13, 13, 12, 13, 10, 12, 11, 9, 11, 10, 14, 12, 11, 12, 12, 11.9, 13, 13, 13, 13, 13, 13, 13, 12, 12.9, 13, 11, 13, 12, 14, 13, 13, 9, 9, 11, 13, 12, 11, 10, 13, 12, 11.9, 12, 13, 12, 11, 11, 11, 14, 11.9, 13, 13, 12, 14, 12, 14, 13, 13, 12.9, 13, 12, 12, 11, 12, 12, 12, 10, 9, 12, 12, 11, 12, 12, 13, 13, 12.9, 13, 13, 13, 13, 12, 13, 12, 12.9, 12, 13, 13, 12, 13, 11, 12, 13, 9, 9, 12, 12, 12, 12, 12, 12, 12.9, 12, 11, 13, 11, 13, 12, 10.9, 11, 12, 12, 11, 10, 10, 12, 13, 10.9, 11, 12, 10, 12, 11, 12, 14, 12, 10.9, 13, 13, 14, 10, 13, 12, 12.9, 12, 13, 13, 14, 13, 13, 14, 13, 12.9, 12, 13, 11, 11, 12, 11, 11, 9, 11, 12, 11, 11, 12, 10, 10, 13, 9, 12, 13, 12, 13, 13, 12, 24.8, 29, 18, 23, 17, 20, 22, 17, 14.9, 14, 14, 14, 13, 13, 12, 13, 11.9, 12, 11, 13, 12, 8, 5, 5, 6.9, 5, 10, 12, 12, 12, 12, 11, 12, 12.9, 12, 10, 12, 10, 12, 11, 12, 12, 11, 11.9, 12, 12, 11, 11, 12, 12, 12, 11.9, 13, 11, 12, 11, 11, 11, 12, 12, 11.9, 11, 11, 12, 12, 11, 12, 11, 10.9, 12, 12, 12, 13, 11, 12, 12, 12, 10.9, 12, 12, 12, 11, 12, 11, 11, 11.9, 12, 11, 11, 12, 11, 11, 12, 11, 10.9, 12, 11, 11, 11, 12, 11, 12, 12, 11.9, 12, 12, 12, 10, 13, 12, 13, 11.9, 12, 11, 12, 13, 12, 12, 13, 11.9, 13, 12, 12, 12, 12, 11, 11, 11.9, 13, 11, 13, 12, 12, 11, 12, 12, 10.9, 12, 11, 12, 12, 9, 11, 12, 11.9, 13, 12, 24, 19, 23, 18.8, 14, 14, 14, 27, 20, 11.9, 13, 13, 11, 10, 16, 22, 20, 16, 28.7, 22, 19, 30, 21.8, 15, 16, 16, 12, 11, 14, 12.9, 13, 15, 11, 15, 13, 14, 13, 12.9, 9, 13, 12, 12, 12, 12, 12, 13, 10.9, 13, 12, 11, 11, 13, 11, 12, 11.9, 13, 12, 13, 13, 12, 11, 13, 10.9, 13, 13, 13, 14, 12, 13, 12, 12, 11.9, 12, 12, 13, 12, 12, 12, 12, 10.9, 12, 13, 12, 12, 12, 12, 11, 12.9, 12, 10, 11, 13, 13, 13, 11, 10.9, 12, 12, 12, 11, 11, 13, 13, 12.9, 13, 13, 13, 13, 11, 12, 12, 11, 12, 11.9, 9, 13, 11, 12, 11, 12, 12, 11, 12.9, 12, 14, 12, 13, 12, 11, 11.9, 11, 18, 20, 13, 12, 26, 18.8, 12, 13, 12, 12, 10.9, 11, 12, 13, 12, 12, 12, 10, 12, 11.9, 12, 13, 13, 12, 10, 13, 13, 13, 12.9, 11, 12, 14, 13, 13, 13, 13, 12, 11.9, 12, 13, 13, 12, 13, 12, 13, 13, 12.9, 12, 13, 12, 13, 12, 13, 12, 12, 12.9, 13, 13, 18, 12, 12, 12, 11, 12.9, 13, 13, 12, 9, 11, 13, 13, 12, 11.9, 14, 12, 13, 13, 13, 13, 12, 14, 12, 12.9, 12, 13, 13, 13, 12, 12, 12, 11.9, 12, 11, 13, 12, 12, 12, 12, 11.9, 12, 13, 12, 9, 9, 12, 13, 13, 12.9, 13, 12, 13, 13, 12, 12, 11.9, 13, 13, 13, 12, 12, 13, 12, 12, 11.9, 11, 12, 10, 11, 10, 10, 10, 10.9, 11, 11, 11, 11, 12, 13, 11, 13, 11.9, 13, 12, 12, 12, 13, 12, 12, 10.9, 13, 12, 12, 10, 10, 12, 12, 11.9, 12, 12, 13, 12, 11, 12, 13, 10.9, 9, 9, 13, 11, 12, 13, 12, 10.9, 13, 11, 13, 12, 12, 11, 12, 12, 11.9, 12, 11, 13, 12, 13, 13, 12, 12, 11.9, 13, 12, 13, 13, 13, 12, 20, 12, 11.9, 11, 30, 15, 11, 12, 10, 11, 11, 10.9, 11, 13, 11, 11, 11, 11, 12, 11.9, 13, 11, 13, 11, 12, 13, 9, 12, 12, 18.8, 26, 12, 11, 12, 12, 10.9, 9, 11, 12, 12, 12, 13, 11, 12, 23.8, 10, 12, 12, 29, 17, 12, 11, 12.9, 12, 10, 11, 12, 13, 11, 12, 12, 11.9, 13, 12, 12, 12, 13, 12, 12, 12.9, 14, 14, 13, 12, 13, 10, 9, 12, 13, 12, 13, 11, 13, 12, 12, 11.9, 18, 12, 12, 13, 15, 10, 17.8, 12, 12, 13, 12, 13, 12, 12, 12, 11.9, 12, 12, 12, 11, 13, 12, 13, 13, 11.9, 12, 12, 11, 13, 12, 9, 13, 11.9, 12, 13, 11, 13, 12, 13, 12, 13, 11.9, 14, 12, 12, 11, 13, 12, 13, 11.9, 11, 12, 13, 20, 12, 12, 13, 12, 13.9, 8, 5, 5, 6, 7, 10, 10.9, 9, 12, 13, 11, 13, 12, 12, 12.9, 13, 10, 12, 12, 12, 13, 13, 11.9, 13, 12, 13, 11, 12, 13, 13, 11.9, 13, 12, 13, 13, 14, 14, 12, 9, 9, 11, 12, 13, 16, 13, 12, 11, 9, 11, 13, 12, 12, 13, 13, 12, 11.9, 11, 12, 13, 12, 12, 12, 13, 12, 11.9, 18, 10, 9, 9, 13, 15, 12, 16.8, 22, 12, 11, 14, 27, 12, 11.9, 12, 12, 12, 11, 13, 12, 10, 12.9, 12, 13, 12, 12, 11, 13, 12, 11.9, 11, 13, 12, 12, 12, 12, 13, 11.9, 13, 12, 11, 11, 12, 13, 11, 11.9, 12, 10, 9, 12, 12, 12, 16, 11.9, 9, 9, 12, 12, 13, 13, 11, 12, 12, 11.9, 12, 13, 12, 12, 12, 13, 10, 12, 11.9, 13, 11, 13, 12, 12, 13, 12, 12, 12.9, 12, 13, 12, 11, 13, 13, 12, 11, 9, 9, 12, 13, 13, 12, 12, 13, 11.9, 13, 12, 12, 13, 12, 13, 12, 12, 11.9, 14, 12, 12, 14, 12, 13, 12, 13, 11.9, 13, 13, 12, 12, 13, 13, 11.9, 13, 13, 13, 12, 12, 13, 12, 11.9, 19, 13, 12, 14, 13, 13, 12, 13, 11.9, 13, 11, 13, 13, 13, 13, 13, 13, 12, 12.9, 11, 11, 12, 12, 12, 13, 12.9, 13, 13, 12, 12, 12, 12, 13, 13, 12.9, 13, 12, 14, 12, 13, 13, 11.9, 9, 9, 13, 12, 10}

				//for every container for CPU
				{
					//agentIp string, containerId string, metricType string, time int64, numberOfPoints int) returns fetched array and time int64
					fetchedDataCPUArray, alignedTimestamp := DataFetcher.GetMetricDataForContainer("192.168.0.28", "mysql", "cpu", time.Now().Unix(), SlidingWindowSize)

					if debug {
						fmt.Println("Fetched Array ", fetchedDataCPUArray)
					}

					// Perform prediction
					predictedArrayCPUHaar := haarPrediction(SlidingWindowSize, PredictionWindowSize, fetchedDataCPUArray, 1)
					if debug {
						fmt.Println("\nPredicted Array - Haar P1", predictedArrayCPUHaar)

					}
					// Check Accuracy of prediciton

					//Pass predicted array to store to influx db

					fmt.Println("Storing Haar predicted array P1 back to db ")
					//SavePredictedData(agentIP string, containerId string, metric string, predictedValues []float32, startTimeStamp int64) {

					err1 := DataFetcher.SavePredictedData("192.168.0.28", "mysql", "cpu_haar_P1", predictedArrayCPUHaar, alignedTimestamp)
					if err1 != nil {
						fmt.Println("ERROR: Could not store predicted data using SavePredictedData")
					}

					// Using GO UP predictor

					// Perform prediction
					predictedArrayCPUHaarGoUp := haarPrediction(SlidingWindowSize, PredictionWindowSize, fetchedDataCPUArray, 2)
					if debug {
						fmt.Println("\nPredicted Array - Haar Go Up ", predictedArrayCPUHaarGoUp)

					}
					// Check Accuracy of prediciton

					//Pass predicted array to store to influx db

					fmt.Println("Storing P1 GoUp predicted array back to db ")
					//SavePredictedData(agentIP string, containerId string, metric string, predictedValues []float32, startTimeStamp int64) {

					err3 := DataFetcher.SavePredictedData("192.168.0.28", "mysql", "cpu_haar_P1_goup", predictedArrayCPUHaarGoUp, alignedTimestamp)
					if err3 != nil {
						fmt.Println("ERROR: Could not store predicted data using SavePredictedData")
					}

					// Using max predictor

					// Perform prediction
					predictedArrayCPUMax := maxPrediction(SlidingWindowSize, PredictionWindowSize, fetchedDataCPUArray, 1)
					if debug {
						fmt.Println("\nPredicted Array - Max ", predictedArrayCPUMax)

					}
					// Check Accuracy of prediciton

					//Pass predicted array to store to influx db

					fmt.Println("Storing Max predicted array back to db ")
					//SavePredictedData(agentIP string, containerId string, metric string, predictedValues []float32, startTimeStamp int64) {

					err2 := DataFetcher.SavePredictedData("192.168.0.28", "mysql", "cpu_max", predictedArrayCPUMax, alignedTimestamp)
					if err2 != nil {
						fmt.Println("ERROR: Could not store predicted data using SavePredictedData")
					}

				}

				// for every container for memory
				{
					//agentIp string, containerId string, metricType string, time int64, numberOfPoints int) returns fetched array and time int64
					fetchedDataMemoryArray, alignedTimestamp := DataFetcher.GetMetricDataForContainer("192.168.0.28", "mysql", "memory", time.Now().Unix(), SlidingWindowSize)

					if debug {
						fmt.Println("Fetched Array ", fetchedDataMemoryArray)
					}

					// Perform prediction
					predictedArrayMemoryHaar := haarPrediction(SlidingWindowSize, PredictionWindowSize, fetchedDataMemoryArray, 1)
					if debug {
						fmt.Println("\nPredicted Array ", predictedArrayMemoryHaar)

					}
					// Check Accuracy of prediciton

					//Pass predicted array to store to influx db

					fmt.Println("Storing predicted memory array back to db ")
					//SavePredictedData(agentIP string, containerId string, metric string, predictedValues []float32, startTimeStamp int64) {

					err1 := DataFetcher.SavePredictedData("192.168.0.28", "mysql", "memory_haar", predictedArrayMemoryHaar, alignedTimestamp)
					if err1 != nil {
						fmt.Println("ERROR: Could not store predicted data using SavePredictedData")
					}

					// Using P1 Go UP

					// Perform prediction
					predictedArrayMemoryHaarGoUp := haarPrediction(SlidingWindowSize, PredictionWindowSize, fetchedDataMemoryArray, 2)
					if debug {
						fmt.Println("\nPredicted Array - Haar P1 Go Up ", predictedArrayMemoryHaarGoUp)

					}
					// Check Accuracy of prediciton

					//Pass predicted array to store to influx db

					fmt.Println("Storing Haar Go Up predicted array back to db ")
					//SavePredictedData(agentIP string, containerId string, metric string, predictedValues []float32, startTimeStamp int64) {

					err3 := DataFetcher.SavePredictedData("192.168.0.28", "mysql", "memory_haar_goup", predictedArrayMemoryHaarGoUp, alignedTimestamp)
					if err3 != nil {
						fmt.Println("ERROR: Could not store predicted data using SavePredictedData")
					}

					// Using max predictor
					// Perform prediction
					predictedArrayMemoryMax := maxPrediction(SlidingWindowSize, PredictionWindowSize, fetchedDataMemoryArray, 1)
					if debug {
						fmt.Println("\nPredicted Array - Max ", predictedArrayMemoryMax)

					}
					// Check Accuracy of prediciton

					//Pass predicted array to store to influx db

					fmt.Println("Storing Max predicted array back to db ")
					//SavePredictedData(agentIP string, containerId string, metric string, predictedValues []float32, startTimeStamp int64) {

					err2 := DataFetcher.SavePredictedData("192.168.0.28", "mysql", "memory_max", predictedArrayMemoryMax, alignedTimestamp)
					if err2 != nil {
						fmt.Println("ERROR: Could not store predicted data using SavePredictedData")
					}

				}

				// Pass all predicted array to migration decider

			}
		}
	}

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

func storeAgentDetails() (err error) {
	// Get Agent Details from the fb
	//Get Agent details - map of container_name : Agent IP
	//fetcher.GetAgentDetails()
	return nil

}
