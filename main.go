package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pion/webrtc/v3"
)

var (
	videoTrack          *webrtc.Track
	peerConnectionCount int64
)

func main() {
	var err error
	videoTrack, err = webrtc.NewTrack(webrtc.DefaultPayloadTypeH264, 5000, "pion-rtsp", "pion-rtsp", webrtc.NewRTPH264Codec(webrtc.DefaultPayloadTypeH264, 90000))
	if err != nil {
		panic(err)
	}

	go serveHTTP()
	go serveStream()
	// go serveReport()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println(sig)
		done <- true
	}()

	log.Println("Server Start Awaiting Signal")
	<-done
	log.Println("Exiting")
}
