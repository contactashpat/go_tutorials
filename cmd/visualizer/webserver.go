package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"go_tutorials/internal/web"
)

// WebServerCommand hosts the HTTP API for visualisations.
type WebServerCommand struct{}

// NewWebServerCommand constructs the web command.
func NewWebServerCommand() *WebServerCommand {
	return &WebServerCommand{}
}

// Run starts the HTTP server.
func (c *WebServerCommand) Run(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	addr := fs.String("addr", ":8080", "address to listen on")
	if err := fs.Parse(args); err != nil {
		return err
	}

	srv := web.NewServer()
	mux := http.NewServeMux()
	srv.Routes(mux)

	log.Printf("HTTP visualiser listening on %s", *addr)
	return http.ListenAndServe(*addr, mux)
}

func init() {
	registerCommand("serve", func() Command { return NewWebServerCommand() })
	registerCommand("help-serve", func() Command {
		return CommandFunc(func([]string) error {
			fmt.Println("Usage: go run ./cmd/visualizer serve --addr :8080")
			return nil
		})
	})
}
