package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var offerClient *websocket.Conn
var answerClient *websocket.Conn

func checkStart() {
	if offerClient != nil && answerClient != nil {
		offerClient.WriteJSON(map[string]string{
			"type": "create_offer",
		})
	}
}
func wsHandler(w http.ResponseWriter, r *http.Request) {
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	fmt.Println("New WS connection")
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	for {
		var obj map[string]any
		err := conn.ReadJSON(&obj)
		if err != nil {
			break
		}
		log.Println("received:", obj)
		switch obj["type"] {
		case "connect":
			if offerClient == nil {
				offerClient = conn
				conn.WriteJSON(map[string]any{
					"type":    "connect",
					"code":    200,
					"message": "connect success",
				})
				checkStart()
			} else if answerClient == nil {
				answerClient = conn
				conn.WriteJSON(map[string]any{
					"type":    "connect",
					"code":    200,
					"message": "connect success",
				})
				checkStart()
			} else {
				conn.WriteJSON(map[string]any{
					"type":    "connect",
					"code":    -1,
					"message": "connect fail",
				})
				conn.Close()
			}
		case "offer":
			if answerClient != nil {
				answerClient.WriteJSON(obj)
			}
		case "answer":
			if offerClient != nil {
				offerClient.WriteJSON(obj)
			}
		case "offer_ice":
			if answerClient != nil {
				answerClient.WriteJSON(obj)
			}
		}
	}

	if conn == offerClient {
		log.Println("remove offerClient")
		offerClient = nil
	} else if conn == answerClient {
		log.Println("remove answerClient")
		answerClient = nil
	}
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New HTTP connection")
	byteData, err := os.ReadFile("index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(byteData)
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", indexHandler)
	log.Println("Server start at :9004")
	log.Fatal(http.ListenAndServe(":9004", nil))
}
