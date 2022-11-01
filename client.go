package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	Conf "github.com/NotRoyadma/BDClient/config"
	"github.com/NotRoyadma/BDClient/logger"
	router "github.com/NotRoyadma/BDClient/router"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		logger.DeleteLogFiles()
		os.Exit(0)
	}()

	http.HandleFunc("/", router.DefaultHandler)
	log.Println("Listening on " + strconv.Itoa(Conf.Conf.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Conf.Conf.Port), nil))
}
