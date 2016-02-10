package game
import (
	"golang.org/x/net/websocket"
	"fmt"
)

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