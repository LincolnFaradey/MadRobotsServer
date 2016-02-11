package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"os/signal"
	"os"
	"syscall"
	"github.com/lincolnfaradey/madrobotsserver/game"
)

var newGame = game.New()


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
	http.Handle("/", websocket.Handler(newGame.Start))
	go supervisor(errCh)
	errCh<- run()

	select {}
}
