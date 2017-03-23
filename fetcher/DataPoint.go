package main

import (
	"fmt"
	"time"
)

type DataPoint struct {
	Timestamp  int64
	Value      float32
	MetricType string
}

func GetMetricDataForContainer(agentIp string, containerId string, metricType string, time int64) {

	dataPoints := GetData(agentIp, containerId, metricType)
	dataPointMap := make(map[int64]float32)

	var latestTimesStamp int64
	var oldestTimesStamp int64

	fmt.Println(len(dataPoints))

	for i, point := range dataPoints {
		dataPointMap[point.Timestamp] = point.Value
		//fmt.Printf("%d : %f \n", point.Timestamp, point.Value)
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

	fmt.Println(oldestTimesStamp)
	fmt.Println(latestTimesStamp)
	fmt.Println(time)

	var points []float32
	for i := oldestTimesStamp; i <= time; i += 1 {
		if value, present := dataPointMap[i]; present {
			points = append(points, value)
		} else {
			points = append(points, 0)
		}
	}

	fmt.Println(len(points))

	// for i, point := range points {
	// 	fmt.Printf("%d : %f \n", i, point)
	// }

}

func FillMissingValues(points []float32) {

}

func main() {

	GetMetricDataForContainer("192.168.0.26", "mysql_container_name", "cpu", time.Now().Unix())

}
