package cmd

import (
	"context"
	"flag"
	"os"

	"gitea.com/gitea/gitea-mcp/operation"
	flagPkg "gitea.com/gitea/gitea-mcp/pkg/flag"
	"gitea.com/gitea/gitea-mcp/pkg/log"
)

var (
	host  string
	port  int
	token string
)

func init() {
	flag.StringVar(
		&flagPkg.Mode,
		"t",
		"stdio",
		"Transport type (stdio or http)",
	)
	flag.StringVar(
		&flagPkg.Mode,
		"transport",
		"stdio",
		"Transport type (stdio or http)",
	)
	flag.StringVar(
		&host,
		"host",
		os.Getenv("GITEA_HOST"),
		"Gitea host",
	)
	flag.IntVar(
		&port,
		"port",
		8080,
		"http port",
	)
	flag.StringVar(
		&token,
		"token",
		"",
		"Your personal access token",
	)
	flag.BoolVar(
		&flagPkg.ReadOnly,
		"read-only",
		false,
		"Read-only mode",
	)
	flag.BoolVar(
		&flagPkg.Debug,
		"d",
		false,
		"debug mode (If -d flag is provided, debug mode will be enabled by default)",
	)
	flag.BoolVar(
		&flagPkg.Insecure,
		"insecure",
		false,
		"ignore TLS certificate errors",
	)

	flag.Parse()

	flagPkg.Host = host
	if flagPkg.Host == "" {
		flagPkg.Host = "https://gitea.com"
	}

	flagPkg.Port = port

	flagPkg.Token = token
	if flagPkg.Token == "" {
		flagPkg.Token = os.Getenv("GITEA_ACCESS_TOKEN")
	}

	if os.Getenv("MCP_MODE") != "" {
		flagPkg.Mode = os.Getenv("MCP_MODE")
	}

	if os.Getenv("GITEA_READONLY") == "true" {
		flagPkg.ReadOnly = true
	}

	if os.Getenv("GITEA_DEBUG") == "true" {
		flagPkg.Debug = true
	}

	// Set insecure mode based on environment variable
	if os.Getenv("GITEA_INSECURE") == "true" {
		flagPkg.Insecure = true
	}
}

func Execute() {
	defer log.Default().Sync()
	if err := operation.Run(); err != nil {
		if err == context.Canceled {
			log.Info("Server shutdown due to context cancellation")
			return
		}
		log.Fatalf("Run Gitea MCP Server Error: %v", err)
	}
}
