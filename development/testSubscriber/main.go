package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const serverUrl string = "http://127.0.0.1:8000"

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("called")
	b := r.Body
	decoder := json.NewDecoder(b)
	var data []WebHookResult
	decoder.Decode(&data)
	fmt.Printf("%+v\n", data)
}

func main() {
	subscribe()
	r := mux.NewRouter()
	r.HandleFunc("/webhook", webhookHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8001", r))
	log.Println("dispatch now data ")
	dispatch()
	// unsubscribe(id.Id)
}

type Subscribe struct {
	Webhook string `json:"webhook,omitempty"`
}
type subscribeResult struct {
	Id string `json:"id,omitempty"`
}

type WebHookResult struct {
	Date int64    `json:"date,omitempty"`
	Data TestData `json:"data,omitempty"`
}
type TestData struct {
	D1 string    `json:"d1,omitempty"`
	D2 TestData2 `json:"d2,omitempty"`
}
type TestData2 struct {
	D22 string `json:"d22,omitempty"`
}

func dispatch() {
	data := &TestData{
		D1: "lalal",
		D2: TestData2{
			D22: "haha",
		},
	}
	byteData, _ := json.Marshal(data)
	reader := *bytes.NewReader(byteData)
	_, err := http.Post(serverUrl+"/dispatch", "application/json", &reader)
	if err != nil {
		log.Panic(err)
	}
}

func unsubscribe(id string) {
	// res, err := http.Post(serverUrl+"/unsubscribe/"+id, "application/json", &reader)
}
func subscribe() subscribeResult {
	data := Subscribe{
		Webhook: "http://127.0.0.1:8001/webhook",
	}
	byteData, _ := json.Marshal(data)
	reader := *bytes.NewReader(byteData)
	res, err := http.Post(serverUrl+"/subscribe", "application/json", &reader)
	if err != nil {
		log.Panic(err)
	}
	log.Println(res)
	b := res.Body
	decoder := json.NewDecoder(b)
	var resBody subscribeResult
	decoder.Decode(&resBody)
	return resBody

}
