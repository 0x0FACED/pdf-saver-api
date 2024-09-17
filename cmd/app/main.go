package main

import "github.com/0x0FACED/pdf-saver-api/internal/server"

func main() {
	if err := server.Start(); err != nil {
		panic(err)
	}
}
