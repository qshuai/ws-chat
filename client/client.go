package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"sync"

	"github.com/gorilla/websocket"
)

var wg sync.WaitGroup

var (
	yellow = "\033[1;33m"
	green  = "\033[0;32m"

	end = "\033[0m"
)

var interrupt = make(chan os.Signal, 1)

func main() {
	signal.Notify(interrupt, os.Interrupt)

	uri := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(uri.String(), nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	err = c.WriteMessage(websocket.TextMessage, []byte("client connected!"))
	if err != nil {
		panic(err)
	}

	go readMsg(c)
	go sendMsg(c)

	wg.Add(1)
	wg.Wait()
}

func readMsg(conn *websocket.Conn) {
	for {
		select {
		case <-interrupt:
			conn.Close()
			return
		default:
			_, rawMsg, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
			}

			if rawMsg[len(rawMsg)-1] == '\n' {
				rawMsg = rawMsg[:len(rawMsg)-1]
			}
			fmt.Println(yellow+"server: "+end, green+string(rawMsg)+end)
		}
	}
}

func sendMsg(conn *websocket.Conn) {
	for {
		select {
		case <-interrupt:
			conn.Close()
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
				}
			}
		}
	}
}
