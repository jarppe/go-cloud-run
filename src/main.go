package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetPrefix("go-cloud-run: ")
	log.Print("Server starting...")


	s := NewServer()
	s.Start()
	WaitForSignal()
	s.Shutdown()

	log.Printf("Server terminated")
	os.Exit(0)
}

func WaitForSignal() {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGUSR2)
	sig := <-sigCh
	signal.Reset(sig)
	log.Printf("Got signal %q, terminating...", sig)
}
