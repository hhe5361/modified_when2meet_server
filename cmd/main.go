package main

import (
	"better-when2meet/internal/server"
)

func main() {
	if err := server.InitServer(); err != nil {
		panic(err)
	}
}
