package pull

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
	MergePullRequestToolName              = "merge_pull_request"
	EditPullRequestToolName               = "edit_pull_request"
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
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(30)),
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
		mcp.WithArray("reviewers", mcp.Description("list of reviewer usernames"), mcp.Items(map[string]any{"type": "string"})),
		mcp.WithArray("team_reviewers", mcp.Description("list of team reviewer names"), mcp.Items(map[string]any{"type": "string"})),
	)

	DeletePullRequestReviewerTool = mcp.NewTool(
		DeletePullRequestReviewerToolName,
		mcp.WithDescription("remove reviewer requests from a pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithArray("reviewers", mcp.Description("list of reviewer usernames to remove"), mcp.Items(map[string]any{"type": "string"})),
		mcp.WithArray("team_reviewers", mcp.Description("list of team reviewer names to remove"), mcp.Items(map[string]any{"type": "string"})),
	)

	ListPullRequestReviewsTool = mcp.NewTool(
		ListPullRequestReviewsToolName,
		mcp.WithDescription("list all reviews for a pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithNumber("page", mcp.Description("page number"), mcp.DefaultNumber(1)),
		mcp.WithNumber("pageSize", mcp.Description("page size"), mcp.DefaultNumber(30)),
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
		mcp.WithArray("comments", mcp.Description("inline review comments (objects with path, body, old_line_num, new_line_num)"), mcp.Items(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path":         map[string]any{"type": "string", "description": "file path to comment on"},
				"body":         map[string]any{"type": "string", "description": "comment body"},
				"old_line_num": map[string]any{"type": "number", "description": "line number in the old file (for deletions/changes)"},
				"new_line_num": map[string]any{"type": "number", "description": "line number in the new file (for additions/changes)"},
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

	MergePullRequestTool = mcp.NewTool(
		MergePullRequestToolName,
		mcp.WithDescription("merge a pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithString("merge_style", mcp.Description("merge style: merge, rebase, rebase-merge, squash, fast-forward-only"), mcp.Enum("merge", "rebase", "rebase-merge", "squash", "fast-forward-only"), mcp.DefaultString("merge")),
		mcp.WithString("title", mcp.Description("custom merge commit title")),
		mcp.WithString("message", mcp.Description("custom merge commit message")),
		mcp.WithBoolean("delete_branch", mcp.Description("delete the branch after merge"), mcp.DefaultBool(false)),
	)

	EditPullRequestTool = mcp.NewTool(
		EditPullRequestToolName,
		mcp.WithDescription("edit a pull request"),
		mcp.WithString("owner", mcp.Required(), mcp.Description("repository owner")),
		mcp.WithString("repo", mcp.Required(), mcp.Description("repository name")),
		mcp.WithNumber("index", mcp.Required(), mcp.Description("pull request index")),
		mcp.WithString("title", mcp.Description("pull request title")),
		mcp.WithString("body", mcp.Description("pull request body content")),
		mcp.WithString("base", mcp.Description("pull request base branch")),
		mcp.WithString("assignee", mcp.Description("username to assign")),
		mcp.WithArray("assignees", mcp.Description("usernames to assign"), mcp.Items(map[string]any{"type": "string"})),
		mcp.WithNumber("milestone", mcp.Description("milestone number")),
		mcp.WithString("state", mcp.Description("pull request state"), mcp.Enum("open", "closed")),
		mcp.WithBoolean("allow_maintainer_edit", mcp.Description("allow maintainer to edit the pull request")),
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
	Tool.RegisterWrite(server.ServerTool{
		Tool:    MergePullRequestTool,
		Handler: MergePullRequestFn,
	})
	Tool.RegisterWrite(server.ServerTool{
		Tool:    EditPullRequestTool,
		Handler: EditPullRequestFn,
	})
}

func GetPullRequestByIndexFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetPullRequestByIndexFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}
	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	pr, _, err := client.GetPullRequest(owner, repo, index)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/pr/%v err: %v", owner, repo, index, err))
	}

	return to.TextResult(slimPullRequest(pr))
}

func GetPullRequestDiffFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetPullRequestDiffFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}
	binary, _ := args["binary"].(bool)

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}
	diffBytes, _, err := client.GetPullRequestDiff(owner, repo, index, gitea_sdk.PullRequestDiffOptions{
		Binary: binary,
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get %v/%v/pr/%v diff err: %v", owner, repo, index, err))
	}

	return to.TextResult(string(diffBytes))
}

func ListRepoPullRequestsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListRepoPullRequests")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	state, _ := args["state"].(string)
	sort := params.GetOptionalString(args, "sort", "recentupdate")
	milestone := params.GetOptionalInt(args, "milestone", 0)
	page, pageSize := params.GetPagination(args, 30)
	opt := gitea_sdk.ListPullRequestsOptions{
		State:     gitea_sdk.StateType(state),
		Sort:      sort,
		Milestone: milestone,
		ListOptions: gitea_sdk.ListOptions{
			Page:     page,
			PageSize: pageSize,
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

	return to.TextResult(slimPullRequests(pullRequests))
}

func CreatePullRequestFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreatePullRequestFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	title, err := params.GetString(args, "title")
	if err != nil {
		return to.ErrorResult(err)
	}
	body, err := params.GetString(args, "body")
	if err != nil {
		return to.ErrorResult(err)
	}
	head, err := params.GetString(args, "head")
	if err != nil {
		return to.ErrorResult(err)
	}
	base, err := params.GetString(args, "base")
	if err != nil {
		return to.ErrorResult(err)
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

	return to.TextResult(slimPullRequest(pr))
}

func CreatePullRequestReviewerFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreatePullRequestReviewerFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}

	reviewers := params.GetStringSlice(args, "reviewers")
	teamReviewers := params.GetStringSlice(args, "team_reviewers")

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	_, err = client.CreateReviewRequests(owner, repo, index, gitea_sdk.PullReviewRequestOptions{
		Reviewers:     reviewers,
		TeamReviewers: teamReviewers,
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create review requests for %v/%v/pr/%v err: %v", owner, repo, index, err))
	}

	// Return a success message instead of the Response object which contains non-serializable functions
	successMsg := map[string]any{
		"message":        "Successfully created review requests",
		"reviewers":      reviewers,
		"team_reviewers": teamReviewers,
		"pr_index":       index,
		"repository":     fmt.Sprintf("%s/%s", owner, repo),
	}

	return to.TextResult(successMsg)
}

func DeletePullRequestReviewerFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeletePullRequestReviewerFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}

	reviewers := params.GetStringSlice(args, "reviewers")
	teamReviewers := params.GetStringSlice(args, "team_reviewers")

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	_, err = client.DeleteReviewRequests(owner, repo, index, gitea_sdk.PullReviewRequestOptions{
		Reviewers:     reviewers,
		TeamReviewers: teamReviewers,
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete review requests for %v/%v/pr/%v err: %v", owner, repo, index, err))
	}

	successMsg := map[string]any{
		"message":        "Successfully deleted review requests",
		"reviewers":      reviewers,
		"team_reviewers": teamReviewers,
		"pr_index":       index,
		"repository":     fmt.Sprintf("%s/%s", owner, repo),
	}

	return to.TextResult(successMsg)
}

func ListPullRequestReviewsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListPullRequestReviewsFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}
	page, pageSize := params.GetPagination(args, 30)

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	reviews, _, err := client.ListPullReviews(owner, repo, index, gitea_sdk.ListPullReviewsOptions{
		ListOptions: gitea_sdk.ListOptions{
			Page:     page,
			PageSize: pageSize,
		},
	})
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list reviews for %v/%v/pr/%v err: %v", owner, repo, index, err))
	}

	return to.TextResult(slimReviews(reviews))
}

func GetPullRequestReviewFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called GetPullRequestReviewFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}
	reviewID, err := params.GetIndex(args, "review_id")
	if err != nil {
		return to.ErrorResult(err)
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	review, _, err := client.GetPullReview(owner, repo, index, reviewID)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get review %v for %v/%v/pr/%v err: %v", reviewID, owner, repo, index, err))
	}

	return to.TextResult(slimReview(review))
}

func ListPullRequestReviewCommentsFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called ListPullRequestReviewCommentsFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}
	reviewID, err := params.GetIndex(args, "review_id")
	if err != nil {
		return to.ErrorResult(err)
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	comments, _, err := client.ListPullReviewComments(owner, repo, index, reviewID)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("list review comments for review %v on %v/%v/pr/%v err: %v", reviewID, owner, repo, index, err))
	}

	return to.TextResult(slimReviewComments(comments))
}

func CreatePullRequestReviewFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called CreatePullRequestReviewFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}

	opt := gitea_sdk.CreatePullReviewOptions{}

	if state, ok := args["state"].(string); ok {
		opt.State = gitea_sdk.ReviewStateType(state)
	}
	if body, ok := args["body"].(string); ok {
		opt.Body = body
	}
	if commitID, ok := args["commit_id"].(string); ok {
		opt.CommitID = commitID
	}

	// Parse inline comments
	if commentsArg, exists := args["comments"]; exists {
		if commentsSlice, ok := commentsArg.([]any); ok {
			for _, comment := range commentsSlice {
				if commentMap, ok := comment.(map[string]any); ok {
					reviewComment := gitea_sdk.CreatePullReviewComment{}
					if path, ok := commentMap["path"].(string); ok {
						reviewComment.Path = path
					}
					if body, ok := commentMap["body"].(string); ok {
						reviewComment.Body = body
					}
					if oldLineNum, ok := params.ToInt64(commentMap["old_line_num"]); ok {
						reviewComment.OldLineNum = oldLineNum
					}
					if newLineNum, ok := params.ToInt64(commentMap["new_line_num"]); ok {
						reviewComment.NewLineNum = newLineNum
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

	review, _, err := client.CreatePullReview(owner, repo, index, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("create review for %v/%v/pr/%v err: %v", owner, repo, index, err))
	}

	return to.TextResult(slimReview(review))
}

func SubmitPullRequestReviewFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called SubmitPullRequestReviewFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}
	reviewID, err := params.GetIndex(args, "review_id")
	if err != nil {
		return to.ErrorResult(err)
	}
	state, err := params.GetString(args, "state")
	if err != nil {
		return to.ErrorResult(err)
	}

	opt := gitea_sdk.SubmitPullReviewOptions{
		State: gitea_sdk.ReviewStateType(state),
	}
	if body, ok := args["body"].(string); ok {
		opt.Body = body
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	review, _, err := client.SubmitPullReview(owner, repo, index, reviewID, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("submit review %v for %v/%v/pr/%v err: %v", reviewID, owner, repo, index, err))
	}

	return to.TextResult(slimReview(review))
}

func DeletePullRequestReviewFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DeletePullRequestReviewFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}
	reviewID, err := params.GetIndex(args, "review_id")
	if err != nil {
		return to.ErrorResult(err)
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	_, err = client.DeletePullReview(owner, repo, index, reviewID)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("delete review %v for %v/%v/pr/%v err: %v", reviewID, owner, repo, index, err))
	}

	successMsg := map[string]any{
		"message":    "Successfully deleted review",
		"review_id":  reviewID,
		"pr_index":   index,
		"repository": fmt.Sprintf("%s/%s", owner, repo),
	}

	return to.TextResult(successMsg)
}

func DismissPullRequestReviewFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called DismissPullRequestReviewFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}
	reviewID, err := params.GetIndex(args, "review_id")
	if err != nil {
		return to.ErrorResult(err)
	}

	opt := gitea_sdk.DismissPullReviewOptions{}
	if message, ok := args["message"].(string); ok {
		opt.Message = message
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	_, err = client.DismissPullReview(owner, repo, index, reviewID, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("dismiss review %v for %v/%v/pr/%v err: %v", reviewID, owner, repo, index, err))
	}

	successMsg := map[string]any{
		"message":    "Successfully dismissed review",
		"review_id":  reviewID,
		"pr_index":   index,
		"repository": fmt.Sprintf("%s/%s", owner, repo),
	}

	return to.TextResult(successMsg)
}

func MergePullRequestFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called MergePullRequestFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}

	mergeStyle := params.GetOptionalString(args, "merge_style", "merge")
	title, _ := args["title"].(string)
	message, _ := args["message"].(string)
	deleteBranch, _ := args["delete_branch"].(bool)

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	opt := gitea_sdk.MergePullRequestOption{
		Style:                  gitea_sdk.MergeStyle(mergeStyle),
		Title:                  title,
		Message:                message,
		DeleteBranchAfterMerge: deleteBranch,
	}

	merged, resp, err := client.MergePullRequest(owner, repo, index, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("merge %v/%v/pr/%v err: %v", owner, repo, index, err))
	}

	if !merged && resp != nil && resp.StatusCode >= 400 {
		return to.ErrorResult(fmt.Errorf("merge %v/%v/pr/%v failed: HTTP %d %s", owner, repo, index, resp.StatusCode, resp.Status))
	}

	if !merged {
		return to.ErrorResult(fmt.Errorf("merge %v/%v/pr/%v returned merged=false", owner, repo, index))
	}

	successMsg := map[string]any{
		"merged":         merged,
		"pr_index":       index,
		"repository":     fmt.Sprintf("%s/%s", owner, repo),
		"merge_style":    mergeStyle,
		"branch_deleted": deleteBranch,
	}

	return to.TextResult(successMsg)
}

func EditPullRequestFn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Debugf("Called EditPullRequestFn")
	args := req.GetArguments()
	owner, err := params.GetString(args, "owner")
	if err != nil {
		return to.ErrorResult(err)
	}
	repo, err := params.GetString(args, "repo")
	if err != nil {
		return to.ErrorResult(err)
	}
	index, err := params.GetIndex(args, "index")
	if err != nil {
		return to.ErrorResult(err)
	}

	opt := gitea_sdk.EditPullRequestOption{}

	if title, ok := args["title"].(string); ok {
		opt.Title = title
	}
	if body, ok := args["body"].(string); ok {
		opt.Body = new(body)
	}
	if base, ok := args["base"].(string); ok {
		opt.Base = base
	}
	if assignee, ok := args["assignee"].(string); ok {
		opt.Assignee = assignee
	}
	if assignees := params.GetStringSlice(args, "assignees"); assignees != nil {
		opt.Assignees = assignees
	}
	if val, exists := args["milestone"]; exists {
		if milestone, ok := params.ToInt64(val); ok {
			opt.Milestone = milestone
		}
	}
	if state, ok := args["state"].(string); ok {
		opt.State = new(gitea_sdk.StateType(state))
	}
	if allowMaintainerEdit, ok := args["allow_maintainer_edit"].(bool); ok {
		opt.AllowMaintainerEdit = new(allowMaintainerEdit)
	}

	client, err := gitea.ClientFromContext(ctx)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("get gitea client err: %v", err))
	}

	pr, _, err := client.EditPullRequest(owner, repo, index, opt)
	if err != nil {
		return to.ErrorResult(fmt.Errorf("edit %v/%v/pr/%v err: %v", owner, repo, index, err))
	}

	return to.TextResult(slimPullRequest(pr))
}
