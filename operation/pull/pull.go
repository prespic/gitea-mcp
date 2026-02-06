package pull

import (
	"context"
	"fmt"

	"gitea.com/gitea/gitea-mcp/pkg/gitea"
	"gitea.com/gitea/gitea-mcp/pkg/log"
	"gitea.com/gitea/gitea-mcp/pkg/to"
	"gitea.com/gitea/gitea-mcp/pkg/tool"

	gitea_sdk "code.gitea.io/sdk/gitea"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var Tool = tool.New()

const (
	GetPullRequestByIndexToolName         = "get_pull_request_by_index"
	GetPullRequestDiffToolName            = "get_pull_request_diff"
	ListRepoPullRequestsToolName          = "list_repo_pull_requests"
	CreatePullRequestToolName             = "create_pull_request"
	CreatePullRequestReviewerToolName     = "create_pull_request_reviewer"
	DeletePullRequestReviewerToolName     = "delete_pull_request_reviewer"
	ListPullRequestReviewsToolName        = "list_pull_request_reviews"
	GetPullRequestReviewToolName          = "get_pull_request_review"
	ListPullRequestReviewCommentsToolName = "list_pull_request_review_comments"
	CreatePullRequestReviewToolName       = "create_pull_request_review"
	SubmitPullRequestReviewToolName       = "submit_pull_request_review"
	DeletePullRequestReviewToolName       = "delete_pull_request_review"
	DismissPullRequestReviewToolName      = "dismiss_pull_request_review"
)

var (
	GetPullRequestByIndexTool = mcp.NewTool(
		GetPullRequestByIndexToolName,
		mcp.WithDescription("get pull request by index"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("repository pull request index")),
	)

	GetPullRequestDiffTool = mcp.NewTool(
		GetPullRequestDiffToolName,
		mcp.WithDescription("get pull request diff"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("repository pull request index")),
		mcp.WithBoolean("binary", mcp.Description("whether to include binary file changes")),
	)

	ListRepoPullRequestsTool = mcp.NewTool(
		ListRepoPullRequestsToolName,
		mcp.WithDescription("List repository pull requests"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("state", mcp.Description("state"), mcp.Enum("open", "closed", "all"), mcp.DefaultString("all")),
		mcp.WithString("sort", mcp.Description("sort"), mcp.Enum("oldest", "recentupdate", "leastupdate", "mostcomment", "leastcomment", "priority"), mcp.DefaultString("recentupdate")),
		mcp.WithNumber("milestone", mcp.Description("milestone")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(100)),
	)

	CreatePullRequestTool = mcp.NewTool(
		CreatePullRequestToolName,
		mcp.WithDescription("create pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithString("title", mcp.Required(), mcp.Description("pull request title")),
		mcp.WithString("body", mcp.Required(), mcp.Description("pull request body")),
		mcp.WithString("head", mcp.Required(), mcp.Description("pull request head")),
		mcp.WithString("base", mcp.Required(), mcp.Description("pull request base")),
	)

	CreatePullRequestReviewerTool = mcp.NewTool(
		CreatePullRequestReviewerToolName,
		mcp.WithDescription("create pull request reviewer"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithArray("reviewers", mcp.Description("list of reviewer usernames"), mcp.Items(map[string]interface{}{"type": "string"})),
		mcp.WithArray("team_reviewers", mcp.Description("list of team reviewer names"), mcp.Items(map[string]interface{}{"type": "string"})),
	)

	DeletePullRequestReviewerTool = mcp.NewTool(
		DeletePullRequestReviewerToolName,
		mcp.WithDescription("remove reviewer requests from a pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithArray("reviewers", mcp.Description("list of reviewer usernames to remove"), mcp.Items(map[string]interface{}{"type": "string"})),
		mcp.WithArray("team_reviewers", mcp.Description("list of team reviewer names to remove"), mcp.Items(map[string]interface{}{"type": "string"})),
	)

	ListPullRequestReviewsTool = mcp.NewTool(
		ListPullRequestReviewsToolName,
		mcp.WithDescription("list all reviews for a pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(100)),
	)

	GetPullRequestReviewTool = mcp.NewTool(
		GetPullRequestReviewToolName,
		mcp.WithDescription("get a specific review for a pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithNumber("review_id", mcp.Required(), mcp.Description("review ID")),
	)

	ListPullRequestReviewCommentsTool = mcp.NewTool(
		ListPullRequestReviewCommentsToolName,
		mcp.WithDescription("list all comments for a specific pull request review"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithNumber("review_id", mcp.Required(), mcp.Description("review ID")),
	)

	CreatePullRequestReviewTool = mcp.NewTool(
		CreatePullRequestReviewToolName,
		mcp.WithDescription("create a review for a pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithString("state", mcp.Description("review state"), mcp.Enum("APPROVED", "REQUEST_CHANGES", "COMMENT", "PENDING")),
		mcp.WithString("body", mcp.Description("review body/comment")),
		mcp.WithString("commit_id", mcp.Description("commit SHA to review")),
		mcp.WithArray("comments", mcp.Description("inline review comments (objects with path, body, old_line_num, new_line_num)"), mcp.Items(map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"path":         map[string]interface{}{"type": "string", "description": "file path to comment on"},
				"body":         map[string]interface{}{"type": "string", "description": "comment body"},
				"old_line_num": map[string]interface{}{"type": "number", "description": "line number in the old file (for deletions/changes)"},
				"new_line_num": map[string]interface{}{"type": "number", "description": "line number in the new file (for additions/changes)"},
			},
		})),
	)

	SubmitPullRequestReviewTool = mcp.NewTool(
		SubmitPullRequestReviewToolName,
		mcp.WithDescription("submit a pending pull request review"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithNumber("review_id", mcp.Required(), mcp.Description("review ID")),
		mcp.WithString("state", mcp.Required(), mcp.Description("final review state"), mcp.Enum("APPROVED", "REQUEST_CHANGES", "COMMENT")),
		mcp.WithString("body", mcp.Description("submission message")),
	)

	DeletePullRequestReviewTool = mcp.NewTool(
		DeletePullRequestReviewToolName,
		mcp.WithDescription("delete a pull request review"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithNumber("review_id", mcp.Required(), mcp.Description("review ID")),
	)

	DismissPullRequestReviewTool = mcp.NewTool(
		DismissPullRequestReviewToolName,
		mcp.WithDescription("dismiss a pull request review"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithNumber("review_id", mcp.Required(), mcp.Description("review ID")),
		mcp.WithString("message", mcp.Description("dismissal reason")),
	)
)

func init() {
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetPullRequestByIndexTool,
		Handler: GetPullRequestByIndexFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetPullRequestDiffTool,
		Handler: GetPullRequestDiffFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    ListRepoPullRequestsTool,
		Handler: ListRepoPullRequestsFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    ListPullRequestReviewsTool,
		Handler: ListPullRequestReviewsFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    GetPullRequestReviewTool,
		Handler: GetPullRequestReviewFn,
	})
	Tool.RegisterRead(server.ServerTool{
		Tool:    ListPullRequestReviewCommentsTool,
		Handler: ListPullRequestReviewCommentsFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreatePullRequestTool,
		Handler: CreatePullRequestFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreatePullRequestReviewerTool,
		Handler: CreatePullRequestReviewerFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    DeletePullRequestReviewerTool,
		Handler: DeletePullRequestReviewerFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    CreatePullRequestReviewTool,
		Handler: CreatePullRequestReviewFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    SubmitPullRequestReviewTool,
		Handler: SubmitPullRequestReviewFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    DeletePullRequestReviewTool,
		Handler: DeletePullRequestReviewFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    DismissPullRequestReviewTool,
		Handler: DismissPullRequestReviewFn,
	})
}

func GetPullRequestByIndexFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetPullRequestByIndexFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	pr, _, err := client.GetPullRequest(owner, repo, int64(index))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/pr/%v err: %v", owner, repo, int64(index), err))
	}

	return to.TextResult(pr)
}

func GetPullRequestDiffFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetPullRequestDiffFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}
	binary, _ := req.GetArguments()["binary"].(bool)

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	diffBytes, _, err := client.GetPullRequestDiff(owner, repo, int64(index), gitea_sdk.PullRequestDiffOptions{
		Binary: binary,
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/pr/%v diff err: %v", owner, repo, int64(index), err))
	}

	result := map[string]interface{}{
		"diff":   string(diffBytes),
		"binary": binary,
		"index":  int64(index),
		"repo":   repo,
		"owner":  owner,
	}
	return to.TextResult(result)
}

func ListRepoPullRequestsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListRepoPullRequests")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	state, _ := req.GetArguments()["state"].(string)
	sort, ok := req.GetArguments()["sort"].(string)
	if !ok {
		sort = "recentupdate"
	}
	milestone, _ := req.GetArguments()["milestone"].(float64)
	page, ok := req.GetArguments()["page"].(float64)
	if !ok {
		page = 1
	}
	pageSize, ok := req.GetArguments()["pageSize"].(float64)
	if !ok {
		pageSize = 100
	}
	opt := gitea_sdk.ListPullRequestsOptions{
		State:     gitea_sdk.StateType(state),
		Sort:      sort,
		Milestone: int64(milestone),
		ListOptions: gitea_sdk.ListOptions{
			Page:     int(page),
			PageSize: int(pageSize),
		},
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	pullRequests, _, err := client.ListRepoPullRequests(owner, repo, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list %v/%v/pull_requests err: %v", owner, repo, err))
	}

	return to.TextResult(pullRequests)
}

func CreatePullRequestFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreatePullRequestFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	title, ok := req.GetArguments()["title"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("title is required"))
	}
	body, ok := req.GetArguments()["body"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("body is required"))
	}
	head, ok := req.GetArguments()["head"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("head is required"))
	}
	base, ok := req.GetArguments()["base"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("base is required"))
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	pr, _, err := client.CreatePullRequest(owner, repo, gitea_sdk.CreatePullRequestOption{
		Title: title,
		Body:  body,
		Head:  head,
		Base:  base,
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create %v/%v/pull_request err: %v", owner, repo, err))
	}

	return to.TextResult(pr)
}

func CreatePullRequestReviewerFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreatePullRequestReviewerFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}

	var reviewers []string
	if reviewersArg, exists := req.GetArguments()["reviewers"]; exists {
		if reviewersSlice, ok := reviewersArg.([]interface{}); ok {
			for _, reviewer := range reviewersSlice {
				if reviewerStr, ok := reviewer.(string); ok {
					reviewers = append(reviewers, reviewerStr)
				}
			}
		}
	}

	var teamReviewers []string
	if teamReviewersArg, exists := req.GetArguments()["team_reviewers"]; exists {
		if teamReviewersSlice, ok := teamReviewersArg.([]interface{}); ok {
			for _, teamReviewer := range teamReviewersSlice {
				if teamReviewerStr, ok := teamReviewer.(string); ok {
					teamReviewers = append(teamReviewers, teamReviewerStr)
				}
			}
		}
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	_, err = client.CreateReviewRequests(owner, repo, int64(index), gitea_sdk.PullReviewRequestOptions{
		Reviewers:     reviewers,
		TeamReviewers: teamReviewers,
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create review requests for %v/%v/pr/%v err: %v", owner, repo, int64(index), err))
	}

	// Return a success message instead of the Response object which contains non-serializable functions
	successMsg := map[string]interface{}{
		"message":        "Successfully created review requests",
		"reviewers":      reviewers,
		"team_reviewers": teamReviewers,
		"pr_index":       int64(index),
		"repository":     fmt.Sprintf("%s/%s", owner, repo),
	}

	return to.TextResult(successMsg)
}

func DeletePullRequestReviewerFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeletePullRequestReviewerFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}

	var reviewers []string
	if reviewersArg, exists := req.GetArguments()["reviewers"]; exists {
		if reviewersSlice, ok := reviewersArg.([]interface{}); ok {
			for _, reviewer := range reviewersSlice {
				if reviewerStr, ok := reviewer.(string); ok {
					reviewers = append(reviewers, reviewerStr)
				}
			}
		}
	}

	var teamReviewers []string
	if teamReviewersArg, exists := req.GetArguments()["team_reviewers"]; exists {
		if teamReviewersSlice, ok := teamReviewersArg.([]interface{}); ok {
			for _, teamReviewer := range teamReviewersSlice {
				if teamReviewerStr, ok := teamReviewer.(string); ok {
					teamReviewers = append(teamReviewers, teamReviewerStr)
				}
			}
		}
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	_, err = client.DeleteReviewRequests(owner, repo, int64(index), gitea_sdk.PullReviewRequestOptions{
		Reviewers:     reviewers,
		TeamReviewers: teamReviewers,
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete review requests for %v/%v/pr/%v err: %v", owner, repo, int64(index), err))
	}

	successMsg := map[string]interface{}{
		"message":        "Successfully deleted review requests",
		"reviewers":      reviewers,
		"team_reviewers": teamReviewers,
		"pr_index":       int64(index),
		"repository":     fmt.Sprintf("%s/%s", owner, repo),
	}

	return to.TextResult(successMsg)
}

func ListPullRequestReviewsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListPullRequestReviewsFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}
	page, ok := req.GetArguments()["page"].(float64)
	if !ok {
		page = 1
	}
	pageSize, ok := req.GetArguments()["pageSize"].(float64)
	if !ok {
		pageSize = 100
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	reviews, _, err := client.ListPullReviews(owner, repo, int64(index), gitea_sdk.ListPullReviewsOptions{
		ListOptions: gitea_sdk.ListOptions{
			Page:     int(page),
			PageSize: int(pageSize),
		},
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list reviews for %v/%v/pr/%v err: %v", owner, repo, int64(index), err))
	}

	return to.TextResult(reviews)
}

func GetPullRequestReviewFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetPullRequestReviewFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}
	reviewID, ok := req.GetArguments()["review_id"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("review_id is required"))
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	review, _, err := client.GetPullReview(owner, repo, int64(index), int64(reviewID))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get review %v for %v/%v/pr/%v err: %v", int64(reviewID), owner, repo, int64(index), err))
	}

	return to.TextResult(review)
}

func ListPullRequestReviewCommentsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListPullRequestReviewCommentsFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}
	reviewID, ok := req.GetArguments()["review_id"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("review_id is required"))
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	comments, _, err := client.ListPullReviewComments(owner, repo, int64(index), int64(reviewID))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list review comments for review %v on %v/%v/pr/%v err: %v", int64(reviewID), owner, repo, int64(index), err))
	}

	return to.TextResult(comments)
}

func CreatePullRequestReviewFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreatePullRequestReviewFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}

	opt := gitea_sdk.CreatePullReviewOptions{}

	if state, ok := req.GetArguments()["state"].(string); ok {
		opt.State = gitea_sdk.ReviewStateType(state)
	}
	if body, ok := req.GetArguments()["body"].(string); ok {
		opt.Body = body
	}
	if commitID, ok := req.GetArguments()["commit_id"].(string); ok {
		opt.CommitID = commitID
	}

	// Parse inline comments
	if commentsArg, exists := req.GetArguments()["comments"]; exists {
		if commentsSlice, ok := commentsArg.([]interface{}); ok {
			for _, comment := range commentsSlice {
				if commentMap, ok := comment.(map[string]interface{}); ok {
					reviewComment := gitea_sdk.CreatePullReviewComment{}
					if path, ok := commentMap["path"].(string); ok {
						reviewComment.Path = path
					}
					if body, ok := commentMap["body"].(string); ok {
						reviewComment.Body = body
					}
					if oldLineNum, ok := commentMap["old_line_num"].(float64); ok {
						reviewComment.OldLineNum = int64(oldLineNum)
					}
					if newLineNum, ok := commentMap["new_line_num"].(float64); ok {
						reviewComment.NewLineNum = int64(newLineNum)
					}
					opt.Comments = append(opt.Comments, reviewComment)
				}
			}
		}
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	review, _, err := client.CreatePullReview(owner, repo, int64(index), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create review for %v/%v/pr/%v err: %v", owner, repo, int64(index), err))
	}

	return to.TextResult(review)
}

func SubmitPullRequestReviewFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called SubmitPullRequestReviewFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}
	reviewID, ok := req.GetArguments()["review_id"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("review_id is required"))
	}
	state, ok := req.GetArguments()["state"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("state is required"))
	}

	opt := gitea_sdk.SubmitPullReviewOptions{
		State: gitea_sdk.ReviewStateType(state),
	}
	if body, ok := req.GetArguments()["body"].(string); ok {
		opt.Body = body
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	review, _, err := client.SubmitPullReview(owner, repo, int64(index), int64(reviewID), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("submit review %v for %v/%v/pr/%v err: %v", int64(reviewID), owner, repo, int64(index), err))
	}

	return to.TextResult(review)
}

func DeletePullRequestReviewFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeletePullRequestReviewFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}
	reviewID, ok := req.GetArguments()["review_id"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("review_id is required"))
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	_, err = client.DeletePullReview(owner, repo, int64(index), int64(reviewID))
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete review %v for %v/%v/pr/%v err: %v", int64(reviewID), owner, repo, int64(index), err))
	}

	successMsg := map[string]interface{}{
		"message":    "Successfully deleted review",
		"review_id":  int64(reviewID),
		"pr_index":   int64(index),
		"repository": fmt.Sprintf("%s/%s", owner, repo),
	}

	return to.TextResult(successMsg)
}

func DismissPullRequestReviewFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DismissPullRequestReviewFn")
	owner, ok := req.GetArguments()["owner"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("owner is required"))
	}
	repo, ok := req.GetArguments()["repo"].(string)
	if !ok {
		return to.ErrorResult(fmt.Errorf("repo is required"))
	}
	index, ok := req.GetArguments()["index"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("index is required"))
	}
	reviewID, ok := req.GetArguments()["review_id"].(float64)
	if !ok {
		return to.ErrorResult(fmt.Errorf("review_id is required"))
	}

	opt := gitea_sdk.DismissPullReviewOptions{}
	if message, ok := req.GetArguments()["message"].(string); ok {
		opt.Message = message
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	_, err = client.DismissPullReview(owner, repo, int64(index), int64(reviewID), opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("dismiss review %v for %v/%v/pr/%v err: %v", int64(reviewID), owner, repo, int64(index), err))
	}

	successMsg := map[string]interface{}{
		"message":    "Successfully dismissed review",
		"review_id":  int64(reviewID),
		"pr_index":   int64(index),
		"repository": fmt.Sprintf("%s/%s", owner, repo),
	}

	return to.TextResult(successMsg)
}
