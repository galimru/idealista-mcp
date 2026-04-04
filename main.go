package main

import (
	"fmt"
	"os"

	"github.com/galimru/idealista-mcp/tools"
	"github.com/mark3labs/mcp-go/server"
)

// Populated by -ldflags during build; see Makefile.
var version = "dev"

func main() {
	s := newServer()
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func newServer() *server.MCPServer {
	s := server.NewMCPServer(
		"idealista-mcp",
		version,
		server.WithToolCapabilities(false),
	)

	runtime := tools.NewRuntimeProvider()
	tools.RegisterLocationTools(s, runtime)
	tools.RegisterSearchTools(s, runtime)
	tools.RegisterDetailTools(s, runtime)
	tools.RegisterStatsTools(s, runtime)

	return s
}
