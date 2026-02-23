package main

import (
	"context"
	"fmt"
	"microdois/internal/infra"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config := infra.Load()
	container := infra.NewContainer(config, ctx)

	fmt.Println("Starting server...", container.Config.Environment)
}
