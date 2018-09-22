package widget

import (
	"path/filepath"
	"testing"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/github"
)

func TestNewGithubIssueWidget(t *testing.T) {
	d, err := filepath.Abs("../") // current dir is ./widget
	if err != nil {
		return
	}
	wi, err := NewGithubIssueWidget(&Option{
		ExecPath:   d,
		Height:     100,
		Width:      100,
		Title:      "test",
		Type:       "github_issue",
		IssueRegex: "i([0-9]*).*",
	})
	if err != nil {
		return
	}
	ui.Init()
	wi.Init()
	c, err := github.NewClient(d, "github.com")
	if err != nil {
		t.Fatalf("error %v", err)
	}
	if err = c.Init(); err != nil {
		t.Fatalf("error %v", err)
	}

	done := make(chan bool)
	wi.SetOption(&AdditionalWidgetOption{
		done:         done,
		GithubClient: c,
	})
	<-done

	issue, comments, err := wi.client.GetIssue(1)
	if err != nil {
		t.Fatalf("error %v", err)
	}
	body, err := wi.buildIssueBody(issue, comments)
	if err != nil {
		t.Fatalf("error %v", err)
	}
	wi.renderer.SetBody(body)
	result := wi.renderer.Render()

	if len(result) == 0 {
		t.Fatalf("error %v", result)
	}

	ui.Close()
}
