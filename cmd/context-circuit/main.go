package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/kaotypr/context-circuit-source/internal/cli"
)

var version = "2.0.0-dev"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	os.Exit(cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr, version))
}
