package http_handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	queue "beeline/pkg/queue"
)

type Handler struct {
	qm      *queue.QueueManager
	timeout int
}

func NewHandler(qm *queue.QueueManager, timeout int) *Handler {
	return &Handler{qm: qm, timeout: timeout}
}

// PutHandler обрабатывает запросы на добавление сообщений в очередь.
func (h *Handler) PutHandler(w http.ResponseWriter, r *http.Request) {
	queueName := r.PathValue("qname")

	var message queue.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	q, err := h.qm.GetQueue(queueName)
	if err != nil {
		http.Error(w, "error getting queue", http.StatusInternalServerError)
		return
	}

	if err := q.PutMessage(message); err != nil {
		http.Error(w, "queue is full", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetHandler обрабатывает запросы на извлечение сообщений из очереди.
func (h Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	queueName := r.PathValue("qname")

	var err error

	timeoutStr := r.URL.Query().Get("timeout")
	timeout := h.timeout
	if timeoutStr != "" {
		timeout, err = strconv.Atoi(timeoutStr)
		if err != nil {
			http.Error(w, "bad timeout", http.StatusBadRequest)
			return
		}
	}

	q, err := h.qm.GetQueue(queueName)
	if err != nil {
		http.Error(w, "queue not found", http.StatusNotFound)
		return
	}

	msgChan := make(chan queue.Message)
	errChan := make(chan error)

	go func() {
		msg, err := q.GetMessage(timeout)
		if err != nil {
			errChan <- err
		} else {
			msgChan <- msg
		}
	}()

	select {
	case <-r.Context().Done():
		http.Error(w, "timeout", http.StatusRequestTimeout)
	case msg := <-msgChan:
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(msg.Message))
	case err := <-errChan:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
