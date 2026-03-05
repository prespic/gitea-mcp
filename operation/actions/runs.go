package actions

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"strconv"

	"gitea.com/gitea/gitea-mcp/pkg/gitea"
	"gitea.com/gitea/gitea-mcp/pkg/log"
	"gitea.com/gitea/gitea-mcp/pkg/params"
	"gitea.com/gitea/gitea-mcp/pkg/to"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	ListRepoActionWorkflowsToolName    = "list_repo_action_workflows"
	GetRepoActionWorkflowToolName      = "get_repo_action_workflow"
	DispatchRepoActionWorkflowToolName = "dispatch_repo_action_workflow"

	ListRepoActionRunsToolName  = "list_repo_action_runs"
	GetRepoActionRunToolName    = "get_repo_action_run"
	CancelRepoActionRunToolName = "cancel_repo_action_run"
	RerunRepoActionRunToolName  = "rerun_repo_action_run"

	ListRepoActionJobsToolName    = "list_repo_action_jobs"
	ListRepoActionRunJobsToolName = "list_repo_action_run_jobs"
)

var (
	ListRepoActionWorkflowsTool = mcp.NewTool(
		ListRepoActionWorkflowsToolName,
		mcp.WithDescription("List repository Actions workflows"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1), mcp.Min(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(30), mcp.Min(1)),
	)

	GetRepoActionWorkflowTool = mcp.NewTool(
		GetRepoActionWorkflowToolName,
		mcp.WithDescription("Get a repository Actions workflow by ID"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("workflow_id", mcp.Required(), mcp.Description("workflow ID or filename (e.g. 'my-workflow.yml')")),
	)

	DispatchRepoActionWorkflowTool = mcp.NewTool(
		DispatchRepoActionWorkflowToolName,
		mcp.WithDescription("Trigger (dispatch) a repository Actions workflow"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("workflow_id", mcp.Required(), mcp.Description("workflow ID or filename (e.g. 'my-workflow.yml')")),
		mcp.WithString("ref", mcp.Required(), mcp.Description("git ref (branch or tag)")),
		mcp.WithObject("inputs", mcp.Description("workflow inputs object")),
	)

	ListRepoActionRunsTool = mcp.NewTool(
		ListRepoActionRunsToolName,
		mcp.WithDescription("List repository Actions workflow runs"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1), mcp.Min(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(30), mcp.Min(1)),
		mcp.WithString("status", mcp.Description("optional status filter")),
	)

	GetRepoActionRunTool = mcp.NewTool(
		GetRepoActionRunToolName,
		mcp.WithDescription("Get a repository Actions run by ID"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("run_id", mcp.Required(), mcp.Description("run ID")),
	)

	CancelRepoActionRunTool = mcp.NewTool(
		CancelRepoActionRunToolName,
		mcp.WithDescription("Cancel a repository Actions run"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("run_id", mcp.Required(), mcp.Description("run ID")),
	)

	RerunRepoActionRunTool = mcp.NewTool(
		RerunRepoActionRunToolName,
		mcp.WithDescription("Rerun a repository Actions run"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("run_id", mcp.Required(), mcp.Description("run ID")),
	)

	ListRepoActionJobsTool = mcp.NewTool(
		ListRepoActionJobsToolName,
		mcp.WithDescription("List repository Actions jobs"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1), mcp.Min(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(30), mcp.Min(1)),
		mcp.WithString("status", mcp.Description("optional status filter")),
	)

	ListRepoActionRunJobsTool = mcp.NewTool(
		ListRepoActionRunJobsToolName,
		mcp.WithDescription("List Actions jobs for a specific run"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("run_id", mcp.Required(), mcp.Description("run ID")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1), mcp.Min(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(30), mcp.Min(1)),
	)
)

func init() {
	Tool.RegisterRead(server.ServerTool{Tool: ListRepoActionWorkflowsTool, Handler: ListRepoActionWorkflowsFn})
	Tool.RegisterRead(server.ServerTool{Tool: GetRepoActionWorkflowTool, Handler: GetRepoActionWorkflowFn})
	Tool.RegisterWrite(server.ServerTool{Tool: DispatchRepoActionWorkflowTool, Handler: DispatchRepoActionWorkflowFn})

	Tool.RegisterRead(server.ServerTool{Tool: ListRepoActionRunsTool, Handler: ListRepoActionRunsFn})
	Tool.RegisterRead(server.ServerTool{Tool: GetRepoActionRunTool, Handler: GetRepoActionRunFn})
	Tool.RegisterWrite(server.ServerTool{Tool: CancelRepoActionRunTool, Handler: CancelRepoActionRunFn})
	Tool.RegisterWrite(server.ServerTool{Tool: RerunRepoActionRunTool, Handler: RerunRepoActionRunFn})

	Tool.RegisterRead(server.ServerTool{Tool: ListRepoActionJobsTool, Handler: ListRepoActionJobsFn})
	Tool.RegisterRead(server.ServerTool{Tool: ListRepoActionRunJobsTool, Handler: ListRepoActionRunJobsFn})
}

func doJSONWithFallback(ctx context.Context, method string, paths []string, query url.Values, body, respOut any) error {
	var lastErr error
	for _, p := range paths {
		_, err := gitea.DoJSON(ctx, method, p, query, body, respOut)
		if err == nil {
			return nil
		}
		lastErr = err
		var httpErr *gitea.HTTPError
		if errors.As(err, &httpErr) && (httpErr.StatusCode == http.StatusNotFound || httpErr.StatusCode == http.StatusMethodNotAllowed) {
			continue
		}
		return err
	}
	return lastErr
}

func ListRepoActionWorkflowsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListRepoActionWorkflowsFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	page, pageSize := params.GetPagination(req.GetArguments(), 30)
	query := url.Values{}
	query.Set("page", strconv.Itoa(page))
	query.Set("limit", strconv.Itoa(pageSize))

	var result any
	err = doJSONWithFallback(ctx, "GET",
		[]string{
			fmt.Sprintf("repos/%s/%s/actions/workflows", url.PathEscape(owner), url.PathEscape(repo)),
		},
		query, nil, &result,
	)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list action workflows err: %v", err))
	}
	return to.TextResult(slimActionWorkflows(result))
}

func GetRepoActionWorkflowFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetRepoActionWorkflowFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	workflowID, err := params.GetString(req.GetArguments(), "workflow_id")
	if err != nil || workflowID == "" {
		return to.ErrorResult(errors.New("workflow_id is required"))
	}

	var result any
	err = doJSONWithFallback(ctx, "GET",
		[]string{
			fmt.Sprintf("repos/%s/%s/actions/workflows/%s", url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(workflowID)),
		},
		nil, nil, &result,
	)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get action workflow err: %v", err))
	}
	return to.TextResult(slimActionWorkflow(result))
}

func DispatchRepoActionWorkflowFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DispatchRepoActionWorkflowFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	workflowID, err := params.GetString(req.GetArguments(), "workflow_id")
	if err != nil || workflowID == "" {
		return to.ErrorResult(errors.New("workflow_id is required"))
	}
	ref, err := params.GetString(req.GetArguments(), "ref")
	if err != nil || ref == "" {
		return to.ErrorResult(errors.New("ref is required"))
	}

	var inputs map[string]any
	if raw, exists := req.GetArguments()["inputs"]; exists {
		if m, ok := raw.(map[string]any); ok {
			inputs = m
		} else if m, ok := raw.(map[string]any); ok {
			inputs = make(map[string]any, len(m))
			maps.Copy(inputs, m)
		}
	}

	body := map[string]any{
		"ref": ref,
	}
	if inputs != nil {
		body["inputs"] = inputs
	}

	err = doJSONWithFallback(ctx, "POST",
		[]string{
			fmt.Sprintf("repos/%s/%s/actions/workflows/%s/dispatches", url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(workflowID)),
			fmt.Sprintf("repos/%s/%s/actions/workflows/%s/dispatch", url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(workflowID)),
		},
		nil, body, nil,
	)
	if err != nil {
		var httpErr *gitea.HTTPError
		if errors.As(err, &httpErr) && (httpErr.StatusCode == http.StatusNotFound || httpErr.StatusCode == http.StatusMethodNotAllowed) {
			return to.ErrorResult(fmt.Errorf("workflow dispatch not supported on this Gitea version (endpoint returned %d). Check https://docs.gitea.com/api/1.24/ for available Actions endpoints", httpErr.StatusCode))
		}
		return to.ErrorResult(fmt.Errorf("dispatch action workflow err: %v", err))
	}
	return to.TextResult(map[string]any{"message": "workflow dispatched"})
}

func ListRepoActionRunsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListRepoActionRunsFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	page, pageSize := params.GetPagination(req.GetArguments(), 30)
	statusFilter, _ := req.GetArguments()["status"].(string)

	query := url.Values{}
	query.Set("page", strconv.Itoa(page))
	query.Set("limit", strconv.Itoa(pageSize))
	if statusFilter != "" {
		query.Set("status", statusFilter)
	}

	var result any
	err = doJSONWithFallback(ctx, "GET",
		[]string{
			fmt.Sprintf("repos/%s/%s/actions/runs", url.PathEscape(owner), url.PathEscape(repo)),
		},
		query, nil, &result,
	)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list action runs err: %v", err))
	}
	return to.TextResult(slimActionRuns(result))
}

func GetRepoActionRunFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetRepoActionRunFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	runID, err := params.GetIndex(req.GetArguments(), "run_id")
	if err != nil || runID <= 0 {
		return to.ErrorResult(errors.New("run_id is required"))
	}

	var result any
	err = doJSONWithFallback(ctx, "GET",
		[]string{
			fmt.Sprintf("repos/%s/%s/actions/runs/%d", url.PathEscape(owner), url.PathEscape(repo), runID),
		},
		nil, nil, &result,
	)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get action run err: %v", err))
	}
	return to.TextResult(slimActionRun(result))
}

func CancelRepoActionRunFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CancelRepoActionRunFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	runID, err := params.GetIndex(req.GetArguments(), "run_id")
	if err != nil || runID <= 0 {
		return to.ErrorResult(errors.New("run_id is required"))
	}

	err = doJSONWithFallback(ctx, "POST",
		[]string{
			fmt.Sprintf("repos/%s/%s/actions/runs/%d/cancel", url.PathEscape(owner), url.PathEscape(repo), runID),
		},
		nil, nil, nil,
	)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("cancel action run err: %v", err))
	}
	return to.TextResult(map[string]any{"message": "run cancellation requested"})
}

func RerunRepoActionRunFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called RerunRepoActionRunFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	runID, err := params.GetIndex(req.GetArguments(), "run_id")
	if err != nil || runID <= 0 {
		return to.ErrorResult(errors.New("run_id is required"))
	}

	err = doJSONWithFallback(ctx, "POST",
		[]string{
			fmt.Sprintf("repos/%s/%s/actions/runs/%d/rerun", url.PathEscape(owner), url.PathEscape(repo), runID),
			fmt.Sprintf("repos/%s/%s/actions/runs/%d/rerun-failed-jobs", url.PathEscape(owner), url.PathEscape(repo), runID),
		},
		nil, nil, nil,
	)
	if err != nil {
		var httpErr *gitea.HTTPError
		if errors.As(err, &httpErr) && (httpErr.StatusCode == http.StatusNotFound || httpErr.StatusCode == http.StatusMethodNotAllowed) {
			return to.ErrorResult(fmt.Errorf("workflow rerun not supported on this Gitea version (endpoint returned %d). Check https://docs.gitea.com/api/1.24/ for available Actions endpoints", httpErr.StatusCode))
		}
		return to.ErrorResult(fmt.Errorf("rerun action run err: %v", err))
	}
	return to.TextResult(map[string]any{"message": "run rerun requested"})
}

func ListRepoActionJobsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListRepoActionJobsFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	page, pageSize := params.GetPagination(req.GetArguments(), 30)
	statusFilter, _ := req.GetArguments()["status"].(string)

	query := url.Values{}
	query.Set("page", strconv.Itoa(page))
	query.Set("limit", strconv.Itoa(pageSize))
	if statusFilter != "" {
		query.Set("status", statusFilter)
	}

	var result any
	err = doJSONWithFallback(ctx, "GET",
		[]string{
			fmt.Sprintf("repos/%s/%s/actions/jobs", url.PathEscape(owner), url.PathEscape(repo)),
		},
		query, nil, &result,
	)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list action jobs err: %v", err))
	}
	return to.TextResult(slimActionJobs(result))
}

func ListRepoActionRunJobsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListRepoActionRunJobsFn")
	owner, err := params.GetString(req.GetArguments(), "owner")
	if err != nil || owner == "" {
		return to.ErrorResult(errors.New("owner is required"))
	}
	repo, err := params.GetString(req.GetArguments(), "repo")
	if err != nil || repo == "" {
		return to.ErrorResult(errors.New("repo is required"))
	}
	runID, err := params.GetIndex(req.GetArguments(), "run_id")
	if err != nil || runID <= 0 {
		return to.ErrorResult(errors.New("run_id is required"))
	}
	page, pageSize := params.GetPagination(req.GetArguments(), 30)

	query := url.Values{}
	query.Set("page", strconv.Itoa(page))
	query.Set("limit", strconv.Itoa(pageSize))

	var result any
	err = doJSONWithFallback(ctx, "GET",
		[]string{
			fmt.Sprintf("repos/%s/%s/actions/runs/%d/jobs", url.PathEscape(owner), url.PathEscape(repo), runID),
		},
		query, nil, &result,
	)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list action run jobs err: %v", err))
	}
	return to.TextResult(slimActionJobs(result))
}
