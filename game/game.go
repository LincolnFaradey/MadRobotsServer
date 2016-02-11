package game
import (
	"golang.org/x/net/websocket"
	"sync"
	"github.com/lincolnfaradey/madrobotsserver/user"
)



var once sync.Once
var instance *game = nil


type game struct {
	Players map[string] user.User
}

type Message struct {
	Method string `json:"method"`
	UserName string `json:"name"`
}

func New() *game {
	once.Do(func() {
		instance = &game{
			Players: make(map[string] user.User),
		}
	})
	return instance
}


func (g *game)Start(ws *websocket.Conn) {
	defer ws.Close()

	var player user.User
	for {
		if err := websocket.JSON.Receive(ws, &player); err != nil {
			panic(err)
			return
		}
		player.Socket = ws
		g.Players[player.Name] = player
		for k, v := range g.Players {
			if k == player.Name {
				continue
			}

			if err := websocket.JSON.Send(v.Socket, player); err != nil {
				delete(g.Players, k)

				go g.SendMessageAll(&Message{
						Method: "destroy",
						UserName: k,
					})
			}
		}
	}
}

func (g *game)SendMessageAll(msg *Message)  {
	for _, p := range g.Players {
		websocket.JSON.Send(p.Socket, &msg)
	}
}