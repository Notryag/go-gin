package main

import "chapter05/routers"

func main() {
	server := routers.NewServer()
	server.Run(":8081")
}
