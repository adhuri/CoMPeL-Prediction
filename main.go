package main

import (
	"flag"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Migration/protocol"
	"github.com/adhuri/Compel-Prediction/fetcher"
)

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

	//Init logging
	predictionFrequencyFlag := flag.Int64("pf", 20, "Prediction frequency in seconds")
	SlidingWindowSizeFlag := flag.Int("slidingwindow", 2048, "Sliding window in seconds- sampling per second")

	PredictionWindowSizeFlag := flag.Int("predictionwindow", 128, "Prediction window in seconds- sampling per second")
	flag.Parse()

	log.WithFields(logrus.Fields{
		"pf":               *predictionFrequencyFlag,
		"slidingwindow":    *SlidingWindowSizeFlag,
		"predictionwindow": *PredictionWindowSizeFlag,
	}).Infoln("Inputs from command line")

	SlidingWindowSize := *SlidingWindowSizeFlag                                  //SlidingWindow
	PredictionWindowSize := *PredictionWindowSizeFlag                            //PredictionWindow
	predictionFrequency := time.Second * time.Duration(*predictionFrequencyFlag) // in seconds
	log.Infoln("Prediction Frequency is ", predictionFrequency)
	DataFetcher := fetcher.NewDataFetcher()
	var timestamp int64
	var predictedValues []float32
	predictionTimer := time.NewTicker(predictionFrequency).C
	for {
		select {
		case <-predictionTimer:
			{
				predictionTime := time.Now()
				log.Infoln("Predicting for time : ", predictionTime.Format("2006-01-02 15:04:05"), " , Unix Time : ", predictionTime.Unix())

				ContainerInfo, err := DataFetcher.GetAgentInformation("127.0.0.1", "9091")
				if err != nil {
					log.Errorln("Error GetAgentInformation did not return client information ", err)
					continue
				}
				log.Infoln("=> Number of Agents for prediction ", len(ContainerInfo.Clients), "\n")
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

				sendDataToMigration(messageToSendToMigration, log)
			}

		}
	}

}
