package dashboard

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode/utf8"

	ui "github.com/gizak/termui"
	"github.com/kr/pty"
	m2s "github.com/mitchellh/mapstructure"
	// "github.com/k0kubun/pp"
)

// MenuWidget is a command launcher
type MenuWidget struct {
	renderer     *ListWrapper
	menus        []Menu
	headerHeight int
	isReady      bool
	disabled     bool
	envs         []map[string]string
}

// NewMenuWidget constructs a New MenuWidget
func NewMenuWidget(wi Widget, envs []map[string]string) (m *MenuWidget, err error) {
	m = new(MenuWidget)
	if err := m2s.Decode(wi.Content, &m.menus); err != nil {
		return nil, err
	}
	h := m.buildHeader()
	m.headerHeight = len(h)
	m.envs = envs
	opt := &ListWrapperOption{
		Title:         wi.Title,
		RealHeight:    wi.RealHeight,
		Header:        h,
		Body:          m.buildBody(),
		LineHighLight: true,
	}
	m.renderer = NewListWrapper(opt)
	m.isReady = true

	return
}

// Activate is the implementation of Widget.Activate
func (m *MenuWidget) Activate() {
	m.setKeyBindings()
	m.renderer.Render()
}

// Deactivate is the implementation of Widget.Activate
func (m *MenuWidget) Deactivate() {
	m.renderer.ResetRender()
}

// IsDisabled is the implementation of Widget.IsDisabled
func (m *MenuWidget) IsDisabled() bool {
	return m.disabled
}

// IsReady is the implementation of Widget.IsReady
func (m *MenuWidget) IsReady() bool {
	return m.isReady
}

// GetHighlightenPos is the implementation of Widget.GetHighlightenPos
func (m *MenuWidget) GetHighlightenPos() int {
	return m.renderer.GetCursor() + m.headerHeight
}

// GetWidget is the implementation of widget.Activate
func (m *MenuWidget) GetWidget() *ui.List {
	return m.renderer.GetWidget()
}

func (m *MenuWidget) setKeyBindings() error {
	// exec command by Enter
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		ui.StopLoop()
		ui.Close()

		cursor := m.renderer.GetCursor()
		// straighten multi line commands
		cmds := strings.Split(m.menus[cursor].Command, "\n")
		cmdStr := ""
		for _, c := range cmds {
			if c != "" {
				cmdStr = cmdStr + c + "; "
			}
		}

		fmt.Println("---------- executing --------------")
		fmt.Println(cmdStr)
		fmt.Println("-----------------------------------")
		fmt.Println("")

		cmd := exec.Command("sh", "-c", cmdStr)

		// load env vars
		cmd.Env = os.Environ()
		for _, env := range m.envs {
			cmd.Env = append(cmd.Env, env["name"]+"="+env["value"])
		}

		tty, err := pty.Start(cmd)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		defer tty.Close()

		go func() {
			scanner := bufio.NewScanner(tty)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()
		go func() {
			io.Copy(tty, os.Stdin)
		}()

		err = cmd.Wait()
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		os.Exit(0)
	})
	return nil
}

func (m *MenuWidget) buildHeader() (header []string) {
	n1, n2, n3 := m.getLongest()
	colCt := m.fillSpaces("CATEGORY ", n1)
	colName := m.fillSpaces("NAME ", n2)
	colDesc := m.fillSpaces("DESCRIPTION ", n3)

	header = []string{
		" [NO" + " | " + colCt + " | " + colName + " | " + colDesc + "](fg-blue)\n",
		" [" + strings.Repeat("-", 500) + "](fg-blue)\n"}
	return
}

func (m *MenuWidget) buildBody() (body []string) {
	n1, n2, _ := m.getLongest()
	for k, v := range m.menus {
		var no string
		if k < 9 {
			no = " 0" + strconv.Itoa(k+1)
		} else {
			no = " " + strconv.Itoa(k+1)
		}
		ct := m.fillSpaces(v.Category, n1)
		name := m.fillSpaces(v.Name, n2)
		desc := v.Description + strings.Repeat(" ", 200)
		r := "[" + no + " |](fg-blue) " + ct + " | " + name + " | " + desc + "\n"
		body = append(body, r)
	}

	return
}

func (m *MenuWidget) fillSpaces(s string, longest int) string {
	var l = longest - utf8.RuneCountInString(s)
	for i := 0; i < l; i++ {
		s += " "
	}
	return s
}

func (m *MenuWidget) getLongest() (n1 int, n2 int, n3 int) {
	n1 = utf8.RuneCountInString("CATEGORY") + 1
	n2 = utf8.RuneCountInString("NAME") + 1
	n3 = utf8.RuneCountInString("DESCRIPTION") + 1
	for _, menu := range m.menus {
		c := utf8.RuneCountInString(menu.Category)
		nm := utf8.RuneCountInString(menu.Name)
		d := utf8.RuneCountInString(menu.Description)
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
