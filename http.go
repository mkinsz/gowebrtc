package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/pion/webrtc/v3"
)

func serveHTTP() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/codec", codecHandler)
	http.HandleFunc("/signal", signalHandler)

	fmt.Println("Open http://localhost:8080 to access this demo")
	panic(http.ListenAndServe(":8080", nil))
}

func codecHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	// b, err := json.Marshal(codecs)
	// if err == nil {
	// _, err := w.Write(b)
	// if err == nil {
	// 	log.Println("Write Codec Info error", err)
	// 	return
	// }
	// }
}

// HTTP Handler that accepts an Offer and returns an Answer
// adds videoTrack to PeerConnection
func signalHandler(w http.ResponseWriter, r *http.Request) {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		if connectionState == webrtc.ICEConnectionStateDisconnected {
			atomic.AddInt64(&peerConnectionCount, -1)
			if err := peerConnection.Close(); err != nil {
				panic(err)
			}
		} else if connectionState == webrtc.ICEConnectionStateConnected {
			atomic.AddInt64(&peerConnectionCount, 1)
		}
	})

	if _, err = peerConnection.AddTrack(videoTrack); err != nil {
		panic(err)
	}

	var offer webrtc.SessionDescription
	if err = json.NewDecoder(r.Body).Decode(&offer); err != nil {
		panic(err)
	}

	if err = peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	gatherCompletePromise := webrtc.GatheringCompletePromise(peerConnection)

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	} else if err = peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	<-gatherCompletePromise

	response, err := json.Marshal(*peerConnection.LocalDescription())
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		panic(err)
	}
}
