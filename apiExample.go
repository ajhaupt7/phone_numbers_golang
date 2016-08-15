package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

// var templates = template.Must(template.ParseFiles("index.html", "phoneData.html"))

type Numverify struct {
	Valid               bool   `json:"valid"`
	Number              string `json:"number"`
	LocalFormat         string `json:"local_format"`
	InternationalFormat string `json:"international_format"`
	CountryCode         string `json:"country_code"`
	CountryName         string `json:"country_name"`
	Location            string `json:"location"`
	Carrier             string `json:"carrier"`
	LineType            string `json:"line_type"`
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func sendRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	phone := r.Form["phone_number"][0]
	record, err := getPhoneData(phone)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}
	fmt.Println("Phone No. = ", record.InternationalFormat)
	fmt.Println("Country   = ", record.CountryName)
	fmt.Println("Location  = ", record.Location)
	fmt.Println("Carrier   = ", record.Carrier)
	fmt.Println("LineType  = ", record.LineType)

	t, _ := template.ParseFiles("phoneData.html")
	t.Execute(w, record)
}

func getPhoneData(phone string) (*Numverify, error) {
	safePhone := url.QueryEscape(phone)
	url := fmt.Sprintf("http://apilayer.net/api/validate?access_key=2d83b310f440ffa8d9d92015c20ad8ed&number=%s", safePhone)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record Numverify

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}

	return &record, nil
}

func main() {
	http.HandleFunc("/", getIndex)
	http.HandleFunc("/sendRequest", sendRequest)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
