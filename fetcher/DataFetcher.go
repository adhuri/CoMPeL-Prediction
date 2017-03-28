package fetcher

import (
	"fmt"
	"sync"
)

type DataFetcher struct {
	sync.RWMutex
	dataCache Cache
}

func NewDataFetcher() *DataFetcher {
	return &DataFetcher{
		dataCache: Cache{},
	}
}

func (dataFetcher *DataFetcher) GetAgentInformation() {

}

func (dataFetcher *DataFetcher) GetMetricDataForContainer(agentIp string, containerId string, metricType string, time int64, numberOfPoints int) ([]float32, int64) {

	dataPoints := getData(agentIp, containerId, metricType)
	dataPointMap := make(map[int64]float32)

	var latestTimesStamp int64
	var oldestTimesStamp int64

	for i, point := range dataPoints {
		dataPointMap[point.Timestamp] = point.Value
		if i == 0 {
			oldestTimesStamp = point.Timestamp
		}
		if point.Timestamp > latestTimesStamp {
			latestTimesStamp = point.Timestamp
		}

		if point.Timestamp < oldestTimesStamp {
			oldestTimesStamp = point.Timestamp
		}
	}

	var points []float32
	for i := oldestTimesStamp; i <= time; i += 1 {
		// if there is break in time series, aligning will be impossible with 2 seconds sampling
		if value, present := dataPointMap[i]; present {
			points = append(points, value)
		} else {
			points = append(points, -1) //some point might have 0 value
		}
	}

	// for i, point := range points {
	// 	fmt.Printf(" %d : %f \n", i, point)
	// }

	fillMissingValues(points)

	// for i, point := range points {
	// 	fmt.Printf(" %d : %f \n", i, point)
	// }

	if len(points) > numberOfPoints {
		// Trim the slice if we have more points than asked for
		numberOfExtraPoints := len(points) - numberOfPoints
		//fmt.Println(len(points[numberOfExtraPoints:]))
		return points[numberOfExtraPoints:], time
	} else if len(points) < numberOfPoints {
		// If points are less than required then pad 0 at the start
		var remainingPoints []float32
		numberOfPointsMissing := numberOfPoints - len(points)
		for i := 0; i < numberOfPointsMissing; i += 1 {
			remainingPoints = append(remainingPoints, 0)
		}
		remainingPoints = append(remainingPoints, points...)

		//fmt.Println(len(remainingPoints))
		return remainingPoints, time

	}

	return points, time

}

func fillMissingValues(points []float32) {
	previousNonZeroIndex := 0
	flag := false

	for i, point := range points {
		if (point != -1) && (!flag) {
			previousNonZeroIndex = i
		} else if point == -1 {
			flag = true
			points[i] = float32(0)
		} else if (point != -1) && (flag) {
			if i-previousNonZeroIndex <= 40 {
				mean := (points[previousNonZeroIndex] + points[i]) / float32(2)
				for j := previousNonZeroIndex + 1; j < i; j++ {
					points[j] = mean
				}
			}
			previousNonZeroIndex = i
			flag = false
		}
	}

	if len(points)-previousNonZeroIndex <= 40 {
		for j := previousNonZeroIndex + 1; j < len(points); j++ {
			points[j] = points[previousNonZeroIndex]
		}
	}

}

func (dataFetcher *DataFetcher) GetMetricDataForAccuracy(agentIp string, containerId string, metricType string, startTimeStamp int64, endTimeStamp int64) ([]float32, error) {

	dataPoints := getData(agentIp, containerId, metricType)
	dataPointMap := make(map[int64]float32)

	for _, point := range dataPoints {
		dataPointMap[point.Timestamp] = point.Value
	}

	var points []float32
	for i := startTimeStamp; i <= endTimeStamp; i++ {
		if value, present := dataPointMap[i]; present {
			//fmt.Println("appending value")
			points = append(points, value)
		} else {
			//fmt.Println("appending value")
			points = append(points, -1) //some point might have 0 value
		}
	}

	fillMissingValues(points)

	fmt.Println(points)
	return points, nil

}

func (dataFetcher *DataFetcher) SavePredictedData(agentIP string, containerId string, metric string, predictedValues []float32, startTimeStamp int64) error {

	var dataPoints []DataPoint
	for _, value := range predictedValues {
		startTimeStamp += 1
		point := DataPoint{
			AgentIp:     agentIP,
			ContainerId: containerId,
			Value:       value,
			Timestamp:   startTimeStamp,
			MetricType:  metric,
		}
		dataPoints = append(dataPoints, point)
	}

	conn := getConnection()
	err := saveData(dataPoints, conn)
	if err != nil {
		return err
	}

	conn.Close()
	return nil

}

//
func (dataFetcher *DataFetcher) GetPredictedData(agentIP string, containerId string, metric string, startTimeStamp int64, endTimeStamp int64) ([]float32, error) {

	dataPoints := getPredictedData(agentIP, containerId, metric)

	dataPointMap := make(map[int64]float32)

	var latestTimesStamp int64
	var oldestTimesStamp int64

	for i, point := range dataPoints {
		dataPointMap[point.Timestamp] = point.Value
		//fmt.Println(point.Value)
		if i == 0 {
			oldestTimesStamp = point.Timestamp
		}
		if point.Timestamp > latestTimesStamp {
			latestTimesStamp = point.Timestamp
		}

		if point.Timestamp < oldestTimesStamp {
			oldestTimesStamp = point.Timestamp
		}

	}
	fmt.Println("start,end", startTimeStamp, endTimeStamp)
	fmt.Println("Oldest ", oldestTimesStamp, "\n latest ", latestTimesStamp)

	var points []float32
	for i := startTimeStamp; i <= endTimeStamp; i++ {
		if value, present := dataPointMap[i]; present {
			//fmt.Println("appending value")
			points = append(points, value)
		} else {
			//fmt.Println("appending 0s")
			points = append(points, 0) //some point might have 0 value
		}
	}
	//fmt.Println(dataPointMap)
	return points, nil

}
