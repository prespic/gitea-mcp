package repo

import (
	"context"
	"fmt"
	"time"

	"gitea.com/gitea/gitea-mcp/pkg/gitea"
	"gitea.com/gitea/gitea-mcp/pkg/log"
	"gitea.com/gitea/gitea-mcp/pkg/ptr"
	"gitea.com/gitea/gitea-mcp/pkg/to"

	gitea_sdk "code.gitea.io/sdk/gitea"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	CreateReleaseToolName    = "create_release"
	DeleteReleaseToolName    = "delete_release"
	GetReleaseToolName       = "get_release"
	GetLatestReleaseToolName = "get_latest_release"
	ListReleasesToolName     = "list_releases"
)

var (
	CreateReleaseTool = mcp.NewTool(
		CreateReleaseToolName,
		mcp.WithDescription("Create release"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("tag_name", mcp.Required(), mcp.Description("tag name")),
		mcp.WithString("target", mcp.Required(), mcp.Description("target commitish")),
		mcp.WithString("title", mcp.Required(), mcp.Description("release title")),
		mcp.WithBoolean("is_draft", mcp.Description("Whether the release is draft"), mcp.DefaultBool(false)),
		mcp.WithBoolean("is_pre_release", mcp.Description("Whether the release is pre-release"), mcp.DefaultBool(false)),
		mcp.WithString("body", mcp.Description("release body")),
	)

	DeleteReleaseTool = mcp.NewTool(
		DeleteReleaseToolName,
		mcp.WithDescription("Delete release"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("release id")),
	)

	GetReleaseTool = mcp.NewTool(
		GetReleaseToolName,
		mcp.WithDescription("Get release"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("release id")),
	)

	GetLatestReleaseTool = mcp.NewTool(
		GetLatestReleaseToolName,
		mcp.WithDescription("Get latest release"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
	)

	ListReleasesTool = mcp.NewTool(
		ListReleasesToolName,
		mcp.WithDescription("List releases"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithBoolean("is_draft", mcp.Description("Whether the release is draft"), mcp.DefaultBool(false)),
		mcp.WithBoolean("is_pre_release", mcp.Description("Whether the release is pre-release"), mcp.DefaultBool(false)),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1), mcp.Min(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(20), mcp.Min(1)),
	)
)

func init() {
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreateReleaseTool,
		Handler: CreateReleaseFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    DeleteReleaseTool,
		Handler: DeleteReleaseFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetReleaseTool,
		Handler: GetReleaseFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetLatestReleaseTool,
		Handler: GetLatestReleaseFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    ListReleasesTool,
		Handler: ListReleasesFn,
	})
}

// To avoid return too many tokens, we need to provide at least information as possible
// llm can call get release to get more information
type ListReleaseResult struct {
	ID           int64     `json:"id"`
	TagName      string    `json:"tag_name"`
	Target       string    `json:"target_commitish"`
	Title        string    `json:"title"`
	IsDraft      bool      `json:"draft"`
	IsPrerelease bool      `json:"prerelease"`
	CreatedAt    time.Time `json:"created_at"`
	PublishedAt  time.Time `json:"published_at"`
}

func CreateReleaseFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateReleasesFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner is required")
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo is required")
	}
	tagName, ok := req.GetArguments()["tag_name"].(string)
	if !ok {
		return nil, fmt.Errorf("tag_name is required")
	}
	target, ok := req.GetArguments()["target"].(string)
	if !ok {
		return nil, fmt.Errorf("target is required")
	}
	title, ok := req.GetArguments()["title"].(string)
	if !ok {
		return nil, fmt.Errorf("title is required")
	}
	isDraft, _ := req.GetArguments()["is_draft"].(bool)
	isPreRelease, _ := req.GetArguments()["is_pre_release"].(bool)
	body, _ := req.GetArguments()["body"].(string)

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	_, _, err = client.CreateRelease(owner, repo, gitea_sdk.CreateReleaseOption{
		TagName:      tagName,
		Target:       target,
		Title:        title,
		Note:         body,
		IsDraft:      isDraft,
		IsPrerelease: isPreRelease,
	})
	if err != nil {
		return nil, fmt.Errorf("create release error: %v", err)
	}

	return mcp.NewToolResultText("Release Created"), nil
}

func DeleteReleaseFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteReleaseFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner is required")
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo is required")
	}
	id, ok := req.GetArguments()["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("id is required")
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	_, err = client.DeleteRelease(owner, repo, int64(id))
	if err != nil {
		return nil, fmt.Errorf("delete release error: %v", err)
	}

	return to.TextResult("Release deleted successfully")
}

func GetReleaseFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetReleaseFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner is required")
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo is required")
	}
	id, ok := req.GetArguments()["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("id is required")
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	release, _, err := client.GetRelease(owner, repo, int64(id))
	if err != nil {
		return nil, fmt.Errorf("get release error: %v", err)
	}

	return to.TextResult(release)
}

func GetLatestReleaseFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetLatestReleaseFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner is required")
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo is required")
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	release, _, err := client.GetLatestRelease(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("get latest release error: %v", err)
	}

	return to.TextResult(release)
}

func ListReleasesFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListReleasesFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return nil, fmt.Errorf("owner is required")
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return nil, fmt.Errorf("repo is required")
	}
	var pIsDraft *bool
	isDraft, ok := req.GetArguments()["is_draft"].(bool)
	if ok {
		pIsDraft = ptr.To(isDraft)
	}
	var pIsPreRelease *bool
	isPreRelease, ok := req.GetArguments()["is_pre_release"].(bool)
	if ok {
		pIsPreRelease = ptr.To(isPreRelease)
	}
	page, _ := req.GetArguments()["page"].(float64)
	pageSize, _ := req.GetArguments()["pageSize"].(float64)

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	releases, _, err := client.ListReleases(owner, repo, gitea_sdk.ListReleasesOptions{
		ListOptions: gitea_sdk.ListOptions{
			Page:     int(page),
			PageSize: int(pageSize),
		},
		IsDraft:      pIsDraft,
		IsPreRelease: pIsPreRelease,
	})
	if err != nil {
		return nil, fmt.Errorf("list releases error: %v", err)
	}

	results := make([]ListReleaseResult, len(releases))
	for _, release := range releases {
		results = append(results, ListReleaseResult{
			ID:           release.ID,
			TagName:      release.TagName,
			Target:       release.Target,
			Title:        release.Title,
			IsDraft:      release.IsDraft,
			IsPrerelease: release.IsPrerelease,
			CreatedAt:    release.CreatedAt,
			PublishedAt:  release.PublishedAt,
		})
	}
	return to.TextResult(results)
}
