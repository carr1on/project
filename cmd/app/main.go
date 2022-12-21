package main

import (
	"Statistic/internal/controller"
	App "Statistic/internal/controller/app"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	var wg = sync.WaitGroup{}
	s := controller.CreateNewServer()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	App.GUIapp(s, &wg, sig)

	wg.Wait()
	log.Print("we have a way out")
	log.Print("we leave")
}
