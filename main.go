package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Message struct {
	Txid          string
	Vout          int
	Value         string
	Height        int
	Confirmations int
}

type Response struct {
	Confirmed   int
	Unconfirmed int
}

// Message Sample
// {
// 	txid: "93f9c795b3fbb5e0b21d16c84f00a42f07a60202be78e6c17b14fcaaf8d4f7a1",
// 	vout: 0,
// 	value: "430137",
// 	confirmations: 0
//  }

var response Response

// Enable CORS func
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// Get bitcoin address transactions
func getAddressInfo(w http.ResponseWriter, r *http.Request) {
	// Sending api request
	params := mux.Vars(r) // Gets params
	var address string = params["address"]
	s := fmt.Sprintf("https://blockbook-bitcoin.tronwallet.me/api/v2/utxo/%s", address)

	resp, err := http.Get(s)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//Convert the body to type message
	var m []Message

	error := json.Unmarshal(body, &m)

	confirmed := 0
	unconfirmed := 0

	for _, element := range m {

		i, err := strconv.Atoi(element.Value)

		if err != nil {
			fmt.Println((err))
		}

		if element.Confirmations < 2 {
			unconfirmed += i
		} else {
			confirmed += i
		}
	}

	fmt.Println("Confirmed", confirmed)
	fmt.Println("Unconfirmed", unconfirmed)

	response = Response{Confirmed: confirmed, Unconfirmed: unconfirmed}

	enableCors(&w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	if error != nil {
		log.Fatal(error)
	}

}

// Main function
func main() {

	port := ":8000"

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/balance/{address}", getAddressInfo).Methods("GET")

	//server message
	fmt.Println("Running server on port number", port)

	// Start server
	log.Fatal(http.ListenAndServe(port, r))
}

// @todo - Unit Tests
