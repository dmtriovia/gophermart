package main

import (
	"sync"

	"github.com/dmitrovia/gophermart/internal/processes/server"
)

func main() {
	waitGroup := new(sync.WaitGroup)

	go server.RunProcess(waitGroup)

	waitGroup.Add(1)
	waitGroup.Wait()
}
