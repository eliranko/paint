package main

import (
	"encoding/json"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

var server *socketio.Server

const (
	room        = "def"
	canvasEvent = "canvas"
)

func getPushNotificationServer() *socketio.Server {
	var err error
	server, err = socketio.NewServer(nil)
	if err != nil {
		log.Panic("socketio creation error: ", err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		s.Join(room)
		log.Println("socketio client connected:", s.ID())
		return nil
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("socketio error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed ", reason)
	})

	go server.Serve()
	return server
}

func pushCanvases(canvas *Canvas) {
	buf, err := json.Marshal(canvas)
	if err != nil {
		log.Println("failed marshaing canvas")
		return
	}

	log.Printf("broadcasting %+v to %s", canvas, canvasEvent)
	server.BroadcastToRoom("", room, canvasEvent, string(buf))
}
