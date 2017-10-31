package dashboard

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"unicode/utf8"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/utils"
	"gopkg.in/src-d/go-git.v4"
	// "github.com/k0kubun/pp"
)

// GitStatusWidget is a command launcher
type GitStatusWidget struct {
	renderer    *ListWrapper
	isReady     bool
	disabled    bool
	statusItems StatusItems
	envs        []map[string]string
	execPath    string
}

// StatusItem is a struct which stores each file status of git status
type StatusItem struct {
	Staged   bool
	StatusNo int
	Stage    string
	Status   string
	Path     string
}

// StatusItems is a collection of StatusItem, and implements sorting
type StatusItems []StatusItem

// Len is interface method of sort
func (s StatusItems) Len() int {
	return len(s)
}

// Swap is interface method of sort
func (s StatusItems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// ByStage is a struct for sorting StatusItems
type ByStage struct {
	StatusItems
}

// Less is interface method of sort
func (s ByStage) Less(i, j int) bool {
	return len(s.StatusItems[i].Stage) < len(s.StatusItems[j].Stage)
}

// ByPath is a struct for sorting StatusItems
type ByPath struct {
	StatusItems
}

// Less is interface method of sort
func (b ByPath) Less(i, j int) bool {
	return b.StatusItems[i].Path < b.StatusItems[j].Path
}

// NewGitStatusWidget constructs a New GitStatusWidget
func NewGitStatusWidget(wi Widget, execPath string, envs []map[string]string) (g *GitStatusWidget, err error) {
	g = new(GitStatusWidget)
	g.envs = envs
	g.execPath = execPath

	header := g.buildHeader()

	opt := &ListWrapperOption{
		Title:         wi.Title,
		RealHeight:    wi.RealHeight,
		Header:        header,
		LineHighLight: true,
	}
	g.renderer = NewListWrapper(opt)
	g.isReady = true

	return
}

func (g *GitStatusWidget) buildHeader() (header []string) {
	n1, n2, n3 := g.getLongest()
	c1 := g.fillSpaces("STAGE ", n1)
	c2 := g.fillSpaces("STATUS ", n2)
	c3 := g.fillSpaces("PATH ", n3)
	header = []string{
		" [" + c1 + " | " + c2 + " | " + c3 + "](fg-blue)\n",
		" [" + strings.Repeat("-", 500) + "](fg-blue)\n"}
	return
}

func (g *GitStatusWidget) buildBody(execPath string) (result []string, err error) {
	status, err := g.getStatus(execPath)
	if err != nil {
		return
	}
	if status.IsClean() {
		return nil, err
	}

	for path, s := range status {
		if s.Staging == git.Unmodified && s.Worktree == git.Unmodified {
			continue
		}
		if s.Staging == git.Renamed {
			path = fmt.Sprintf("%s -> %s", path, s.Extra)
		}

		var gs git.StatusCode
		var stage string
		if s.Staging != git.Unmodified {
			gs = s.Staging
			stage = "STAGED"
		} else {
			gs = s.Worktree
			stage = "UNSTAGED"
		}
		var statusStr string
		var statusNo int
		if gs == git.Untracked {
			statusStr = "New File"
			statusNo = 1
		} else if gs == git.Modified {
			statusStr = "Modified"
			statusNo = 2
		} else if gs == git.Added {
			statusStr = "Added"
			statusNo = 3
		} else if gs == git.Deleted {
			statusStr = "Deleted"
			statusNo = 4
		} else if gs == git.Renamed {
			statusStr = "Renamed"
			statusNo = 5
		} else if gs == git.Copied {
			statusStr = "Copied"
			statusNo = 6
		} else if gs == git.UpdatedButUnmerged {
			statusStr = "UpdatedButUnmerged"
			statusNo = 7
		}
		g.statusItems = append(g.statusItems, StatusItem{
			Staged:   s.Staging != git.Unmodified,
			Stage:    stage,
			Status:   statusStr,
			StatusNo: statusNo,
			Path:     path,
		})
	}
	// sort
	sort.Sort(ByPath{g.statusItems})
	sort.Sort(ByStage{g.statusItems})

	// build body
	n1, n2, _ := g.getLongest()
	for _, statusItem := range g.statusItems {
		s1 := g.fillSpaces(statusItem.Stage, n1)
		s2 := g.fillSpaces(statusItem.Status, n2)
		s3 := statusItem.Path + strings.Repeat(" ", 200)
		var st string
		if statusItem.Staged {
			st = "[" + s1 + "](fg-green)"
		} else {
			st = "[" + s1 + "](fg-red)"
		}
		result = append(result, " "+st+" [|](fg-blue) "+s2+" [|](fg-blue) "+s3)
	}

	return
}

func (g *GitStatusWidget) getStatus(execPath string) (status git.Status, err error) {
	// Load worktree status
	dotGitPath, err := utils.GetDotGitPath(execPath)
	r, err := git.PlainOpen(dotGitPath)
	if err != nil {
		return
	}
	w, err := r.Worktree()
	if err != nil {
		return
	}
	status, err = w.Status()
	if err != nil {
		return
	}
	return
}

func (g *GitStatusWidget) fillSpaces(s string, longest int) string {
	var l = longest - utf8.RuneCountInString(s)
	for i := 0; i < l; i++ {
		s += " "
	}
	return s
}

func (g *GitStatusWidget) getLongest() (n1 int, n2 int, n3 int) {
	n1 = utf8.RuneCountInString("STAGE") + 1
	n2 = utf8.RuneCountInString("STATUS") + 1
	n3 = utf8.RuneCountInString("PATH") + 1
	for _, statusItem := range g.statusItems {
		c := utf8.RuneCountInString(statusItem.Stage)
		nm := utf8.RuneCountInString(statusItem.Status)
		d := utf8.RuneCountInString(statusItem.Path)
		if n1 < c {
			n1 = c
		}
		if n2 < nm {
			n2 = nm
		}
		if n3 < d {
			n3 = d
		}
	}
	return n1, n2, n3
}

// Activate is the implementation of Widget.Activate
func (g *GitStatusWidget) Activate() {
	g.setKeyBindings()
	g.renderer.Render()
}

// Deactivate is the implementation of Widget.Activate
func (g *GitStatusWidget) Deactivate() {
	g.renderer.ResetRender()
}

// IsDisabled is the implementation of Widget.IsDisabled
func (g *GitStatusWidget) IsDisabled() bool {
	return g.disabled
}

// IsReady is the implementation of Widget.IsReady
func (g *GitStatusWidget) IsReady() bool {
	return g.isReady
}

// GetHighlightenPos is the implementation of Widget.GetHighlightenPos
func (g *GitStatusWidget) GetHighlightenPos() int {
	return g.renderer.GetCursor()
}

// GetWidget is the implementation of widget.Activate
func (g *GitStatusWidget) GetWidget() []ui.GridBufferer {
	return []ui.GridBufferer{g.renderer.GetWidget()}
}

func (g *GitStatusWidget) setKeyBindings() error {
	// exec command by Enter
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		ui.StopLoop()
		ui.Close()

		cursor := g.renderer.GetCursor()
		hasEditor := os.Getenv("EDITOR") != ""
		editorCmd := ""
		if hasEditor {
			editorCmd = os.Getenv("EDITOR")
		}
		for _, env := range g.envs {
			if env["name"] == "EDITOR" {
				hasEditor = true
				editorCmd = env["value"]
			}
		}

		if !hasEditor {
			log.Println("Set an enviromental variable \"EDITOR\" to open file")
			os.Exit(0)
		} else {
			cmd := exec.Command(editorCmd, g.statusItems[cursor].Path)
			// load env vars
			cmd.Env = os.Environ()
			for _, env := range g.envs {
				cmd.Env = append(cmd.Env, env["name"]+"="+env["value"])
			}
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				os.Exit(1)
			}
			os.Exit(0)
		}
	})
	return nil
}

// Render is the implementation of widget.Render
func (g *GitStatusWidget) Render() (err error) {
	body, err := g.buildBody(g.execPath)
	if err != nil {
		return
	}
	if body == nil {
		body = []string{
			"Worktree is clean",
		}
	}

	g.renderer.SetBody(body)
	g.renderer.ResetRender()

	return
}

// GetWidth is the implementation of widget.Render
func (g *GitStatusWidget) GetWidth() int {
	return g.renderer.GetWidth()
}

// GetHeight is the implementation of widget.Render
func (g *GitStatusWidget) GetHeight() int {
	return g.renderer.GetHeight()
}
