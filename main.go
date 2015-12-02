package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

var activePlayers = make(map[string]Player)

type Player struct {
	Name   string          `json:"name"`
	X      float32         `json:"x"`
	Y      float32         `json:"y"`
	Socket *websocket.Conn `json:"-"`
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
					for _, p := range activePlayers {
						websocket.JSON.Send(p.Socket, &Player{
							Name: k,
							X:    0,
							Y:    0,
						})
					}
				}()
			}
		}
	}
}

func run() error {
	return http.ListenAndServe(":8080", nil)
}

//func supervisor() {
//	sig := make(chan os.Signal)
//	signal.Notify(sig, syscall.SIGTERM)
//
//	go func() {
//		for {
//			<-sig
//			fmt.Fprintf(os.Stdout, "Restarting...")
//			if err := run(); err != nil {
//				panic(err)
//			}
//		}
//	}()
//}

func main() {
	http.Handle("/", websocket.Handler(Game))
//	go supervisor()
	for {
		if err := run(); err != nil {
			panic(err)
		}
	}
}
