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
	MyDB     = "square_holes"
	username = "bubba"
	password = "bumblebeetuna"
)

func GetConnection() influx.Client {

	// Create a new HTTPClient
	conn, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     "http://localhost:10090",
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

func GetData(agentIp string, containerId string, metric string) []DataPoint {

	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     "http://localhost:10090",
		Username: username,
		Password: password,
	})

	defer c.Close()

	if err != nil {
		log.Fatal(err)
	}

	q := fmt.Sprintf("select %s from container_data where agent = '%s' and container = '%s' ORDER BY time DESC LIMIT 20", metric, agentIp, containerId)

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
