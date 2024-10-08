package main

import (
	"fmt"
	"log"
	"net/http"

	"basili4-1982/web-queue/internal/http-handler"
	"basili4-1982/web-queue/pkg/args"
	"basili4-1982/web-queue/pkg/queue"
)

func main() {
	cfg, err := args.GetArgs(args.Args{
		Port:        0,
		MaxQueues:   10,
		MaxMessages: 100,
		Timeout:     5,
	})

	if err != nil {
		log.Fatalf("can't get args: %s", err.Error())
	}

	queueManager := queue.NewQueueManager(cfg.MaxQueues, cfg.MaxMessages)

	handler := http_handler.NewHandler(queueManager, cfg.Timeout)

	http.HandleFunc("PUT /queue/{qname}", handler.PutHandler)
	http.HandleFunc("GET /queue/{qname}", handler.GetHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}
