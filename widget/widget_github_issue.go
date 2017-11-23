package widget

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	go_github "github.com/google/go-github/github"
	"github.com/qmu/mcc/github"
	"github.com/qmu/mcc/widget/listable"
	"golang.org/x/text/width"
)

// GithubIssueWidget is a stack which shows a issue
// of the current branch referering its name including issue id
type GithubIssueWidget struct {
	options    *Option
	renderer   *listable.ListWrapper
	active     bool
	client     *github.Client
	timezone   string
	indent     int
	isReady    bool
	issueRegex string
}

// NewGithubIssueWidget constructs a New GithubIssueWidget
func NewGithubIssueWidget(opt *Option) (g *GithubIssueWidget, err error) {
	g = new(GithubIssueWidget)
	g.options = opt
	return
}

// Init renders the issue contents
func (g *GithubIssueWidget) Init() (err error) {
	g.indent = 9
	g.timezone = g.options.Timezone
	lopt := &listable.ListWrapperOption{
		Title:      g.options.GetTitle(),
		RealHeight: g.options.GetHeight(),
	}
	g.issueRegex = g.options.IssueRegex
	g.renderer = listable.NewListWrapper(lopt)
	return
}

// Activate is the implementation of Widget.Activate
func (g *GithubIssueWidget) Activate() {
	g.active = true
	if g.isReady {
		g.renderer.Activate()
	}
}

// Deactivate is the implementation of Widget.Deactivate
func (g *GithubIssueWidget) Deactivate() {
	g.active = false
	if g.isReady {
		g.renderer.Deactivate()
	}
}

// IsReady is the implementation of Widget.IsReady
func (g *GithubIssueWidget) IsReady() bool {
	return g.isReady
}

// GetGridBufferers is the implementation of Widget.Activate
func (g *GithubIssueWidget) GetGridBufferers() []ui.GridBufferer {
	return []ui.GridBufferer{g.renderer.GetWidget()}
}

func (g *GithubIssueWidget) buildIssueBody(issue *go_github.Issue, comments []*go_github.IssueComment) (body []string, err error) {
	desc := g.overflow(issue.GetBody())
	desc = g.putIndent(desc)

	// labels
	lbls := ""
	for _, lbl := range issue.Labels {
		lbls = lbls + "[" + lbl.GetName() + "] "
	}
	// milestone
	milestone := issue.Milestone.GetTitle()

	text := " [TITLE :](fg-blue) " + issue.GetTitle() + "\n"
	text += " [NO    :](fg-blue) " + "#" + fmt.Sprint(issue.GetNumber()) + "\n"
	text += " [BY    :](fg-blue) " + issue.User.GetLogin() + "\n"
	text += " [URL   :](fg-blue) " + issue.GetHTMLURL() + "\n"
	if lbls != "" {
		text += " [LABEL :](fg-blue) " + lbls + "\n"
	}
	if milestone != "" {
		text += " [MILE  :](fg-blue) " + milestone + "\n"
	}
	text += " [" + strings.Repeat("-", 300) + "](fg-blue) \n"
	text += " [DESC  :](fg-blue) " + desc
	text += " [" + strings.Repeat("-", 300) + "](fg-blue)"

	if len(comments) > 0 {
		commentText := "\n"
		for i, c := range comments {
			t := c.GetCreatedAt()
			loc, err := time.LoadLocation(g.timezone)
			if err != nil {
				return body, nil
			}
			if i > 0 {
				commentText += "[" + strings.Repeat(". ", 150) + "](fg-blue) \n\n"
			}
			commentText += "@" + c.User.GetLogin() + " [COMMENTED ON " + fmt.Sprint(t.In(loc)) + "](fg-blue)" + "\n\n"
			b := c.GetBody()
			commentText += g.overflow(b) + "\n"
			commentText += "\n"
		}

		commentText = g.putIndent(commentText)
		text += " [      :](fg-blue) " + commentText
	}

	body = strings.Split(text, "\n")

	return
}

