package queue

import (
	"errors"
	"sync"
	"time"
)

// Queue представляет очередь с сообщениями.
type Queue struct {
	messages    []Message
	mu          sync.Mutex
	maxMessages int
}

func NewQueue(maxMessages int) *Queue {
	return &Queue{maxMessages: maxMessages}
}

// QueueManager управляет очередями.
type QueueManager struct {
	queues      map[string]*Queue
	maxQueues   int
	maxMessages int
	queuesMu    sync.Mutex
}

// NewQueueManager создает новый QueueManager.
func NewQueueManager(maxQueues, maxMessages int) *QueueManager {
	return &QueueManager{
		queues:      make(map[string]*Queue),
		maxQueues:   maxQueues,
		maxMessages: maxMessages,
	}
}

// GetQueue возвращает очередь с указанным именем или создает новую, если она не существует.
func (qm *QueueManager) GetQueue(name string) (*Queue, error) {
	qm.queuesMu.Lock()
	defer qm.queuesMu.Unlock()

	if q, exists := qm.queues[name]; exists {
		return q, nil
	}

	if len(qm.queues) >= qm.maxQueues {
		return nil, errors.New("max number of queues reached")
	}

	qm.queues[name] = NewQueue(qm.maxMessages)
	return qm.queues[name], nil
}

// PutMessage добавляет сообщение в очередь.
func (q *Queue) PutMessage(message Message) error {
	if len(q.messages) >= q.maxMessages {
		return errors.New("max number of messages reached")
	}
	q.messages = append(q.messages, message)

	return nil
}

// GetMessage извлекает сообщение из очереди.
func (q *Queue) GetMessage(timeout int) (Message, error) {
	q.mu.Lock() // синхронизация доступа к очереди
	//Если 2 клиента обращаются к этой очереди одновременно, то они выстраиваются
	//в очередь по порядку обращений,
	//НО если их 3 и более то в момент когда произойдет разблокировка кто-то из них случайным
	//образом получит сообщение.
	defer q.mu.Unlock()

	// в очереди нет сообщений следует запустить таймер ожидания
	if len(q.messages) == 0 {
		now := time.Now()
		// буду раз в секунду проверять появились ли сообщения до тех пор, пока не появится сообщение или не истекает таймаут
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if len(q.messages) > 0 { // есть сообщение надо его отдать и выходить
				break
			}
			if timeout > 0 && now.Add(time.Second*time.Duration(timeout)).Before(time.Now()) { // истек таймаут
				return Message{}, errors.New("timeout")
			}
		}
	}
	msg := q.messages[0]
	q.messages = q.messages[1:]
	return msg, nil
}
