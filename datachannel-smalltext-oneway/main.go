package main

import (
	"sync/atomic"
	"time"

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
	defer profile.Start().Stop()

	offerPC, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	answerPC, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	var msgCount uint64
	answerPC.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			atomic.AddUint64(&msgCount, 1)
		})
	})

	dc, err := offerPC.CreateDataChannel("benchChannel", nil)
	if err != nil {
		panic(err)
	}

	hasOpened := make(chan interface{})
	dc.OnOpen(func() {
		close(hasOpened)
	})

	doSignaling(offerPC, answerPC)
	<-hasOpened

	for atomic.LoadUint64(&msgCount) < 10000 {
		if err := dc.SendText("foobar"); err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Millisecond)
	}

	offerPC.Close()
	answerPC.Close()
}
