package main

import "kurento-client-go-generator/kurento"

var pipeline = new(kurento.MediaPipeline)
var master = new(kurento.WebRtcEndPoint)
var server = kurento.NewConnection('ws://127.0.0.1:8888')

func main() {
	
}
