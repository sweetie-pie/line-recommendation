package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"io"
	"log"
	"net/url"
)

var (
	addr           string
	peerConnection *webrtc.PeerConnection
)

func main() {
	flag.StringVar(&addr, "a", "localhost:7000", "address to use")
	flag.Parse()

	u := url.URL{
		Scheme: "ws",
		Host:   addr,
		Path:   "/ws",
	}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Fatal("close fatal:", err)
		}
	}(c)

	done := make(chan struct{})

	go readMessage(c, done)

	<-done

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
	}

	mediaEngine := webrtc.MediaEngine{}
}

func readMessage(connection *websocket.Conn, done chan struct{}) {
	defer close(done)

	for {
		_, message, err := connection.ReadMessage()
		if err != nil || err == io.EOF {
			log.Fatal("Error reading: ", err)
			return
		}

		fmt.Printf("recv: %s", message)
	}
}
