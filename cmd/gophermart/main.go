package main

import "github.com/dmitrovia/gophermart/internal/processes/server"

func main() {
	go server.RunProcess()
}
