package server

import (
	"math/rand"

	"github.com/pion/webrtc/v3"
)

var (
	videoTrack          *webrtc.Track
	peerConnectionCount int64
)

func init() {
	var err error
	videoTrack, err = webrtc.NewTrack(webrtc.DefaultPayloadTypeH264, rand.Uint32(), "pion-rtsp", "pion-rtsp", webrtc.NewRTPH264Codec(webrtc.DefaultPayloadTypeH264, 90000))
	if err != nil {
		panic(err)
	}
}

// Run run server
func Run() {
	go serveHTTP()
	go serveStream()
	// go serveReport()
}
