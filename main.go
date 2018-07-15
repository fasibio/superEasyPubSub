package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type handler struct {
	client *mongo.Client
	db     *mongo.Database
}

type subscriber struct {
	Webhook string `json:"webhook,omitempty"`
}

func (h handler) root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Super Easy Pub Sub"))
}

type deleteSubscribe struct {
	_id objectid.ObjectID `json:"_id,omitempty"`
}

type Dispatch = struct {
	Date string
	Data interface{}
}

func (h handler) dispatch(w http.ResponseWriter, r *http.Request) {
	b := r.Body
	decoder := json.NewDecoder(b)
	var data interface{}
	decoder.Decode(&data)
	d := Dispatch{
		Date: time.Now().Format("20060102150405"),
		Data: data,
	}
	_, err := h.db.Collection("data").InsertOne(context.Background(), d)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	cur, err := h.db.Collection("subscriber").Find(context.Background(), nil)
	if err != nil {
		log.Println(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var subs subscriber
		err := cur.Decode(&subs)
		if err != nil {
			log.Println(err)
		}
		log.Println(subs.Webhook)
	}

}

func (h handler) unsubscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	coll := h.db.Collection("subscriber")
	hex, err := objectid.FromHex(id)
	if err != nil {
		log.Println(err)
	}
	d := deleteSubscribe{
		_id: hex,
	}
	res, err := coll.DeleteOne(context.TODO(), d)
	if err != nil {
		log.Println(err)
	}
	if res.DeletedCount > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

/**
To subscribe a new Service to confirm by new Data
Need a JSON Body
{
	"webhook": "URL subscribe to"
}
*/
type subscribeResult struct {
	Id string `json:"id,omitempty"`
}

func (h handler) subscribe(w http.ResponseWriter, r *http.Request) {
	b := r.Body
	decoder := json.NewDecoder(b)
	var sub subscriber
	decoder.Decode(&sub)
	coll := h.db.Collection("subscriber")
	res, err := coll.InsertOne(context.Background(), sub)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	json, err := json.Marshal(subscribeResult{
		Id: res.InsertedID.(objectid.ObjectID).Hex(),
	})
	w.Write(json)
}

func main() {

	client := getDbHandler("mongodb://localhost:27018")

	var h = handler{
		client: client,
		db:     client.Database("simpleSubPub"),
	}
	log.Println("Starting Application")
	r := mux.NewRouter()
	r.HandleFunc("/", h.root).Methods("GET")
	r.HandleFunc("/subscribe", h.subscribe).Methods("POST")
	r.HandleFunc("/subscribe/{id}", h.unsubscribe).Methods("DELETE")
	r.HandleFunc("/dispatch", h.dispatch).Methods("POST")
	log.Println("Listen on port :8000")
	log.Fatal(http.ListenAndServe(":8000", r))

}

func getDbHandler(conStr string) *mongo.Client {

	client, err := mongo.NewClient(conStr)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	return client
}
