package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type jsonPacket struct {
	Data string
	Id   string
}

type HelloPacket struct {
	jsonPacket
}

type jsonPacketHandler struct {
	Id string
}

type jsonHandlePacketler interface {
	handlePacket(string) (jsonPacket, error)
}

var jsonPacketHandlers = make(map[string]jsonHandlePacketler)

func loadPacketHandlers() {
	loadTestPacketHandler()
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	if len(jsonPacketHandlers) == 0 {
		loadPacketHandlers()
	}

	fmt.Println("Incoming web socket request:", r.URL.Path)
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		fmt.Println(err)
		return
	}

	if err := conn.WriteJSON(&HelloPacket{jsonPacket{Id: "Hello"}}); err != nil {
		fmt.Println("Could not send JSON:", err.Error())
		return
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Could not read message:", err.Error())
			return
		}

		var data jsonPacket
		if err := json.Unmarshal(p, &data); err != nil {
			fmt.Println("Could not read json:", err.Error())
			return
		}
		fmt.Println(data)

		// MiaTODO: check if contains
		packet, err := jsonPacketHandlers[data.Id].handlePacket(data.Data)
		if err != nil {
			return
		}
		if packet.Id != nil {
			conn.WriteJSON(packet)
		}
	}
}
