package milestone

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
	GetMilestoneToolName    = "get_milestone"
	ListMilestonesToolName  = "list_milestones"
	CreateMilestoneToolName = "create_milestone"
	EditMilestoneToolName   = "edit_milestone"
	DeleteMilestoneToolName = "delete_milestone"
)

var (
	GetMilestoneTool = mcp.NewTool(
		GetMilestoneToolName,
		mcp.WithDescription("get milestone by id"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("milestone id")),
	)

	ListMilestonesTool = mcp.NewTool(
		ListMilestonesToolName,
		mcp.WithDescription("List milestones"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("state", mcp.Description("milestone state"), mcp.DefaultString("all")),
		mcp.WithString("name", mcp.Description("milestone name")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(30)),
	)

	CreateMilestoneTool = mcp.NewTool(
		CreateMilestoneToolName,
		mcp.WithDescription("create milestone"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("title", mcp.Required(), mcp.Description("milestone title")),
		mcp.WithString("description", mcp.Description("milestone description")),
		mcp.WithString("due_on", mcp.Description("due date")),
	)

	EditMilestoneTool = mcp.NewTool(
		EditMilestoneToolName,
		mcp.WithDescription("edit milestone"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("milestone id")),
		mcp.WithString("title", mcp.Description("milestone title")),
		mcp.WithString("description", mcp.Description("milestone description")),
		mcp.WithString("due_on", mcp.Description("due date")),
		mcp.WithString("state", mcp.Description("milestone state, one of open, closed")),
	)

	DeleteMilestoneTool = mcp.NewTool(
		DeleteMilestoneToolName,
		mcp.WithDescription("delete milestone"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("milestone id")),
	)
)

func init() {
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetMilestoneTool,
		Handler: GetMilestoneFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    ListMilestonesTool,
		Handler: ListMilestonesFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreateMilestoneTool,
		Handler: CreateMilestoneFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    EditMilestoneTool,
		Handler: EditMilestoneFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    DeleteMilestoneTool,
		Handler: DeleteMilestoneFn,
	})
}

func GetMilestoneFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetMilestoneFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	id, err := params.GetIndex(req.GetArguments(), "id")
	if err != nil {
		return to.ErrorResult(err)
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	milestone, _, err := client.GetMilestone(owner, repo, id)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/milestone/%v err: %v", owner, repo, id, err))
	}

	return to.TextResult(slimMilestone(milestone))
}

func ListMilestonesFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListMilestonesFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	state := params.GetOptionalString(req.GetArguments(), "state", "all")
	name := params.GetOptionalString(req.GetArguments(), "name", "")
	page, pageSize := params.GetPagination(req.GetArguments(), 30)
	opt := gitea_sdk.ListMilestoneOption{
		State: gitea_sdk.StateType(state),
		Name:  name,
		ListOptions: gitea_sdk.ListOptions{
			Page:     page,
			PageSize: pageSize,
		},
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	milestones, _, err := client.ListRepoMilestones(owner, repo, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/milestones err: %v", owner, repo, err))
	}
	return to.TextResult(slimMilestones(milestones))
}

func CreateMilestoneFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreateMilestoneFn")
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

	opt := gitea_sdk.CreateMilestoneOption{
		Title: title,
	}

	description, ok := req.GetArguments()["description"].(string)
	if ok {
		opt.Description = description
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	milestone, _, err := client.CreateMilestone(owner, repo, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create %v/%v/milestone err: %v", owner, repo, err))
	}

	return to.TextResult(slimMilestone(milestone))
}

func EditMilestoneFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called EditMilestoneFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	id, err := params.GetIndex(req.GetArguments(), "id")
	if err != nil {
		return to.ErrorResult(err)
	}

	opt := gitea_sdk.EditMilestoneOption{}

	title, ok := req.GetArguments()["title"].(string)
	if ok {
		opt.Title = title
	}
	description, ok := req.GetArguments()["description"].(string)
	if ok {
		opt.Description = new(description)
	}
	state, ok := req.GetArguments()["state"].(string)
	if ok {
		opt.State = new(gitea_sdk.StateType(state))
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	milestone, _, err := client.EditMilestone(owner, repo, id, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("edit %v/%v/milestone/%v err: %v", owner, repo, id, err))
	}

	return to.TextResult(slimMilestone(milestone))
}

func DeleteMilestoneFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeleteMilestoneFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	id, err := params.GetIndex(req.GetArguments(), "id")
	if err != nil {
		return to.ErrorResult(err)
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	_, err = client.DeleteMilestone(owner, repo, id)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete %v/%v/milestone/%v err: %v", owner, repo, id, err))
	}

	return to.TextResult("Milestone deleted successfully")
}
