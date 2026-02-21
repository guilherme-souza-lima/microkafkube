package main

import (
	"fmt"
	"microum/infra"
)

func main() {
	config := infra.Load()
	container := infra.NewContainer(config)
	fmt.Printf("run server %s:%s\n", container.Config.ServerName, container.Config.ServerPort)
}
