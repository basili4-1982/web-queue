package main

import (
	"fmt"
	"log"
	"net/http"

	"beeline/pkg/args"
	http_handler "beeline/pkg/http-handler"
	"beeline/pkg/queue"
)

func main() {
	cfg, err := args.GetArgs(args.Args{
		Port:        0,
		MaxQueues:   10,
		MaxMessages: 100,
		Timeout:     5,
	})

	if err != nil {
		log.Fatalf("can't get args: %w", err)
	}

	queueManager := queue.NewQueueManager(cfg.MaxQueues, cfg.MaxMessages)

	handler := http_handler.NewHandler(queueManager, cfg.Timeout)

	http.HandleFunc("PUT /queue/{qname}", handler.PutHandler)
	http.HandleFunc("GET /queue/{qname}", handler.GetHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}
