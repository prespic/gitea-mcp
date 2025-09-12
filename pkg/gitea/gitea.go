package gitea

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"code.gitea.io/sdk/gitea"
	mcpContext "gitea.com/gitea/gitea-mcp/pkg/context"
	"gitea.com/gitea/gitea-mcp/pkg/flag"
)

func NewClient(token string) (*gitea.Client, error) {
	httpClient := &http.Client{
		Transport: http.DefaultTransport,
	}

	opts := []gitea.ClientOption{
		gitea.SetToken(token),
	}
	if flag.Insecure {
		httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	opts = append(opts, gitea.SetHTTPClient(httpClient))
	if flag.Debug {
		opts = append(opts, gitea.SetDebugMode())
	}
	client, err := gitea.NewClient(flag.Host, opts...)
	if err != nil {
		return nil, fmt.Errorf("create gitea client err: %w", err)
	}

	// Set user agent for the client
	client.SetUserAgent(fmt.Sprintf("gitea-mcp-server/%s", flag.Version))
	return client, nil
}

func ClientFromContext(ctx context.Context) (*gitea.Client, error) {
	token, ok := ctx.Value(mcpContext.TokenContextKey).(string)
	if !ok {
		token = flag.Token
	}
	return NewClient(token)
}