func (g *GithubIssueWidget) buildPrBody(pr *go_github.PullRequest, comments []*go_github.IssueComment) (body []string, err error) {
	desc := g.overflow(pr.GetBody())
	desc = g.putIndent(desc)

	// milestone
	milestone := pr.Milestone.GetTitle()
	//
	text := " [TITLE :](fg-blue) " + pr.GetTitle() + "\n"
	text += " [NO    :](fg-blue) " + "#" + fmt.Sprint(pr.GetNumber()) + "\n"
	text += " [BY    :](fg-blue) " + pr.User.GetLogin() + "\n"
	text += " [URL   :](fg-blue) " + pr.GetHTMLURL() + "\n"
	if milestone != "" {
		text += " [MILE  :](fg-blue) " + milestone + "\n"
	}
	text += " [" + strings.Repeat("-", 300) + "](fg-blue) \n"
	text += " [DESC  :](fg-blue) " + desc
	text += " [" + strings.Repeat("-", 300) + "](fg-blue)"

	if len(comments) > 0 {
		commentText := "\n"
		for i, c := range comments {
			t := c.GetCreatedAt()
			loc, err := time.LoadLocation(g.timezone)
			if err != nil {
				return body, nil
			}
			if i > 0 {
				commentText += "[" + strings.Repeat(". ", 150) + "](fg-blue) \n\n"
			}
			commentText += "@" + c.User.GetLogin() + " [COMMENTED ON " + fmt.Sprint(t.In(loc)) + "](fg-blue)" + "\n\n"
			b := c.GetBody()
			commentText += g.overflow(b) + "\n"
			commentText += "\n"
		}

		commentText = g.putIndent(commentText)
		text += " [      :](fg-blue) " + commentText
	}

	body = strings.Split(text, "\n")

	return
}

func (g *GithubIssueWidget) overflow(text string) (result string) {
	lines := strings.Split(text, "\n")
	splitlen := g.renderer.GetWidth() - 2 - g.indent
	for _, line := range lines {
		cnt := 0
		for _, c := range line {
			ctype := width.LookupRune(c).Kind().String()
			t := string(ctype)
			if t == "Neutral" || t == "EastAsianNarrow" || t == "EastAsianHalfwidth" {
				cnt++
			} else {
				cnt += 2
			}
			if cnt+3 > splitlen {
				result += string(c) + "\n"
				cnt = 0
			} else {
				result += string(c)
			}
		}
		result += "\n"
	}
	return
}

func (g *GithubIssueWidget) putIndent(text string) (result string) {
	sp := strings.Split(text, "\n")
	indent := g.indent - 3
	for i, s := range sp {
		if i == 0 {
			result += s + "\n"
			continue
		}
		result += " [" + strings.Repeat(" ", indent) + ": ](fg-blue)" + s + "\n"
	}
	return
}

// SetOption is
func (g *GithubIssueWidget) SetOption(opt *AdditionalWidgetOption) {
	g.client = opt.GithubClient
	if g.client == nil {
		return
	}
	go func() {
		body := []string{}
		if g.options.Type == "github_issue" {
			err := g.client.SetIssueNoRegex(g.issueRegex)
			if err != nil {
				return
			}
			issue, comments, err := g.client.GetIssue(g.client.IssueID)
			if err != nil {
				return
			}
			body, err = g.buildIssueBody(issue, comments)
			if err != nil {
				return
			}
		} else if g.options.Type == "github_pr" {
			pr, comments, err := g.client.GetPR(g.client.IssueID)
			if err != nil {
				return
			}
			body, err = g.buildPrBody(pr, comments)
			if err != nil {
				return
			}
		} else {
			// TODO: should return error
			return
		}
		g.renderer.SetBody(body)
		g.renderer.ResetRender()
		g.isReady = true
		opt.done <- true
	}()
}
