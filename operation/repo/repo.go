package repo

import (
	"context"
	"fmt"

	"gitea.com/gitea/gitea-mcp/pkg/gitea"
	"gitea.com/gitea/gitea-mcp/pkg/log"
	"gitea.com/gitea/gitea-mcp/pkg/params"
	"gitea.com/gitea/gitea-mcp/pkg/to"
	"gitea.com/gitea/gitea-mcp/pkg/tool"

	gitea_sdk "code.gitea.io/sdk/gitea"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var Tool = tool.New()

const (
	CreateRepoToolName  = "create_repo"
	ForkRepoToolName    = "fork_repo"
	ListMyReposToolName = "list_my_repos"
)

var (
	CreateRepoTool = mcp.NewTool(
		CreateRepoToolName,
		mcp.WithDescription("Create repository in personal account or organization"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the repository to create")),
		mcp.WithString("description", mcp.Description("Description of the repository to create")),
		mcp.WithBoolean("private", mcp.Description("Whether the repository is private")),
		mcp.WithString("issue_labels", mcp.Description("Issue Label set to use")),
		mcp.WithBoolean("auto_init", mcp.Description("Whether the repository should be auto-intialized?")),
		mcp.WithBoolean("template", mcp.Description("Whether the repository is template")),
		mcp.WithString("gitignores", mcp.Description("Gitignores to use")),
		mcp.WithString("license", mcp.Description("License to use")),
		mcp.WithString("readme", mcp.Description("Readme of the repository to create")),
		mcp.WithString("default_branch", mcp.Description("DefaultBranch of the repository (used when initializes and in template)")),
		mcp.WithString("organization", mcp.Description("Organization name to create repository in (optional - defaults to personal account)")),
	)

	ForkRepoTool = mcp.NewTool(
		ForkRepoToolName,
		mcp.WithDescription("Fork repository"),
		mcp.WithString("user", mcp.Required(), mcp.Description("User name of the repository to fork")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("Repository name to fork")),
		mcp.WithString("organization", mcp.Description("Organization name to fork")),
		mcp.WithString("name", mcp.Description("Name of the forked repository")),
	)

	ListMyReposTool = mcp.NewTool(
		ListMyReposToolName,
		mcp.WithDescription("List my repositories"),
		mcp.WithNumber("page", mcp.Required(), mcp.Description("Page number"), mcp.DefaultNumber(1), mcp.Min(1)),
		mcp.WithNumber("pageSize", mcp.Required(), mcp.Description("Page size number"), mcp.DefaultNumber(30), mcp.Min(1)),
	)
)

func init() {
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreateRepoTool,
		Handler: CreateRepoFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    ForkRepoTool,
		Handler: ForkRepoFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    ListMyReposTool,
		Handler: ListMyReposFn,
	})
}

func RegisterTool(s *server.MCPServer) {
	s.AddTool(CreateRepoTool, CreateRepoFn)
	s.AddTool(ForkRepoTool, ForkRepoFn)
	s.AddTool(ListMyReposTool, ListMyReposFn)

	// File
	s.AddTool(GetFileContentTool, GetFileContentFn)
	s.AddTool(CreateFileTool, CreateFileFn)
	s.AddTool(UpdateFileTool, UpdateFileFn)
	s.AddTool(DeleteFileTool, DeleteFileFn)

	// Branch
	s.AddTool(CreateBranchTool, CreateBranchFn)
	s.AddTool(DeleteBranchTool, DeleteBranchFn)
	s.AddTool(ListBranchesTool, ListBranchesFn)

	// Release
	s.AddTool(CreateReleaseTool, CreateReleaseFn)
	s.AddTool(DeleteReleaseTool, DeleteReleaseFn)
	s.AddTool(GetReleaseTool, GetReleaseFn)
	s.AddTool(GetLatestReleaseTool, GetLatestReleaseFn)
	s.AddTool(ListReleasesTool, ListReleasesFn)

	// Tag
	s.AddTool(CreateTagTool, CreateTagFn)
	s.AddTool(DeleteTagTool, DeleteTagFn)
	s.AddTool(GetTagTool, GetTagFn)
	s.AddTool(ListTagsTool, ListTagsFn)

	// Commit
	s.AddTool(ListRepoCommitsTool, ListRepoCommitsFn)
}

func CreateRepoFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateRepoFn")
	args := req.GetArguments()
	name, err := params.GetString(args, "name")
	if err != nil {
		return to.ErrorResult(err)
	}
	description, _ := args["description"].(string)
	private, _ := args["private"].(bool)
	issueLabels, _ := args["issue_labels"].(string)
	autoInit, _ := args["auto_init"].(bool)
	template, _ := args["template"].(bool)
	gitignores, _ := args["gitignores"].(string)
	license, _ := args["license"].(string)
	readme, _ := args["readme"].(string)
	defaultBranch, _ := args["default_branch"].(string)
	organization, _ := args["organization"].(string)

	opt := gitea_sdk.CreateRepoOption{
		Name:          name,
		Description:   description,
		Private:       private,
		IssueLabels:   issueLabels,
		AutoInit:      autoInit,
		Template:      template,
		Gitignores:    gitignores,
		License:       license,
		Readme:        readme,
		DefaultBranch: defaultBranch,
	}

	var repo *gitea_sdk.Repository
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	if organization != "" {
		repo, _, err = client.CreateOrgRepo(organization, opt)
		if err != nil {
			return to.ErrorResult(fmt.Errorf("create organization repository '%s' in '%s' err: %v", name, organization, err))
		}
	} else {
		repo, _, err = client.CreateRepo(opt)
		if err != nil {
			return to.ErrorResult(fmt.Errorf("create repository '%s' err: %v", name, err))
		}
	}
	return to.TextResult(slimRepo(repo))
}

func ForkRepoFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ForkRepoFn")
	args := req.GetArguments()
	user, err := params.GetString(args, "user")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	organization, ok := args["organization"].(string)
	organizationPtr := new(organization)
	if !ok || organization == "" {
		organizationPtr = nil
	}
	name, ok := args["name"].(string)
	namePtr := new(name)
	if !ok || name == "" {
		namePtr = nil
	}
	opt := gitea_sdk.CreateForkOption{
		Organization: organizationPtr,
		Name:         namePtr,
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	_, _, err = client.CreateFork(user, repo, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("fork repository error: %v", err))
	}
	return to.TextResult("Fork success")
}

func ListMyReposFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListMyReposFn")
	page, pageSize := params.GetPagination(req.GetArguments(), 30)
	opt := gitea_sdk.ListReposOptions{
		ListOptions: gitea_sdk.ListOptions{
			Page:     page,
			PageSize: pageSize,
		},
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	repos, _, err := client.ListMyRepos(opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list my repositories error: %v", err))
	}

	return to.TextResult(slimRepos(repos))
}
