package main

import (
	"errors"
	"testing"
	"time"

	"github.com/adhuri/Compel-Prediction/fetcher"
	"github.com/adhuri/Compel-Prediction/predictor"
)

type Results struct {
	name   string
	result []ResultForOneDuration
}

type ResultForOneDuration struct {
	startTime                      int64
	endTime                        int64
	withinThresholdEstimatePercent float32
	overThresholdEstimatePercent   float32
	underThresholdEstimatePercent  float32
	rmseOverThresholdEstimate      float32
	rmseUnderThresholdEstimate     float32
}

func TestAccuracyForPredictedData(t *testing.T) {

	TestIp := "10.10.3.183"
	TestContainer := "d122412887e4"

	predictionResults := Results{name: "CPU Haar P1"}

	numberOfSlidingWindows := 8
	slidingWindow := 1024
	slidingWindowDuration := time.Duration(slidingWindow)

	endTime := time.Now().Add(time.Minute * -1)
	startTime := endTime.Add(-1 * time.Second * slidingWindowDuration)
	log.Infoln("================ Accuracy Graph Data Generator ===============")
	log.Infoln("Running Prediction ", predictionResults.name, " for ", numberOfSlidingWindows, " times")
	for i := 0; i < numberOfSlidingWindows; i++ {
		log.Infoln("For Start time :", startTime, ", End Time :", endTime)

		res1, err := getResults(TestIp, TestContainer, "cpu", "cpu_haar_P1_goup", slidingWindow, startTime, endTime)
		if err != nil {
			t.Error(err)
		}

		predictionResults.result = append(predictionResults.result, res1)

		endTime = startTime

		startTime = endTime.Add(-1 * time.Second * slidingWindowDuration)

	}

	resultPrinter(predictionResults)

}

// metric = cpu predictedMetric = "cpu_haar_P1"
func getResults(ip string, containerName string, metric string, predictedMetric string, slidingWindow int, startTime time.Time, endTime time.Time) (r ResultForOneDuration, e error) {
	DataFetcher := fetcher.NewDataFetcher()
	res := ResultForOneDuration{}

	//GetMetricDataForAccuracy(agentIp string, containerId string, metricType string, startTimeStamp int64, endTimeStamp int64) ([]float32, error) {
	actualData, err := DataFetcher.GetMetricDataForAccuracy(ip, containerName, metric, startTime.Unix(), endTime.Unix())
	if err != nil {
		return res, errors.New("Unable to fetch data from GetMetricDataForAccuracy ")
	}

	//func (dataFetcher *DataFetcher) GetPredictedData(agentIP string, containerId string, metric string, startTimeStamp int64, endTimeStamp int64) ([]float32, error) {
	predictedData, err1 := DataFetcher.GetPredictedData(ip, containerName, predictedMetric, startTime.Unix(), endTime.Unix())
	if err1 != nil {
		return res, errors.New("Unable to fetch data from GetPredictedData ")
	}

	accuracyThreshold := float32(1)
	res.withinThresholdEstimatePercent, res.underThresholdEstimatePercent, res.overThresholdEstimatePercent, res.rmseOverThresholdEstimate, res.rmseUnderThresholdEstimate, err = predictor.AccuracyChecker(actualData, predictedData, len(actualData), accuracyThreshold, log)
	if err != nil {
		return res, errors.New("Accuracy checker failed " + err.Error())
	}
	res.startTime = startTime.Unix()
	res.endTime = endTime.Unix()
	return res, nil

}

func resultPrinter(r Results) {
	log.Infoln("=============Results for ", r.name, "===============")
	var startTimestampArray []int64
	var withinThresholdEstimatePercentArray []float32
	var overThresholdEstimatePercentArray []float32
	var underThresholdEstimatePercentArray []float32
	var rmseOverThresholdEstimateArray []float32
	var rmseUnderThresholdEstimateArray []float32

	for _, el := range r.result {
		startTimestampArray = append([]int64{el.startTime}, startTimestampArray...)
		withinThresholdEstimatePercentArray = append([]float32{el.withinThresholdEstimatePercent}, withinThresholdEstimatePercentArray...)
		overThresholdEstimatePercentArray = append([]float32{el.overThresholdEstimatePercent}, overThresholdEstimatePercentArray...)
		underThresholdEstimatePercentArray = append([]float32{el.underThresholdEstimatePercent}, underThresholdEstimatePercentArray...)
		rmseOverThresholdEstimateArray = append([]float32{el.rmseOverThresholdEstimate}, rmseOverThresholdEstimateArray...)
		rmseUnderThresholdEstimateArray = append([]float32{el.rmseUnderThresholdEstimate}, rmseUnderThresholdEstimateArray...)

	}
	log.Infoln("----------------------------------------------------------------------------")
	log.Infoln("startTimestampArray ", startTimestampArray)

	log.Infoln("withinThresholdEstimatePercentArray ", withinThresholdEstimatePercentArray)
	log.Infoln("overThresholdEstimatePercentArray", overThresholdEstimatePercentArray)
	log.Infoln("underThresholdEstimatePercentArray", underThresholdEstimatePercentArray)
	log.Infoln("rmseOverThresholdEstimateArray", rmseOverThresholdEstimateArray)
	log.Infoln("rmseUnderThresholdEstimateArray", rmseUnderThresholdEstimateArray, "\n\n")

}
