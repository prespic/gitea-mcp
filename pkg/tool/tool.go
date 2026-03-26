package tool

import (
	"gitea.com/gitea/gitea-mcp/pkg/flag"
	"github.com/mark3labs/mcp-go/server"
)

type Tool struct {
	write []server.ServerTool
	read  []server.ServerTool
}

func New() *Tool {
	return &Tool{
		write: make([]server.ServerTool, 0, 100),
		read:  make([]server.ServerTool, 0, 100),
	}
}

func (t *Tool) RegisterWrite(s server.ServerTool) {
	t.write = append(t.write, s)
}

func (t *Tool) RegisterRead(s server.ServerTool) {
	t.read = append(t.read, s)
}

func (t *Tool) Tools() []server.ServerTool {
	tools := make([]server.ServerTool, 0, len(t.write)+len(t.read))
	if flag.ReadOnly {
		tools = append(tools, t.read...)
	} else {
		tools = append(tools, t.write...)
		tools = append(tools, t.read...)
	}
	if len(flag.AllowedTools) == 0 {
		return tools
	}
	allowed := make(map[string]bool, len(flag.AllowedTools))
	for _, name := range flag.AllowedTools {
		allowed[name] = true
	}
	filtered := make([]server.ServerTool, 0, len(allowed))
	for _, st := range tools {
		if allowed[st.Tool.Name] {
			filtered = append(filtered, st)
		}
	}
	return filtered
}
