package fetcher

import "fmt"

type DataPoint struct {
	Timestamp  int64
	Value      float32
	MetricType string
}

func GetAgentInformation() {

}

func GetMetricDataForContainer(agentIp string, containerId string, metricType string, time int64, numberOfPoints int) []float32 {

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
		// if there is break in time series, aligning will be impossible with 2 seconds sampling
		if value, present := dataPointMap[i]; present {
			points = append(points, value)
		} else {
			points = append(points, -1) //some point might have 0 value
		}
	}

	fmt.Println(len(points))

	// for i, point := range points {
	// 	fmt.Printf(" %d : %f \n", i, point)
	// }

	FillMissingValues(points)

	// for i, point := range points {
	// 	fmt.Printf(" %d : %f \n", i, point)
	// }

	if len(points) > numberOfPoints {
		// Trim the slice if we have more points than asked for
		numberOfExtraPoints := len(points) - numberOfPoints
		fmt.Println(len(points[numberOfExtraPoints:]))
		return points[numberOfExtraPoints:]
	} else if len(points) < numberOfPoints {
		// If points are less than required then pad 0 at the start
		var remainingPoints []float32
		numberOfPointsMissing := numberOfPoints - len(points)
		for i := 0; i < numberOfPointsMissing; i += 1 {
			remainingPoints = append(remainingPoints, 0)
		}
		remainingPoints = append(remainingPoints, points...)

		fmt.Println(len(remainingPoints))
		return remainingPoints

	}

	return points

}

func FillMissingValues(points []float32) {
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
