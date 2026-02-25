package search

import (
	"context"
	"errors"
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
	SearchUsersToolName    = "search_users"
	SearchOrgTeamsToolName = "search_org_teams"
	SearchReposToolName    = "search_repos"
)

var (
	SearchUsersTool = mcp.NewTool(
		SearchUsersToolName,
		mcp.WithDescription("search users"),
		mcp.WithString("keyword", mcp.Required(), mcp.Description("Keyword")),
		mcp.WithNumber("page", mcp.Description("Page"), mcp.DefaultNumber(1)),
		mcp.WithNumber("pageSize", mcp.Description("PageSize"), mcp.DefaultNumber(100)),
	)

	SearOrgTeamsTool = mcp.NewTool(
		SearchOrgTeamsToolName,
		mcp.WithDescription("search organization teams"),
		mcp.WithString("org", mcp.Required(), mcp.Description("organization name")),
		mcp.WithString("query", mcp.Required(), mcp.Description("search organization teams")),
		mcp.WithBoolean("includeDescription", mcp.Description("include description?")),
		mcp.WithNumber("page", mcp.Description("Page"), mcp.DefaultNumber(1)),
		mcp.WithNumber("pageSize", mcp.Description("PageSize"), mcp.DefaultNumber(100)),
	)

	SearchReposTool = mcp.NewTool(
		SearchReposToolName,
		mcp.WithDescription("search repos"),
		mcp.WithString("keyword", mcp.Required(), mcp.Description("Keyword")),
		mcp.WithBoolean("keywordIsTopic", mcp.Description("KeywordIsTopic")),
		mcp.WithBoolean("keywordInDescription", mcp.Description("KeywordInDescription")),
		mcp.WithNumber("ownerID", mcp.Description("OwnerID")),
		mcp.WithBoolean("isPrivate", mcp.Description("IsPrivate")),
		mcp.WithBoolean("isArchived", mcp.Description("IsArchived")),
		mcp.WithString("sort", mcp.Description("Sort")),
		mcp.WithString("order", mcp.Description("Order")),
		mcp.WithNumber("page", mcp.Description("Page"), mcp.DefaultNumber(1)),
		mcp.WithNumber("pageSize", mcp.Description("PageSize"), mcp.DefaultNumber(100)),
	)
)

func init() {
	Tool.RegisterRead(server.ServerTool{
		Tool:    SearchUsersTool,
		Handler: UsersFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    SearOrgTeamsTool,
		Handler: OrgTeamsFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    SearchReposTool,
		Handler: ReposFn,
	})
}

func UsersFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called UsersFn")
	keyword, ok := req.GetArguments()["keyword"].(string)
	if !ok {
		return to.ErrorResult(errors.New("keyword is required"))
	}
	page := params.GetOptionalInt(req.GetArguments(), "page", 1)
	pageSize := params.GetOptionalInt(req.GetArguments(), "pageSize", 100)
	opt := gitea_sdk.SearchUsersOption{
		KeyWord: keyword,
		ListOptions: gitea_sdk.ListOptions{
			Page:     int(page),
			PageSize: int(pageSize),
		},
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	users, _, err := client.SearchUsers(opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("search users err: %v", err))
	}
	return to.TextResult(users)
}

func OrgTeamsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called OrgTeamsFn")
	org, ok := req.GetArguments()["org"].(string)
	if !ok {
		return to.ErrorResult(errors.New("organization is required"))
	}
	query, ok := req.GetArguments()["query"].(string)
	if !ok {
		return to.ErrorResult(errors.New("query is required"))
	}
	includeDescription, _ := req.GetArguments()["includeDescription"].(bool)
	page := params.GetOptionalInt(req.GetArguments(), "page", 1)
	pageSize := params.GetOptionalInt(req.GetArguments(), "pageSize", 100)
	opt := gitea_sdk.SearchTeamsOptions{
		Query:              query,
		IncludeDescription: includeDescription,
		ListOptions: gitea_sdk.ListOptions{
			Page:     int(page),
			PageSize: int(pageSize),
		},
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	teams, _, err := client.SearchOrgTeams(org, &opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("search organization teams error: %v", err))
	}
	return to.TextResult(teams)
}

func ReposFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ReposFn")
	keyword, ok := req.GetArguments()["keyword"].(string)
	if !ok {
		return to.ErrorResult(errors.New("keyword is required"))
	}
	keywordIsTopic, _ := req.GetArguments()["keywordIsTopic"].(bool)
	keywordInDescription, _ := req.GetArguments()["keywordInDescription"].(bool)
	ownerID := params.GetOptionalInt(req.GetArguments(), "ownerID", 0)
	var pIsPrivate *bool
	isPrivate, ok := req.GetArguments()["isPrivate"].(bool)
	if ok {
		pIsPrivate = new(isPrivate)
	}
	var pIsArchived *bool
	isArchived, ok := req.GetArguments()["isArchived"].(bool)
	if ok {
		pIsArchived = new(isArchived)
	}
	sort, _ := req.GetArguments()["sort"].(string)
	order, _ := req.GetArguments()["order"].(string)
	page := params.GetOptionalInt(req.GetArguments(), "page", 1)
	pageSize := params.GetOptionalInt(req.GetArguments(), "pageSize", 100)
	opt := gitea_sdk.SearchRepoOptions{
		Keyword:              keyword,
		KeywordIsTopic:       keywordIsTopic,
		KeywordInDescription: keywordInDescription,
		OwnerID:              ownerID,
		IsPrivate:            pIsPrivate,
		IsArchived:           pIsArchived,
		Sort:                 sort,
		Order:                order,
		ListOptions: gitea_sdk.ListOptions{
			Page:     int(page),
			PageSize: int(pageSize),
		},
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	repos, _, err := client.SearchRepos(opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("search repos error: %v", err))
	}
	return to.TextResult(repos)
}
