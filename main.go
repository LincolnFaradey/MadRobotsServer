package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"os/signal"
	"os"
	"syscall"
)

var activePlayers = make(map[string]Player)

type Player struct {
	Name   string          `json:"name"`
	X      float32         `json:"x"`
	Y      float32         `json:"y"`
	Socket *websocket.Conn `json:"-"`
}

type Message struct {
	Method string `json:"method"`
	UserName string `json:"name"`
}

func Game(ws *websocket.Conn) {
	defer ws.Close()

	var player Player
	for {
		if err := websocket.JSON.Receive(ws, &player); err != nil {
			panic(err)
			return
		}
		player.Socket = ws
		activePlayers[player.Name] = player
		fmt.Println(player)
		for k, v := range activePlayers {
			if k == player.Name {
				continue
			}

			if err := websocket.JSON.Send(v.Socket, player); err != nil {
				delete(activePlayers, k)

				go func() {
					sendMessageAll(&Message{
						Method: "destroy",
						UserName: k,
					})
				}()
			}
		}
	}
}

func sendMessageAll(msg *Message)  {
	for _, p := range activePlayers {
		websocket.JSON.Send(p.Socket, &msg)
	}
}

func run() error {
	return http.ListenAndServe(":8080", nil)
}

func supervisor(errCh <-chan error) {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM)

	go func() {
		for {
			select {
				case <-errCh:
					run()
				case <-sig:
					run()
			}

			fmt.Fprintf(os.Stdout, "Restarting...\n")
		}
	}()
}

func main() {
	errCh := make(chan error)
	http.Handle("/", websocket.Handler(Game))
	go supervisor(errCh)
	errCh<- run()

	select {}
}
