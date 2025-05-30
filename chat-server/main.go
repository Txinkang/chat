package main

import "chat-server/initialize"

func main() {
	initialize.Initialize()

	select {}
}
