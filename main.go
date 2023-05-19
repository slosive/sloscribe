/*
Copyright © 2023 Oluwole Fadeyi
*/
package main

import (
	"context"
	"github.com/tfadeyi/sloth-simple-comments/cmd"
	"github.com/tfadeyi/sloth-simple-comments/internal/logging"
	"os"
	"os/signal"
	"syscall"
)

//@aloe name sloth-simple-comments
//@aloe url https://github.com/tfadeyi/sloth-simple-comments
//@aloe version v0.0.1
//@aloe description This is a CLI tool to parse source code and generate a Sloth definition manifest.

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	log := logging.NewStandardLogger()
	ctx = logging.ContextWithLogger(ctx, log)

	cmd.Execute(ctx)
}
