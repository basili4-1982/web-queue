package args

import (
	"fmt"
	"os"
	"strconv"
)

// GetArgs извлекает аргументы командной строки для заполнения структуры Args
func GetArgs(defaultArgs Args) (Args, error) {
	args := os.Args[1:]

	var (
		res Args
		err error
	)

	res = defaultArgs
	if len(args) == 0 {
		return res, nil
	}

	// Перебираю все аргументы
	for i := 0; i < len(args); i++ {
		if args[i] == "--port" || args[i] == "-p" {
			if i+1 >= len(args) {
				return Args{}, fmt.Errorf("port is required")
			}
			res.Port, err = strconv.Atoi(args[i+1])
			if err != nil {
				return Args{}, fmt.Errorf("bad port: %s", args[i+1])
			}
		} else if args[i] == "--timeout" || args[i] == "-t" {
			if i+1 <= len(args) {
				res.Timeout, err = strconv.Atoi(args[i+1])
				if err != nil {
					return Args{}, fmt.Errorf("bad timeout: %s", args[i+1])
				}
			}
		} else if args[i] == "--max-queues" || args[i] == "-mq" {
			if i+1 <= len(args) {
				res.MaxQueues, err = strconv.Atoi(args[i+1])
				if err != nil {
					return Args{}, fmt.Errorf("bad max queues: %s", args[i+1])
				}
			}
		} else if args[i] == "--max-messages" || args[i] == "-m" {
			if i+1 <= len(args) {
				res.MaxMessages, err = strconv.Atoi(args[i+1])
				if err != nil {
					return Args{}, fmt.Errorf("bad max messages: %s", args[i+1])
				}
			}
		}
	}

	if res.Port < 0 || res.Port > 65535 {
		return Args{}, fmt.Errorf("bad port: %d", res.Port)
	}

	if res.MaxQueues < 1 {
		return Args{}, fmt.Errorf("bad max queues: %d", res.MaxQueues)
	}

	if res.MaxMessages < 1 {
		return Args{}, fmt.Errorf("bad max messages: %d", res.MaxMessages)
	}

	if res.Port == 0 {
		return Args{}, fmt.Errorf("port is required")
	}

	return res, nil
}
