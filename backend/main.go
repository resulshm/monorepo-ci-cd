package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/monorepo-backend/config"
	"example.com/monorepo-backend/web"

	log "github.com/sirupsen/logrus"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	quitChan := make(chan interface{})

	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "01-02 15:04:05.000"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(log.DebugLevel)

	err := config.ReadConfig("config.json")
	if err != nil {
		log.WithError(err).Panic("error reading config file")
	}

	setupServer(quitChan, signalChan, config.Conf)
}

func setupServer(quit chan interface{}, signalChan chan os.Signal, conf *config.Config) {
	s := web.NewServer()
	r := web.NewRouter(s)

	srv := &http.Server{
		Addr:         conf.ListenAddress,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 45 * time.Second,
	}

	listener, err := net.Listen("tcp", conf.ListenAddress)
	if err != nil {
		log.WithError(err).Panic(" setting up listener")
		return
	}

	log.WithField("listen", conf.ListenAddress).Info("Starting HTTP API Server")
	fmt.Println("<--START-SERVER-->")
	go startServer(srv, listener)

	for {
		select {
		case <-quit:
			log.Warn("quit channel closed, closing listener")
			err = srv.Close()
			if err != nil {
				log.WithError(err).Error("error during HTTP Server close")
			}
			err = listener.Close()
			if err != nil {
				log.WithError(err).Error("error during TCP Listener close")
			}
			return
		case sig := <-signalChan:
			switch sig {
			case os.Interrupt, os.Kill, syscall.SIGTERM:
				log.Info("interrupt signal received, sending Quit signal")
				close(quit)
			default:
				log.WithField("signal", sig).Info("signal received")
			}
		}
	}
}

func startServer(srv *http.Server, listener net.Listener) {
	err := srv.Serve(listener)
	if err != nil {
		log.WithError(err).Error("HTTP server Error")
	}
}
