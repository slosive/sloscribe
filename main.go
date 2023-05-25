/*
MIT License

Copyright © 2023 Oluwole Fadeyi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
//go:generate aloe-cli spec generate main.go --format yaml,markdown
//go:generate gomarkdoc --output "{{.Dir}}/README.md" ./internal/... ./cmd/...

package main

import (
	"context"
	"github.com/tfadeyi/sloth-simple-comments/cmd"
	"github.com/tfadeyi/sloth-simple-comments/internal/logging"
	"os"
	"os/signal"
	"syscall"
)

// @aloe name slotalk
// @aloe url https://tfadeyi.github.io
// @aloe version v0.0.1
// @aloe description Generate Sloth SLO/SLI definitions from code annotations.

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer cancel()

	log := logging.NewStandardLogger()
	ctx = logging.ContextWithLogger(ctx, log)

	cmd.Execute(ctx)
}
