package dashboard

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/github"
	"golang.org/x/text/width"
)

// GithubIssueWidget is a widget which shows a issue
// of the current branch referering its name including issue id
type GithubIssueWidget struct {
	renderer   *ListWrapper
	active     bool
	client     *github.Client
	timezone   string
	indent     int
	isReady    bool
	disabled   bool
	issueRegex string
}

// NewGithubIssueWidget constructs a New GithubIssueWidget
func NewGithubIssueWidget(wi Widget, timezone string) (g *GithubIssueWidget, err error) {
	g = new(GithubIssueWidget)
	g.indent = 9
	g.timezone = timezone
	opt := &ListWrapperOption{
		Title:      wi.Title,
		RealHeight: wi.RealHeight,
	}
	g.issueRegex = wi.IssueRegex
	g.renderer = NewListWrapper(opt)

	return
}

// Activate is the implementation of Widget.Activate
func (g *GithubIssueWidget) Activate() {
	g.active = true
	if g.isReady {
		g.renderer.Render()
	}
}

// Deactivate is the implementation of Widget.Deactivate
func (g *GithubIssueWidget) Deactivate() {
	g.active = false
	if g.isReady {
		g.renderer.ResetRender()
	}
}

// IsDisabled is the implementation of Widget.IsDisabled
func (g *GithubIssueWidget) IsDisabled() bool {
	return g.disabled
}

// IsReady is the implementation of Widget.IsReady
func (g *GithubIssueWidget) IsReady() bool {
	return g.isReady
}

// GetHighlightenPos is the implementation of Widget.GetHighlightenPos
func (g *GithubIssueWidget) GetHighlightenPos() int {
	return g.renderer.GetCursor()
}

// Disable sets a GithubIssueWidget instance as disabled
func (g *GithubIssueWidget) Disable() {
	g.disabled = true
	g.renderer.SetBody([]string{"Could not load issue number from branch name..."})
	ui.Render(ui.Body)
}

// Render renders the issue contents
func (g *GithubIssueWidget) Render(client *github.Client) (err error) {
	g.client = client

	body, err := g.buildBody()
	if err != nil {
		return
	}

	g.renderer.SetBody(body)
	g.isReady = true

	if g.active {
		g.renderer.Render()
	} else {
		g.renderer.ResetRender()
	}
	return
}

// GetWidget is the implementation of Widget.Activate
func (g *GithubIssueWidget) GetWidget() *ui.List {
	return g.renderer.GetWidget()
}

func (g *GithubIssueWidget) buildBody() (body []string, err error) {
	issue, comments, err := g.client.GetIssue(g.issueRegex)
	if err != nil {
		return
	}
	desc := g.overflow(issue.GetBody())
	desc = g.putIndent(desc)
	commentText := ""
	for i, c := range comments {
		t := c.GetCreatedAt()
		loc, err := time.LoadLocation(g.timezone)
		if err != nil {
			return body, nil
		}
		if i > 0 {
			commentText += "[" + strings.Repeat(". ", 150) + "](fg-blue) \n\n"
		}
		commentText += "[COMMENTED BY ](fg-blue)" + c.User.GetLogin() + " [ON " + fmt.Sprint(t.In(loc)) + "](fg-blue)" + "\n"
		b := c.GetBody()
		commentText += g.overflow(b) + "\n"
		commentText += "\n"
	}

	// labels
	lbls := ""
	for _, lbl := range issue.Labels {
		lbls = lbls + "[" + lbl.GetName() + "] "
	}
	// milestone
	milestone := issue.Milestone.GetTitle()

	commentText = g.putIndent(commentText)
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
	text += " [" + strings.Repeat("-", 300) + "](fg-blue) \n"
	text += " [      :](fg-blue) " + commentText
	text += " [" + strings.Repeat("-", 300) + "](fg-blue)"

	body = strings.Split(text, "\n")

	return
}

func (g *GithubIssueWidget) overflow(text string) (result string) {
	lines := strings.Split(text, "\n")
	w := g.GetWidget()
	splitlen := w.Width - 2 - g.indent
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
