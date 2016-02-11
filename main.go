package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"os/signal"
	"os"
	"syscall"
	"github.com/LincolnFaradey/MadRobotsServer/game"
	"time"
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

			fmt.Fprintf(os.Stdout, "Restarting... %s\n", time.Now().Format(time.Stamp))
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
