package actions

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"gitea.com/gitea/gitea-mcp/pkg/gitea"
	"gitea.com/gitea/gitea-mcp/pkg/log"
	"gitea.com/gitea/gitea-mcp/pkg/params"
	"gitea.com/gitea/gitea-mcp/pkg/to"

	gitea_sdk "code.gitea.io/sdk/gitea"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	ListRepoActionSecretsToolName  = "list_repo_action_secrets"
	UpsertRepoActionSecretToolName = "upsert_repo_action_secret"
	DeleteRepoActionSecretToolName = "delete_repo_action_secret"
	ListOrgActionSecretsToolName   = "list_org_action_secrets"
	UpsertOrgActionSecretToolName  = "upsert_org_action_secret"
	DeleteOrgActionSecretToolName  = "delete_org_action_secret"
)

type secretMeta struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitzero"`
}

var (
	ListRepoActionSecretsTool = mcp.NewTool(
		ListRepoActionSecretsToolName,
		mcp.WithDescription("List repository Actions secrets (metadata only; secret values are never returned)"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1), mcp.Min(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(30), mcp.Min(1)),
	)

	UpsertRepoActionSecretTool = mcp.NewTool(
		UpsertRepoActionSecretToolName,
		mcp.WithDescription("Create or update (upsert) a repository Actions secret"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("name", mcp.Required(), mcp.Description("secret name")),
		mcp.WithString("data", mcp.Required(), mcp.Description("secret value")),
		mcp.WithString("description", mcp.Description("secret description")),
	)

	DeleteRepoActionSecretTool = mcp.NewTool(
		DeleteRepoActionSecretToolName,
		mcp.WithDescription("Delete a repository Actions secret"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("secretName", mcp.Required(), mcp.Description("secret name")),
	)

	ListOrgActionSecretsTool = mcp.NewTool(
		ListOrgActionSecretsToolName,
		mcp.WithDescription("List organization Actions secrets (metadata only; secret values are never returned)"),
		mcp.WithString("org", mcp.Required(), mcp.Description("organization name")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1), mcp.Min(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(30), mcp.Min(1)),
	)

	UpsertOrgActionSecretTool = mcp.NewTool(
		UpsertOrgActionSecretToolName,
		mcp.WithDescription("Create or update (upsert) an organization Actions secret"),
		mcp.WithString("org", mcp.Required(), mcp.Description("organization name")),
		mcp.WithString("name", mcp.Required(), mcp.Description("secret name")),
		mcp.WithString("data", mcp.Required(), mcp.Description("secret value")),
		mcp.WithString("description", mcp.Description("secret description")),
	)

	DeleteOrgActionSecretTool = mcp.NewTool(
		DeleteOrgActionSecretToolName,
		mcp.WithDescription("Delete an organization Actions secret"),
		mcp.WithString("org", mcp.Required(), mcp.Description("organization name")),
		mcp.WithString("secretName", mcp.Required(), mcp.Description("secret name")),
	)
)

func init() {
	Tool.RegisterRead(server.ServerTool{Tool: ListRepoActionSecretsTool, Handler: ListRepoActionSecretsFn})
	Tool.RegisterWrite(server.ServerTool{Tool: UpsertRepoActionSecretTool, Handler: UpsertRepoActionSecretFn})
	Tool.RegisterWrite(server.ServerTool{Tool: DeleteRepoActionSecretTool, Handler: DeleteRepoActionSecretFn})

	Tool.RegisterRead(server.ServerTool{Tool: ListOrgActionSecretsTool, Handler: ListOrgActionSecretsFn})
	Tool.RegisterWrite(server.ServerTool{Tool: UpsertOrgActionSecretTool, Handler: UpsertOrgActionSecretFn})
	Tool.RegisterWrite(server.ServerTool{Tool: DeleteOrgActionSecretTool, Handler: DeleteOrgActionSecretFn})
}

func ListRepoActionSecretsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListRepoActionSecretsFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	page, pageSize := params.GetPagination(req.GetArguments(), 30)

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	secrets, _, err := client.ListRepoActionSecret(owner, repo, gitea_sdk.ListRepoActionSecretOption{
		ListOptions: gitea_sdk.ListOptions{Page: page, PageSize: pageSize},
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list repo action secrets err: %v", err))
	}

	metas := make([]secretMeta, 0, len(secrets))
	for _, s := range secrets {
		if s == nil {
			continue
		}
		metas = append(metas, secretMeta{
			Name:        s.Name,
			Description: s.Description,
			CreatedAt:   s.Created,
		})
	}
	return to.TextResult(metas)
}

func UpsertRepoActionSecretFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called UpsertRepoActionSecretFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	name, err := params.GetString(req.GetArguments(), "name")
	if err != nil || name == "" {
		return to.ErrorResult(errors.New("name is required"))
	}
	data, err := params.GetString(req.GetArguments(), "data")
	if err != nil || data == "" {
		return to.ErrorResult(errors.New("data is required"))
	}
	description, _ := req.GetArguments()["description"].(string)

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	resp, err := client.CreateRepoActionSecret(owner, repo, gitea_sdk.CreateSecretOption{
		Name:        name,
		Data:        data,
		Description: description,
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("upsert repo action secret err: %v", err))
	}
	return to.TextResult(map[string]any{"message": "secret upserted", "status": resp.StatusCode})
}

func DeleteRepoActionSecretFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteRepoActionSecretFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	secretName, err := params.GetString(req.GetArguments(), "secretName")
	if err != nil || secretName == "" {
		return to.ErrorResult(errors.New("secretName is required"))
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	resp, err := client.DeleteRepoActionSecret(owner, repo, secretName)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete repo action secret err: %v", err))
	}
	return to.TextResult(map[string]any{"message": "secret deleted", "status": resp.StatusCode})
}

func ListOrgActionSecretsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListOrgActionSecretsFn")
	org, err := params.GetString(req.GetArguments(), "org")
	if err != nil || org == "" {
		return to.ErrorResult(errors.New("org is required"))
	}
	page, pageSize := params.GetPagination(req.GetArguments(), 30)

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	secrets, _, err := client.ListOrgActionSecret(org, gitea_sdk.ListOrgActionSecretOption{
		ListOptions: gitea_sdk.ListOptions{Page: page, PageSize: pageSize},
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list org action secrets err: %v", err))
	}

	metas := make([]secretMeta, 0, len(secrets))
	for _, s := range secrets {
		if s == nil {
			continue
		}
		metas = append(metas, secretMeta{
			Name:        s.Name,
			Description: s.Description,
			CreatedAt:   s.Created,
		})
	}
	return to.TextResult(metas)
}

func UpsertOrgActionSecretFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called UpsertOrgActionSecretFn")
	org, err := params.GetString(req.GetArguments(), "org")
	if err != nil || org == "" {
		return to.ErrorResult(errors.New("org is required"))
	}
	name, err := params.GetString(req.GetArguments(), "name")
	if err != nil || name == "" {
		return to.ErrorResult(errors.New("name is required"))
	}
	data, err := params.GetString(req.GetArguments(), "data")
	if err != nil || data == "" {
		return to.ErrorResult(errors.New("data is required"))
	}
	description, _ := req.GetArguments()["description"].(string)

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	resp, err := client.CreateOrgActionSecret(org, gitea_sdk.CreateSecretOption{
		Name:        name,
		Data:        data,
		Description: description,
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("upsert org action secret err: %v", err))
	}
	return to.TextResult(map[string]any{"message": "secret upserted", "status": resp.StatusCode})
}

func DeleteOrgActionSecretFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteOrgActionSecretFn")
	org, err := params.GetString(req.GetArguments(), "org")
	if err != nil || org == "" {
		return to.ErrorResult(errors.New("org is required"))
	}
	secretName, err := params.GetString(req.GetArguments(), "secretName")
	if err != nil || secretName == "" {
		return to.ErrorResult(errors.New("secretName is required"))
	}

	escapedOrg := url.PathEscape(org)
	escapedSecret := url.PathEscape(secretName)
	_, err = gitea.DoJSON(ctx, "DELETE", fmt.Sprintf("orgs/%s/actions/secrets/%s", escapedOrg, escapedSecret), nil, nil, nil)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete org action secret err: %v", err))
	}
	return to.TextResult(map[string]any{"message": "secret deleted"})
}
