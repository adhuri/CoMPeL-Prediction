package fetcher

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

const (
	MyDB            = "square_holes"
	username        = "bubba"
	password        = "bumblebeetuna"
	InfluxDBAddress = "http://influxdb:8086"
)

func getConnection() influx.Client {

	// Create a new HTTPClient
	conn, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     InfluxDBAddress,
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatal(err)
	}

	return conn
}

// queryDB convenience function to query the database
func queryDB(clnt influx.Client, cmd string) (res []influx.Result, err error) {
	q := influx.Query{
		Command:  cmd,
		Database: MyDB,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func getData(agentIp string, containerId string, metric string) []DataPoint {

	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     InfluxDBAddress,
		Username: username,
		Password: password,
	})

	defer c.Close()

	if err != nil {
		log.Fatal(err)
	}

	q := fmt.Sprintf("select %s from container_data where agent = '%s' and container = '%s' ORDER BY time DESC LIMIT 10000", metric, agentIp, containerId)

	res, err := queryDB(c, q)
	if err != nil {
		log.Fatal(err)
	}
	if len(res) == 0 {
		panic("Result is empty for given query")
	}

	if len(res[0].Series) == 0 {
		panic("Series is empty for given query")
	}

	var dataPoints []DataPoint
	for _, value := range res[0].Series[0].Values {
		timeStamp, ok := value[0].(string)
		if ok {
			t, err := time.Parse(time.RFC3339, timeStamp)
			if err != nil {
				panic("Unable to parse date")
			}
			tm := t.Unix()

			value, ok := value[1].(json.Number)
			if !ok {
				continue
			}
			floatValue, err := strconv.ParseFloat(value.String(), 32)
			if err != nil {
				continue
			}

			dataPoint := DataPoint{
				Timestamp:  tm,
				Value:      float32(floatValue),
				MetricType: "cpu",
			}
			dataPoints = append(dataPoints, dataPoint)
		}
	}
	return dataPoints
}

func saveData(dataPoints []DataPoint, conn influx.Client) error {
	// Create a new point batch
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, dataPoint := range dataPoints {
		tags := map[string]string{
			"agent":     dataPoint.AgentIp,
			"container": dataPoint.ContainerId,
			"metric":    dataPoint.MetricType,
		}
		fields := map[string]interface{}{
			"value": dataPoint.Value,
		}

		tm := time.Unix(dataPoint.Timestamp, 0)

		pt, err := influx.NewPoint("predicted_container_data", tags, fields, tm)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)

		// Write the batch
		if err := conn.Write(bp); err != nil {
			return err
		}
	}

	return nil

}

func getPredictedData(agentIp string, containerId string, metric string) []DataPoint {

	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     InfluxDBAddress,
		Username: username,
		Password: password,
	})

	defer c.Close()

	if err != nil {
		log.Fatal(err)
	}

	q := fmt.Sprintf("select value from predicted_container_data where agent = '%s' and container = '%s' and metric = '%s' ORDER BY time DESC LIMIT 11000", agentIp, containerId, metric)

	res, err := queryDB(c, q)
	if err != nil {
		log.Fatal(err)
	}
	if len(res) == 0 {
		panic("Result is empty for given query")
	}

	//fmt.Println(res)

	if len(res[0].Series) == 0 {
		panic("Series is empty for given query")
	}

	var dataPoints []DataPoint
	for _, value := range res[0].Series[0].Values {
		timeStamp, ok := value[0].(string)
		if ok {
			t, err := time.Parse(time.RFC3339, timeStamp)
			if err != nil {
				panic("Unable to parse date")
			}
			tm := t.Unix()

			value, ok := value[1].(json.Number)
			if !ok {
				continue
			}
			floatValue, err := strconv.ParseFloat(value.String(), 32)
			if err != nil {
				continue
			}

			dataPoint := DataPoint{
				Timestamp:  tm,
				Value:      float32(floatValue),
				MetricType: metric,
			}
			dataPoints = append(dataPoints, dataPoint)
		}
	}
	//fmt.Println(dataPoints)
	return dataPoints
}
