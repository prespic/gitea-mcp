package operation

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitea.com/gitea/gitea-mcp/operation/issue"
	"gitea.com/gitea/gitea-mcp/operation/label"
	"gitea.com/gitea/gitea-mcp/operation/pull"
	"gitea.com/gitea/gitea-mcp/operation/repo"
	"gitea.com/gitea/gitea-mcp/operation/search"
	"gitea.com/gitea/gitea-mcp/operation/user"
	"gitea.com/gitea/gitea-mcp/operation/version"
	"gitea.com/gitea/gitea-mcp/operation/wiki"
	mcpContext "gitea.com/gitea/gitea-mcp/pkg/context"
	"gitea.com/gitea/gitea-mcp/pkg/flag"
	"gitea.com/gitea/gitea-mcp/pkg/log"

	"github.com/mark3labs/mcp-go/server"
)

var mcpServer *server.MCPServer

func RegisterTool(s *server.MCPServer) {
	// User Tool
	s.AddTools(user.Tool.Tools()...)

	// Repo Tool
	s.AddTools(repo.Tool.Tools()...)

	// Issue Tool
	s.AddTools(issue.Tool.Tools()...)

	// Label Tool
	s.AddTools(label.Tool.Tools()...)

	// Pull Tool
	s.AddTools(pull.Tool.Tools()...)

	// Search Tool
	s.AddTools(search.Tool.Tools()...)

	// Version Tool
	s.AddTools(version.Tool.Tools()...)

	// Wiki Tool
	s.AddTools(wiki.Tool.Tools()...)

	s.DeleteTools("")
}

// parseBearerToken extracts the Bearer token from an Authorization header.
// Returns the token and true if valid, empty string and false otherwise.
func parseBearerToken(authHeader string) (string, bool) {
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", false
	}

	token := strings.TrimSpace(authHeader[len(bearerPrefix):])
	if token == "" {
		return "", false
	}

	return token, true
}

func getContextWithToken(ctx context.Context, r *http.Request) context.Context {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ctx
	}

	token, ok := parseBearerToken(authHeader)
	if !ok {
		return ctx
	}

	return context.WithValue(ctx, mcpContext.TokenContextKey, token)
}

func Run() error {
	mcpServer = newMCPServer(flag.Version)
	RegisterTool(mcpServer)
	switch flag.Mode {
	case "stdio":
		if err := server.ServeStdio(
			mcpServer,
		); err != nil {
			return err
		}
	case "http":
		httpServer := server.NewStreamableHTTPServer(
			mcpServer,
			server.WithLogger(log.New()),
			server.WithHeartbeatInterval(30*time.Second),
			server.WithHTTPContextFunc(getContextWithToken),
		)
		log.Infof("Gitea MCP HTTP server listening on :%d", flag.Port)
		if err := httpServer.Start(fmt.Sprintf(":%d", flag.Port)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid transport type: %s. Must be 'stdio' or 'http'", flag.Mode)
	}
	return nil
}

func newMCPServer(version string) *server.MCPServer {
	return server.NewMCPServer(
		"Gitea MCP Server",
		version,
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithRecovery(),
	)
}

