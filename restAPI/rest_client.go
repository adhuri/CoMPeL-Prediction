package restAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type QueryResponse struct {
	Clients []Client `json:"clients,omitempty"`
}

type Client struct {
	ClientIp   string   `json:"agentId,omitempty"`
	Containers []string `json:"containers,omitempty"`
}

func GetContainerDetails() *QueryResponse {
	qr := new(QueryResponse)
	err := restAPIClient("http://127.0.0.1:9091/query", qr)
	if err != nil {
		fmt.Println("Error", err)
	}
	//fmt.Println("Query response", qr)
	return qr
}

func restAPIClient(url string, q *QueryResponse) error {
	var myHTTPClient = http.Client{Timeout: 10 * time.Second}

	r, err := myHTTPClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	op, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Json received from Http 9091/query", string(op))

	err = json.Unmarshal(op, &q)
	if err != nil {
		return err
	}
	return nil
	//return json.NewDecoder(r.Body).Decode(&q)

}
