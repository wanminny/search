package testing

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"expvar"
)

const NotAvailableMessage = "Not available"

var CommitHash string
var StartTime time.Time

type HeartbeatMessage struct {
	Status string `json:"status"`
	Build  string `json:"build"`
	Uptime string `json:"uptime"`
}

func init() {
	StartTime = time.Now()
}

var Counter = expvar.NewInt("visitor")


func handler(rw http.ResponseWriter, r *http.Request) {

	Counter.Add(1)

	hash := CommitHash
	if hash == "" {
		hash = NotAvailableMessage
	}
	uptime := time.Since(StartTime).String()
	err := json.NewEncoder(rw).Encode(HeartbeatMessage{"running", hash, uptime})
	if err != nil {
		log.Fatalf("Failed to write heartbeat message. Reason: %s", err.Error())
	}
}

func RunHeartbeatService(address string) {

	h := expvar.Handler()
	http.Handle("/stat",h)

	http.HandleFunc("/heartbeat", handler)
	log.Println(http.ListenAndServe(address, nil))
}

func Get(address string) (HeartbeatMessage, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", address, nil)
	//log.Println(req)
	resp, err := client.Do(req)
	if err != nil {
		return HeartbeatMessage{}, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return HeartbeatMessage{}, errors.New(fmt.Sprintf("Wrong status code: %d", resp.StatusCode))
	}
	message := HeartbeatMessage{}
	err = json.Unmarshal(b, &message)
	if err != nil {
		log.Println("Error occured unmarshalling the response")
	}
	return message, nil
}
