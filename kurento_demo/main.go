package main

import "kurento-client-go-generator/kurento"
import "net/http"

var pipeline = new(kurento.MediaPipeline)
var master = new(kurento.WebRtcEndPoint)
var server = kurento.NewConnection('ws://127.0.0.1:8888')

func main() {
	http.HandleFunc("/static/", staticExecute)
	http.ListenAndServe(":8080", nil)
}


func staticExecute(response http.Response, request *http.Request)  {
	requestUrl := request.URL.String()
	filepath := requestUrl[len("/static"):]
}