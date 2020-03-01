package main

import (
	"fmt"
	"io"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v2"
	"github.com/pkg/profile"
)

func doSignaling(offerPC, answerPC *webrtc.PeerConnection) {
	offer, err := offerPC.CreateOffer(nil)
	if err != nil {
		panic(err)
	}
	if err := offerPC.SetLocalDescription(offer); err != nil {
		panic(err)
	}
	if err := answerPC.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	answer, err := answerPC.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}
	if err := answerPC.SetLocalDescription(answer); err != nil {
		panic(err)
	}
	if err := offerPC.SetRemoteDescription(answer); err != nil {
		panic(err)
	}

}

func main() {
	defer profile.Start(profile.CPUProfile).Stop()

	offerPC, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	answerPC, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	var videoReceptionCount uint64
	answerPC.OnTrack(func(t *webrtc.Track, r *webrtc.RTPReceiver) {
		for {
			if _, err := t.ReadRTP(); err != nil {
				if err != io.EOF {
					panic(err)
				}
			}
			atomic.AddUint64(&videoReceptionCount, 1)
		}
	})

	track, err := offerPC.NewTrack(webrtc.DefaultPayloadTypeVP8, rand.Uint32(), fmt.Sprintf("video-%d", rand.Uint32()), fmt.Sprintf("video-%d", rand.Uint32()))
	if err != nil {
		panic(err)
	}

	if _, err := offerPC.AddTrack(track); err != nil {
		panic(err)
	}

	if _, err = answerPC.AddTransceiver(webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}

	doSignaling(offerPC, answerPC)
	rtpPacket := &rtp.Packet{
		Header: rtp.Header{
			Version: 2,
			SSRC:    track.SSRC(),
		},
	}

	for _ = range time.Tick(5 * time.Millisecond) {
		if atomic.LoadUint64(&videoReceptionCount) > 10000 {
			break
		}
		if err := track.WriteRTP(rtpPacket); err != nil {
			panic(err)
		}
	}

	offerPC.Close()
	answerPC.Close()
}
