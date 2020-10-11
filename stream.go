package main

import (
	"io"
	"log"
	"time"

	"github.com/deepch/vdk/codec/h264parser"
	"github.com/deepch/vdk/format/rtsp"
	"github.com/pion/webrtc/v3/pkg/media"
)

// Connect to a RTSP URL and pull media
// Convert H264 to Annex-B, thesn write to `videoTrack` which sends to all PeerConnections
const rtspURL = "rtsp://10.67.24.94:8554/a170432"

func serveStream() {
	annexbNALUStartCode := func() []byte { return []byte{0x00, 0x00, 0x00, 0x01} }

	for {
		session, err := rtsp.Dial(rtspURL)
		if err != nil {
			panic(err)
		}
		session.RtpKeepAliveTimeout = 10 * time.Second

		codecs, err := session.Streams()
		if err != nil {
			panic(err)
		}
		// else if len(codecs) != 1 || codecs[0].Type() != av.H264 {
		// 	panic("RTSP feed must be a single H264 codec")
		// }

		var previousTime time.Duration
		for {
			pkt, err := session.ReadPacket()
			if err != nil {
				break
			}
			pkt.Data = pkt.Data[4:]

			// For every key-frame pre-pend the SPS and PPS
			if pkt.IsKeyFrame {
				pkt.Data = append(annexbNALUStartCode(), pkt.Data...)
				pkt.Data = append(codecs[0].(h264parser.CodecData).PPS(), pkt.Data...)
				pkt.Data = append(annexbNALUStartCode(), pkt.Data...)
				pkt.Data = append(codecs[0].(h264parser.CodecData).SPS(), pkt.Data...)
				pkt.Data = append(annexbNALUStartCode(), pkt.Data...)
			}

			bufferDuration := pkt.Time - previousTime
			if pkt.Idx == 0 {
				samples := media.NSamples(bufferDuration, 90000)
				if err = videoTrack.WriteSample(media.Sample{Data: pkt.Data, Samples: samples}); err != nil && err != io.ErrClosedPipe {
					panic(err)
				}
				previousTime = pkt.Time
				// log.Println("Stream: ", pkt.Time, len(pkt.Data), samples, previousTime)
			}
		}

		if err = session.Close(); err != nil {
			log.Println("session Close error", err)
		}

		time.Sleep(5 * time.Second)
	}
}
