package user

import "golang.org/x/net/websocket"

type User struct {
	Name   string          `json:"name"`
	X      float32         `json:"x"`
	Y      float32         `json:"y"`
	Socket *websocket.Conn `json:"-"`
}