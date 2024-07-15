package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type clientData struct {
	ClientID  string  `json:"clientId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", handler)

	http.HandleFunc("GET /ws", wsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error to upgrade:", err)
	}

	defer conn.Close()

	clientId, err := uuid.NewRandom()
	if err != nil {
		log.Println("Error to generating clientId:", err)
		return
	}
    
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		var data clientData
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("JSON Unmarshal error:", err)
			return
		}

		data.ClientID = clientId.String()

		dataJSON, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshalling data message:", err)
			return
		}

		log.Printf("Received clientId: %s location: Latitude: %f, Longitude: %f\n", data.ClientID, data.Latitude, data.Longitude)

		err = conn.WriteMessage(messageType, dataJSON)
		if err != nil {
			log.Println("Write:", err)
			return
		}
	}
}
