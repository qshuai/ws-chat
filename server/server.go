package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var (
	yellow = "\033[1;33m"
	green  = "\033[0;32m"

	end = "\033[0m"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	fmt.Println(yellow + "Start conversation..." + end)
	http.ListenAndServe("localhost:8080", nil)
}

var done = make(chan struct{})

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go readMsg(conn)
	go sendMsg(conn)
}

func readMsg(conn *websocket.Conn) {
	for {
		select {
		case <-done:
			// end up the ws connection
			return
		default:
			_, rawMsg, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
				close(done)
			}

			if rawMsg[len(rawMsg)-1] == '\n' {
				rawMsg = rawMsg[:len(rawMsg)-1]
			}
			fmt.Println(yellow+"client: "+end, green+string(rawMsg)+end)
		}
	}
}

func sendMsg(conn *websocket.Conn) {
	for {
		select {
		case <-done:
			return
		default:
			input := bufio.NewReader(os.Stdin)
			msg, err := input.ReadString('\n')
			if err != nil {
				return
			}

			if len(msg) > 1 {
				err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					conn.Close()
					close(done)
				}
			}
		}
	}
}
