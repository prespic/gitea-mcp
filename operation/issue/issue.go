package issue

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
	GetIssueByIndexToolName         = "get_issue_by_index"
	ListRepoIssuesToolName          = "list_repo_issues"
	CreateIssueToolName             = "create_issue"
	CreateIssueCommentToolName      = "create_issue_comment"
	EditIssueToolName               = "edit_issue"
	EditIssueCommentToolName        = "edit_issue_comment"
	GetIssueCommentsByIndexToolName = "get_issue_comments_by_index"
)

var (
	GetIssueByIndexTool = mcp.NewTool(
		GetIssueByIndexToolName,
		mcp.WithDescription("get issue by index"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("repository issue index")),
	)

	ListRepoIssuesTool = mcp.NewTool(
		ListRepoIssuesToolName,
		mcp.WithDescription("List repository issues"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("state", mcp.Description("issue state"), mcp.DefaultString("all")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(30)),
	)

	CreateIssueTool = mcp.NewTool(
		CreateIssueToolName,
		mcp.WithDescription("create issue"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("title", mcp.Required(), mcp.Description("issue title")),
		mcp.WithString("body", mcp.Required(), mcp.Description("issue body")),
	)

	CreateIssueCommentTool = mcp.NewTool(
		CreateIssueCommentToolName,
		mcp.WithDescription("create issue comment"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("repository issue index")),
		mcp.WithString("body", mcp.Required(), mcp.Description("issue comment body")),
	)

	EditIssueTool = mcp.NewTool(
		EditIssueToolName,
		mcp.WithDescription("edit issue"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("repository issue index")),
		mcp.WithString("title", mcp.Description("issue title"), mcp.DefaultString("")),
		mcp.WithString("body", mcp.Description("issue body content")),
		mcp.WithArray("assignees", mcp.Description("usernames to assign to this issue"), mcp.Items(map[string]any{"type": "string"})),
		mcp.WithNumber("milestone", mcp.Description("milestone number")),
		mcp.WithString("state", mcp.Description("issue state, one of open, closed, all")),
	)

	EditIssueCommentTool = mcp.NewTool(
		EditIssueCommentToolName,
		mcp.WithDescription("edit issue comment"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("commentID", mcp.Required(), mcp.Description("id of issue comment")),
		mcp.WithString("body", mcp.Required(), mcp.Description("issue comment body")),
	)

	GetIssueCommentsByIndexTool = mcp.NewTool(
		GetIssueCommentsByIndexToolName,
		mcp.WithDescription("get issue comment by index"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("repository issue index")),
	)
)

func init() {
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetIssueByIndexTool,
		Handler: GetIssueByIndexFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    ListRepoIssuesTool,
		Handler: ListRepoIssuesFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreateIssueTool,
		Handler: CreateIssueFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreateIssueCommentTool,
		Handler: CreateIssueCommentFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    EditIssueTool,
		Handler: EditIssueFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    EditIssueCommentTool,
		Handler: EditIssueCommentFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetIssueCommentsByIndexTool,
		Handler: GetIssueCommentsByIndexFn,
	})
}

func GetIssueByIndexFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetIssueByIndexFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(req.GetArguments(), "index")
	if err != nil {
		return to.ErrorResult(err)
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	issue, _, err := client.GetIssue(owner, repo, index)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/issue/%v err: %v", owner, repo, index, err))
	}

	return to.TextResult(slimIssue(issue))
}

func ListRepoIssuesFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListIssuesFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	state, ok := req.GetArguments()["state"].(string)
	if !ok {
		state = "all"
	}
	page, pageSize := params.GetPagination(req.GetArguments(), 30)
	opt := gitea_sdk.ListIssueOption{
		State: gitea_sdk.StateType(state),
		ListOptions: gitea_sdk.ListOptions{
			Page:     page,
			PageSize: pageSize,
		},
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	issues, _, err := client.ListRepoIssues(owner, repo, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/issues err: %v", owner, repo, err))
	}
	return to.TextResult(slimIssues(issues))
}

func CreateIssueFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateIssueFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	title, err := params.GetString(req.GetArguments(), "title")
	if err != nil {
		return to.ErrorResult(err)
	}
	body, err := params.GetString(req.GetArguments(), "body")
	if err != nil {
		return to.ErrorResult(err)
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	issue, _, err := client.CreateIssue(owner, repo, gitea_sdk.CreateIssueOption{
		Title: title,
		Body:  body,
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create %v/%v/issue err: %v", owner, repo, err))
	}

	return to.TextResult(slimIssue(issue))
}

func CreateIssueCommentFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateIssueCommentFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(req.GetArguments(), "index")
	if err != nil {
		return to.ErrorResult(err)
	}
	body, err := params.GetString(req.GetArguments(), "body")
	if err != nil {
		return to.ErrorResult(err)
	}
	opt := gitea_sdk.CreateIssueCommentOption{
		Body: body,
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	issueComment, _, err := client.CreateIssueComment(owner, repo, index, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create %v/%v/issue/%v/comment err: %v", owner, repo, index, err))
	}

	return to.TextResult(slimComment(issueComment))
}

func EditIssueFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called EditIssueFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(req.GetArguments(), "index")
	if err != nil {
		return to.ErrorResult(err)
	}

	opt := gitea_sdk.EditIssueOption{}

	title, ok := req.GetArguments()["title"].(string)
	if ok {
		opt.Title = title
	}
	body, ok := req.GetArguments()["body"].(string)
	if ok {
		opt.Body = new(body)
	}
	opt.Assignees = params.GetStringSlice(req.GetArguments(), "assignees")
	if val, exists := req.GetArguments()["milestone"]; exists {
		if milestone, ok := params.ToInt64(val); ok {
			opt.Milestone = new(milestone)
		}
	}
	state, ok := req.GetArguments()["state"].(string)
	if ok {
		opt.State = new(gitea_sdk.StateType(state))
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	issue, _, err := client.EditIssue(owner, repo, index, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("edit %v/%v/issue/%v err: %v", owner, repo, index, err))
	}

	return to.TextResult(slimIssue(issue))
}

func EditIssueCommentFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called EditIssueCommentFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	commentID, err := params.GetIndex(req.GetArguments(), "commentID")
	if err != nil {
		return to.ErrorResult(err)
	}
	body, err := params.GetString(req.GetArguments(), "body")
	if err != nil {
		return to.ErrorResult(err)
	}
	opt := gitea_sdk.EditIssueCommentOption{
		Body: body,
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	issueComment, _, err := client.EditIssueComment(owner, repo, commentID, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("edit %v/%v/issues/comments/%v err: %v", owner, repo, commentID, err))
	}

	return to.TextResult(slimComment(issueComment))
}

func GetIssueCommentsByIndexFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetIssueCommentsByIndexFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(req.GetArguments(), "index")
	if err != nil {
		return to.ErrorResult(err)
	}
	opt := gitea_sdk.ListIssueCommentOptions{}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	issue, _, err := client.ListIssueComments(owner, repo, index, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/issues/%v/comments err: %v", owner, repo, index, err))
	}

	return to.TextResult(slimComments(issue))
}
