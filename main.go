package main

import (
	"golang.org/x/net/websocket"
	"net/http"
	"fmt"
)

var activePlayers = make(map[string]*websocket.Conn)

type Player struct {
	Name string `json:"name"`
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func Game(ws *websocket.Conn)  {
	defer ws.Close()

	var player Player

	for {
		if err := websocket.JSON.Receive(ws, &player); err != nil {
			panic(err)
			return
		}

		activePlayers[player.Name] = ws
		fmt.Println(player)
		for k, v := range activePlayers {
			if k == player.Name { continue }

			if err := websocket.JSON.Send(v, player); err != nil {
				delete(activePlayers, k)

				go func() {
					for _, sock := range activePlayers {
						websocket.JSON.Send(sock, &Player{
							Name: k,
							X: 0,
							Y: 0,
						})
					}
				}()
			}
		}
	}
}

func main() {
	http.Handle("/", websocket.Handler(Game))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

