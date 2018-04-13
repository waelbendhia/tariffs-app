package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/waelbendhia/tariffs-app/app"
)

func main() {
	var (
		sigs = make(chan os.Signal, 1)
		done = make(chan bool, 1)
	)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println("Received signal")
		log.Println(sig)
		done <- true
	}()

	app.Start(done)
}
