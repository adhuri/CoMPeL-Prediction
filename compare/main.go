package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Prediction/fetcher"
	"github.com/adhuri/Compel-Prediction/predictor"
)

type ResultForDuration struct {
	startTime                      int64
	endTime                        int64
	withinThresholdEstimatePercent float32
	overThresholdEstimatePercent   float32
	underThresholdEstimatePercent  float32
	rmseOverThresholdEstimate      float32
	rmseUnderThresholdEstimate     float32
}

var (
	log *logrus.Logger
)

func init() {

	log = logrus.New()

	// Output logging to stdout
	log.Out = os.Stdout

	// Only log the info severity or above.
	log.Level = logrus.InfoLevel
	// Microseconds level logging
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	customFormatter.FullTimestamp = true

	log.Formatter = customFormatter

}

func main() {

	DataFetcher := fetcher.NewDataFetcher()
	startTime := time.Date(2017, time.April, 25, 23, 00, 0, 0, time.Local)
	endTime := time.Date(2017, time.April, 26, 1, 37, 0, 0, time.Local)

	ip := "152.1.13.223"
	metric := "cpu"
	predictedMetric := "cpuhaargoup"
	containerName := "459e40b6fc5b"

	actualData, err := DataFetcher.GetMetricDataForAccuracy(ip, containerName, metric, startTime.Unix(), endTime.Unix())
	if err != nil {
		fmt.Println("Unable to fetch data from GetMetricDataForAccuracy ")
	}

	//func (dataFetcher *DataFetcher) GetPredictedData(agentIP string, containerId string, metric string, startTimeStamp int64, endTimeStamp int64) ([]float32, error) {
	predictedData, err1 := DataFetcher.GetPredictedData(ip, containerName, predictedMetric, startTime.Unix(), endTime.Unix())
	if err1 != nil {
		fmt.Println("Unable to fetch data from GetPredictedData ")

	}

	fmt.Println("Actual Data", len(actualData), actualData[0:1])

	fmt.Println("Predicted Data", len(predictedData), predictedData[0:1])

	accuracyThreshold := float32(1)

	res := ResultForDuration{}

	res.withinThresholdEstimatePercent, res.underThresholdEstimatePercent, res.overThresholdEstimatePercent, res.rmseOverThresholdEstimate, res.rmseUnderThresholdEstimate, err = predictor.AccuracyChecker(actualData, predictedData, len(actualData), accuracyThreshold, log)
	if err != nil {
		fmt.Println("Accuracy checker failed " + err.Error())
	}
	resultPrinter(res)
}

func resultPrinter(el ResultForDuration) {

	var startTimestampArray []int64
	var withinThresholdEstimatePercentArray []float32
	var overThresholdEstimatePercentArray []float32
	var underThresholdEstimatePercentArray []float32
	var rmseOverThresholdEstimateArray []float32
	var rmseUnderThresholdEstimateArray []float32

	startTimestampArray = append([]int64{el.startTime}, startTimestampArray...)
	withinThresholdEstimatePercentArray = append([]float32{el.withinThresholdEstimatePercent}, withinThresholdEstimatePercentArray...)
	overThresholdEstimatePercentArray = append([]float32{el.overThresholdEstimatePercent}, overThresholdEstimatePercentArray...)
	underThresholdEstimatePercentArray = append([]float32{el.underThresholdEstimatePercent}, underThresholdEstimatePercentArray...)
	rmseOverThresholdEstimateArray = append([]float32{el.rmseOverThresholdEstimate}, rmseOverThresholdEstimateArray...)
	rmseUnderThresholdEstimateArray = append([]float32{el.rmseUnderThresholdEstimate}, rmseUnderThresholdEstimateArray...)

	log.Infoln("----------------------------------------------------------------------------")
	log.Infoln("startTimestampArray ", startTimestampArray)

	log.Infoln("withinThresholdEstimatePercentArray ", withinThresholdEstimatePercentArray)
	log.Infoln("overThresholdEstimatePercentArray", overThresholdEstimatePercentArray)
	log.Infoln("underThresholdEstimatePercentArray", underThresholdEstimatePercentArray)
	log.Infoln("rmseOverThresholdEstimateArray", rmseOverThresholdEstimateArray)
	log.Infoln("rmseUnderThresholdEstimateArray", rmseUnderThresholdEstimateArray, "\n\n")

}
