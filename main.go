package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"
)

type podInfo struct {
	Name   string
	IPAddr string
	Node   string
}

var pod podInfo
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func getPodInfoHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, pod)
}

func main() {
	pod = podInfo{
		Name:   os.Getenv("POD_NAME"),
		IPAddr: os.Getenv("IPADDR"),
		Node:   os.Getenv("NODE_NAME"),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", getPodInfoHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("server started")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("shudown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown:", err)
	}

	log.Println("server exiting")
}
