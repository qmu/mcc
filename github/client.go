package github

import (
	"context"
	"errors"
	"regexp"
	"strconv"

	go_github "github.com/google/go-github/github"
	"github.com/qmu/mcc/utils"
	"gopkg.in/src-d/go-git.v4"
	// "github.com/k0kubun/pp"
)

// Client is a GitHub OAuth API wrapper
type Client struct {
	client     *go_github.Client
	host       string
	dotGitPath string
	repoOwner  string
	repoName   string
	branch     string
	IssueID    int
	auth       *AuthService
}

// NewClient constructs a new Client
func NewClient(execPath string, host string) (g *Client, err error) {
	g = new(Client)
	g.dotGitPath, err = utils.GetDotGitPath(execPath)
	g.host = host
	if err != nil {
		return
	}

	r, err := git.PlainOpen(g.dotGitPath)
	if err != nil {
		return
	}
	// get branch info
	ref, err := r.Head()
	if err != nil {
		return
	}
	g.branch = ref.Name().Short()
	// get github info
	remotes, err := r.Remotes()
	if err != nil {
		return
	}
	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			u := remote.Config().URLs[0]
			rep1 := regexp.MustCompile(`.*:(.*)/.*`)
			g.repoOwner = rep1.ReplaceAllString(u, "$1")
			rep2 := regexp.MustCompile(`.*/(.*)\.git`)
			g.repoName = rep2.ReplaceAllString(u, "$1")
			break
		}
	}
	return
}

// Init initialize Client
func (g *Client) Init() (err error) {
	// check if it's public repository
	ctx := context.Background()
	g.client = go_github.NewClient(nil)
	_, _, err = g.client.Repositories.Get(ctx, g.repoOwner, g.repoName)
	// if it's private, authenticate first
	if err != nil {
		g.auth, err = NewAuthService(g.host)
		if err != nil {
			return
		}
		g.client, err = g.auth.InitClient()
	}
	return
}

// GetIssue requests an issue and comments by refering current branch name which includes IssueID
func (g *Client) GetIssue(issueID int) (issue *go_github.Issue, comments []*go_github.IssueComment, err error) {
	// get a issue
	ctx := context.Background()
	issue, _, err = g.client.Issues.Get(ctx, g.repoOwner, g.repoName, issueID)
	if err != nil {
		return
	}
	opt := new(go_github.IssueListCommentsOptions)
	comments, _, err = g.client.Issues.ListComments(ctx, g.repoOwner, g.repoName, issueID, opt)

	return
}

// GetPR is
func (g *Client) GetPR(issueID int) (pr *go_github.PullRequest, comments []*go_github.IssueComment, err error) {
	ctx := context.Background()
	br, _, err := g.client.Repositories.GetBranch(ctx, g.repoOwner, g.repoName, g.branch)
	if err != nil {
		return
	}

	prs, _, err := g.client.PullRequests.List(ctx, g.repoOwner, g.repoName, &go_github.PullRequestListOptions{
		Head: br.Commit.GetSHA(),
	})
	if err != nil {
		return
	}
	if len(prs) == 0 {
		return nil, nil, errors.New("Fetched not only one pull request")
	}
	pr = prs[0]
	opt := new(go_github.IssueListCommentsOptions)
	ctx = context.Background()
	comments, _, err = g.client.Issues.ListComments(ctx, g.repoOwner, g.repoName, pr.GetNumber(), opt)
	return
}

// SetIssueNoRegex sets
func (g *Client) SetIssueNoRegex(issueNoRegex string) (err error) {
	rep0 := regexp.MustCompile(issueNoRegex)
	issueID := rep0.ReplaceAllString(g.branch, "$1")
	g.IssueID, err = strconv.Atoi(issueID)
	if err != nil {
		return
	}
	return
}
