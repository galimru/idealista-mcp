package main

import (
	"fmt"
	"os"

	"github.com/galimru/idealista-mcp/auth"
	"github.com/galimru/idealista-mcp/client"
	"github.com/galimru/idealista-mcp/internal/session"
	"github.com/galimru/idealista-mcp/internal/signing"
	"github.com/galimru/idealista-mcp/tools"
	"github.com/mark3labs/mcp-go/server"
)

// Populated by -ldflags during build; see Makefile.
var version = "dev"

func main() {
	sess, err := session.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "session error: %v\n", err)
		os.Exit(1)
	}

	authClient, err := auth.New(sess)
	if err != nil {
		fmt.Fprintf(os.Stderr, "auth error: %v\n", err)
		os.Exit(1)
	}

	signer, err := signing.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "signing error: %v\n", err)
		os.Exit(1)
	}

	apiClient := client.New(authClient, signer, sess.DeviceIdentifier)

	s := server.NewMCPServer(
		"idealista-mcp",
		version,
		server.WithToolCapabilities(false),
	)

	tools.RegisterLocationTools(s, apiClient)
	tools.RegisterSearchTools(s, apiClient)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
